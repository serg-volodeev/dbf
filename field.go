package xbase

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	maxFieldNameLen = 10
	maxCFieldLen    = 254
	maxNFieldLen    = 19
)

const (
	defaultLFieldLen = 1
	defaultDFieldLen = 8
)

type field struct {
	Name   [11]byte
	Type   byte
	Offset uint32
	Len    byte
	Dec    byte
	Filler [14]byte
}

func (f *field) name() string {
	i := bytes.IndexByte(f.Name[:], 0)
	return string(f.Name[:i])
}

// String utils

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// New field

func newField(name string, typ string, length, dec int) (*field, error) {
	f := &field{}
	// do not change the call order
	if err := f.setName(name); err != nil {
		return nil, err
	}
	if err := f.setType(typ); err != nil {
		return nil, err
	}
	if err := f.setLen(length); err != nil {
		return nil, err
	}
	if err := f.setDec(dec); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *field) setName(name string) error {
	name = strings.ToUpper(strings.TrimSpace(name))
	if len(name) == 0 {
		return fmt.Errorf("empty field name")
	}
	if len(name) > maxFieldNameLen {
		return fmt.Errorf("too long field name: %q, max len %d", name, maxFieldNameLen)
	}
	copy(f.Name[:], name)
	return nil
}

func (f *field) setType(typ string) error {
	typ = strings.ToUpper(strings.TrimSpace(typ))
	if len(typ) == 0 {
		return fmt.Errorf("empty field type")
	}
	t := typ[0]
	if bytes.IndexByte([]byte("CNLD"), t) < 0 {
		return fmt.Errorf("invalid field type: got %s, want C, N, L, D", string(t))
	}
	f.Type = t
	return nil
}

func (f *field) setLen(length int) error {
	switch f.Type {
	case 'C':
		if length <= 0 || length > maxCFieldLen {
			return fmt.Errorf("invalid field len: got %d, want 0 < len <= %d", length, maxCFieldLen)
		}
	case 'N':
		if length <= 0 || length > maxNFieldLen {
			return fmt.Errorf("invalid field len: got %d, want 0 < len <= %d", length, maxNFieldLen)
		}
	case 'L':
		length = defaultLFieldLen
	case 'D':
		length = defaultDFieldLen
	default:
		return fmt.Errorf("field type not defined")
	}
	f.Len = byte(length)
	return nil
}

func (f *field) setDec(dec int) error {
	if f.Type == 'N' {
		if dec < 0 {
			return fmt.Errorf("invalid field dec: got %d, want dec > 0", dec)
		}
		length := int(f.Len)
		if length <= 2 && dec > 0 {
			return fmt.Errorf("invalid field dec: got %d, want 0", dec)
		}
		if length > 2 && (dec > length-2) {
			return fmt.Errorf("invalid field dec: got %d, want dec <= %d", dec, length-2)
		}
	} else {
		dec = 0
	}
	f.Dec = byte(dec)
	return nil
}

// Read/write

func (f *field) read(reader io.Reader) error {
	return binary.Read(reader, binary.LittleEndian, f)
}

func (f *field) write(writer io.Writer) error {
	tmp := f.Offset
	f.Offset = 0
	defer func() { f.Offset = tmp }()
	return binary.Write(writer, binary.LittleEndian, f)
}
