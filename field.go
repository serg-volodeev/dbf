package dbf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
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

func newLogicalField(name string) (*field, error) {
	if err := checkName(name); err != nil {
		return nil, err
	}
	f := &field{}
	f.setName(name)
	f.Type = 'L'
	f.Len = 1
	return f, nil
}

func newDateField(name string) (*field, error) {
	if err := checkName(name); err != nil {
		return nil, err
	}
	f := &field{}
	f.setName(name)
	f.Type = 'D'
	f.Len = 8
	return f, nil
}

func newCharacterField(name string, length int) (*field, error) {
	if err := checkName(name); err != nil {
		return nil, err
	}
	if length <= 0 || length > maxCharacterLen {
		return nil, fmt.Errorf("field len %d, want 0 < len <= %d", length, maxCharacterLen)
	}
	f := &field{}
	f.setName(name)
	f.Type = 'C'
	f.Len = byte(length)
	return f, nil
}

func newNumericField(name string, length, dec int) (*field, error) {
	if err := checkName(name); err != nil {
		return nil, err
	}
	if length <= 0 || length > maxNumericLen {
		return nil, fmt.Errorf("field len %d, want 0 < len <= %d", length, maxNumericLen)
	}
	if dec < 0 {
		return nil, fmt.Errorf("field dec %d, want dec > 0", dec)
	}
	if length <= 2 && dec > 0 {
		return nil, fmt.Errorf("field dec %d, want 0", dec)
	}
	if length > 2 && (dec > length-2) {
		return nil, fmt.Errorf("field dec %d, want dec <= %d", dec, length-2)
	}
	f := &field{}
	f.setName(name)
	f.Type = 'N'
	f.Len = byte(length)
	f.Dec = byte(dec)
	return f, nil
}

// Field name

func checkName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return fmt.Errorf("empty field name")
	}
	if len(name) > maxNameLen {
		return fmt.Errorf("too long field name %q, max len %d", name, maxNameLen)
	}
	return nil
}

func (f *field) name() string {
	i := bytes.IndexByte(f.Name[:], 0)
	return string(f.Name[:i])
}

