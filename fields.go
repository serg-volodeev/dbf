package dbf

import (
	"fmt"
	"io"

	"golang.org/x/text/encoding"
)

// Fields for creating file structure.
type Fields struct {
	items   []*field
	curOffs int
}

// NewFields returns a pointer to a structure Fields.
func NewFields() *Fields {
	return &Fields{curOffs: 1}
}

// Count returns the number of fields.
func (f *Fields) Count() int {
	return len(f.items)
}

func (f *Fields) addItem(item *field) {
	item.Offset = uint32(f.curOffs)
	f.curOffs += int(item.Len)
	f.items = append(f.items, item)
}

// AddLogicalField adds a logical field to the structure.
func (f *Fields) AddLogicalField(name string) {
	f.addItem(newLogicalField(name))
}

// AddDateField adds a date field to the structure.
func (f *Fields) AddDateField(name string) {
	f.addItem(newDateField(name))
}

// AddCharacterField adds a character field to the structure.
func (f *Fields) AddCharacterField(name string, length int) {
	f.addItem(newCharacterField(name, length))
}

// AddNumericField adds a numeric field to the structure.
func (f *Fields) AddNumericField(name string, length, dec int) {
	f.addItem(newNumericField(name, length, dec))
}

// FieldInfo returns field information by index.
func (f *Fields) FieldInfo(i int) (name, typ string, length, dec int) {
	item := f.items[i]
	name = item.name()
	typ = string(item.Type)
	length = int(item.Len)
	dec = int(item.Dec)
	return
}

func (f *Fields) write(w io.Writer) error {
	offset := 1 // deleted mark
	for _, item := range f.items {
		item.Offset = uint32(offset)
		if err := item.write(w); err != nil {
			return err
		}
		offset += int(item.Len)
	}
	return nil
}

func (f *Fields) read(r io.Reader, count int) error {
	offset := 1 // deleted mark
	for i := 0; i < count; i++ {
		item := &field{}
		if err := item.read(r); err != nil {
			return err
		}
		item.Offset = uint32(offset)
		f.items = append(f.items, item)
		offset += int(item.Len)
	}
	return nil
}

func (f *Fields) calcRecSize() uint16 {
	size := 1 // deleted mark
	for _, item := range f.items {
		size += int(item.Len)
	}
	return uint16(size)
}

func (f *Fields) copyRecordToBuf(buf []byte, record []interface{}, encoder *encoding.Encoder) error {
	if len(record) != f.Count() {
		return fmt.Errorf("the record does not match the number of fields")
	}
	for i, item := range f.items {
		var s string
		var err error

		value := record[i]

		switch item.Type {
		case 'C':
			s, err = item.characterToString(value, encoder)
		case 'L':
			s, err = item.logicalToString(value)
		case 'D':
			s, err = item.dateToString(value)
		case 'N':
			s, err = item.numericToString(value)
		default:
			err = fmt.Errorf("invalid field type: got %s, want C, N, L, D", string(item.Type))
		}
		if err != nil {
			return fmt.Errorf("field %q: %w", item.name(), err)
		}
		item.copyValueToBuf(buf, s)
	}
	return nil
}
