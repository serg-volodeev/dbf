package xbase

import (
	"bytes"
	"fmt"
	"strings"
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