func (f *field) setName(name string) {
	name = strings.ToUpper(strings.TrimSpace(name))
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

// Check field

func (f *field) checkLen(value string) error {
	if len(value) > int(f.Len) {
		return fmt.Errorf("field value %q overflow: value len %d, field len %d", value, len(value), int(f.Len))
	}
	return nil
}

func (f *field) checkType(t byte) error {
	if t != f.Type {
		return fmt.Errorf("field type %q, want: %q", f.Type, t)
	}
	return nil
}

// Get field value

func (f *field) fieldBuf(recordBuf []byte) []byte {
	return recordBuf[int(f.Offset) : int(f.Offset)+int(f.Len)]
}

func (f *field) stringFieldValue(recordBuf []byte, decoder *encoding.Decoder) (string, error) {
	buf := f.fieldBuf(recordBuf)
	switch f.Type {
	case 'C':
		var err error
		s := trimRight(buf)
		if decoder != nil && !isASCII(s) {
			s, err = decoder.String(s)
			if err != nil {
				return "", err
			}
		}
		return s, nil
	case 'L', 'D', 'N':
		return trimLeft(buf), nil
	}
	return "", fmt.Errorf("unknow type %q, want 'C', 'L', 'D', 'N'", f.Type)
}

func (f *field) boolFieldValue(recordBuf []byte) (bool, error) {
	if err := f.checkType('L'); err != nil {
		return false, err
	}
	buf := f.fieldBuf(recordBuf)
	switch buf[0] {
	case 'T', 't', 'Y', 'y':
		return true, nil
	}
	return false, nil
}

func (f *field) dateFieldValue(recordBuf []byte) (time.Time, error) {
	if err := f.checkType('D'); err != nil {
		return time.Time{}, err
	}
	buf := f.fieldBuf(recordBuf)
	if isEmpty(buf) {
		return time.Time{}, nil
	}
	return time.Parse("20060102", string(buf))
}

func (f *field) intFieldValue(recordBuf []byte) (int64, error) {
	if err := f.checkType('N'); err != nil {
		return 0, err
	}
	buf := f.fieldBuf(recordBuf)
	if f.Dec != 0 {
		buf = buf[:len(buf)-int(f.Dec)-1]
	}
	s := trimLeft(buf)
	if s == "" {
		s = "0"
	}
	return strconv.ParseInt(s, 10, 64)
}

func (f *field) floatFieldValue(recordBuf []byte) (float64, error) {
	if err := f.checkType('N'); err != nil {
		return 0, err
	}
	buf := f.fieldBuf(recordBuf)
	s := trimLeft(buf)
	if s == "" {
		s = "0"
	}
	return strconv.ParseFloat(s, 64)
}

// Set field value

func (f *field) setFieldBuf(recordBuf []byte, value string) {
	copy(recordBuf[int(f.Offset):int(f.Offset)+int(f.Len)], value)
}

func (f *field) setStringFieldValue(recordBuf []byte, value string, encoder *encoding.Encoder) error {
	switch f.Type {
	case 'C':
		var err error
		s := value
		if encoder != nil && !isASCII(s) {
			s, err = encoder.String(s)
			if err != nil {
				return err
			}
		}
		if err := f.checkLen(s); err != nil {
			return err
		}
		s = padRight(s, int(f.Len))
		f.setFieldBuf(recordBuf, s)
	case 'L':
		s := strings.TrimSpace(value)
		if s == "" {
			s = " "
		} else {
			switch s[0] {
			case 'T', 't', 'Y', 'y':
				s = "T"
			default:
				s = "F"
			}
		}
		f.setFieldBuf(recordBuf, s)
	case 'D':
		s := strings.TrimSpace(value)
		if s == "" {
			s = strings.Repeat(" ", int(f.Len))
		} else {
			_, err := time.Parse("20060102", s)
			if err != nil {
				return err
			}
		}
		f.setFieldBuf(recordBuf, s)
	case 'N':
		s := strings.TrimSpace(value)
		if s == "" {
			s = "0"
		}
		if f.Dec == 0 {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			return f.setIntFieldValue(recordBuf, n)
		} else {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return err
			}
			return f.setFloatFieldValue(recordBuf, n)
		}
	default:
		return fmt.Errorf("unknow type %q, want 'C', 'L', 'D', 'N'", f.Type)
	}
	return nil
}

func (f *field) setBoolFieldValue(recordBuf []byte, value bool) error {
	if err := f.checkType('L'); err != nil {
		return err
	}
	s := "F"
	if value {
		s = "T"
	}
	f.setFieldBuf(recordBuf, s)
	return nil
}

func (f *field) setDateFieldValue(recordBuf []byte, value time.Time) error {
	if err := f.checkType('D'); err != nil {
		return err
	}
	s := value.Format("20060102")
	f.setFieldBuf(recordBuf, s)
	return nil
}

func (f *field) setIntFieldValue(recordBuf []byte, value int64) error {
	if err := f.checkType('N'); err != nil {
		return err
	}
	s := strconv.FormatInt(value, 10)
	if f.Dec > 0 {
		s += "." + strings.Repeat("0", int(f.Dec))
	}
	if err := f.checkLen(s); err != nil {
		return err
	}
	s = padLeft(s, int(f.Len))
	f.setFieldBuf(recordBuf, s)
	return nil
}

func (f *field) setFloatFieldValue(recordBuf []byte, value float64) error {
	if err := f.checkType('N'); err != nil {
		return err
	}
	s := strconv.FormatFloat(value, 'f', int(f.Dec), 64)
	if err := f.checkLen(s); err != nil {
		return err
	}
	s = padLeft(s, int(f.Len))
	f.setFieldBuf(recordBuf, s)
	return nil
}
