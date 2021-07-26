package dbf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/text/encoding"
)

const (
	fieldSize = 32

	maxNameLen      = 10
	maxCharacterLen = 254
	maxNumericLen   = 19
)

type field struct {
	Name   [11]byte
	Type   byte
	Offset uint32
	Len    byte
	Dec    byte
	Filler [14]byte
}

// New field

func newLogicalField(name string) *field {
	f := &field{}
	f.setName(name)
	f.Type = 'L'
	f.Len = 1
	f.Dec = 0
	return f
}

func newDateField(name string) *field {
	f := &field{}
	f.setName(name)
	f.Type = 'D'
	f.Len = 8
	f.Dec = 0
	return f
}

func newCharacterField(name string, length int) *field {
	if length <= 0 || length > maxCharacterLen {
		panic(fmt.Errorf("invalid field len: got %d, want 0 < len <= %d", length, maxCharacterLen))
	}
	f := &field{}
	f.setName(name)
	f.Type = 'C'
	f.Len = byte(length)
	f.Dec = 0
	return f
}

func newNumericField(name string, length, dec int) *field {
	if length <= 0 || length > maxNumericLen {
		panic(fmt.Errorf("invalid field len: got %d, want 0 < len <= %d", length, maxNumericLen))
	}
	if dec < 0 {
		panic(fmt.Errorf("invalid field dec: got %d, want dec > 0", dec))
	}
	if length <= 2 && dec > 0 {
		panic(fmt.Errorf("invalid field dec: got %d, want 0", dec))
	}
	if length > 2 && (dec > length-2) {
		panic(fmt.Errorf("invalid field dec: got %d, want dec <= %d", dec, length-2))
	}
	f := &field{}
	f.setName(name)
	f.Type = 'N'
	f.Len = byte(length)
	f.Dec = byte(dec)
	return f
}

// Field name

func (f *field) name() string {
	i := bytes.IndexByte(f.Name[:], 0)
	return string(f.Name[:i])
}

func (f *field) setName(name string) {
	name = strings.ToUpper(strings.TrimSpace(name))
	if len(name) == 0 {
		panic(fmt.Errorf("empty field name"))
	}
	if len(name) > maxNameLen {
		panic(fmt.Errorf("too long field name: %q, max len %d", name, maxNameLen))
	}
	copy(f.Name[:], name)
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

func (f *field) checkLen(value string) error {
	if len(value) > int(f.Len) {
		return fmt.Errorf("field value overflow: value len %d, field len %d", len(value), int(f.Len))
	}
	return nil
}

// Value to buffer

func (f *field) characterToBuf(value interface{}, encoder *encoding.Encoder) (string, error) {
	var err error
	v, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("error convert %v to string", value)
	}
	if encoder != nil && !isASCII(v) {
		v, err = encoder.String(v)
		if err != nil {
			return "", err
		}
	}
	if err := f.checkLen(v); err != nil {
		return "", err
	}
	v = padRight(v, int(f.Len))
	return v, nil
}

func (f *field) logicalToBuf(value interface{}) (string, error) {
	v, ok := value.(bool)
	if !ok {
		return "", fmt.Errorf("error convert %v to bool", value)
	}
	if v {
		return "T", nil
	}
	return "F", nil
}

func (f *field) dateToBuf(value interface{}) (string, error) {
	v, ok := value.(time.Time)
	if !ok {
		return "", fmt.Errorf("error convert %v to date", value)
	}
	return v.Format("20060102"), nil
}

func (f *field) numericToBuf(value interface{}) (string, error) {
	var s string
	if f.Dec == 0 {
		s = fmt.Sprintf("%d", value)
	} else {
		format := fmt.Sprintf("%%.%df", f.Dec)
		s = fmt.Sprintf(format, value)
	}
	if err := f.checkLen(s); err != nil {
		return "", err
	}
	s = padLeft(s, int(f.Len))
	return s, nil
}
