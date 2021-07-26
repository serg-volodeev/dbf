package dbf

import "io"

// Fields for creating file structure.
type Fields struct {
	items []*field
}

// NewFields returns a pointer to a structure Fields.
func NewFields() *Fields {
	return &Fields{}
}

// Count returns the number of fields.
func (f *Fields) Count() int {
	return len(f.items)
}

// AddLogicalField adds a logical field to the structure.
func (f *Fields) AddLogicalField(name string) {
	f.items = append(f.items, newLogicalField(name))
}

// AddDateField adds a date field to the structure.
func (f *Fields) AddDateField(name string) {
	f.items = append(f.items, newDateField(name))
}

// AddCharacterField adds a character field to the structure.
func (f *Fields) AddCharacterField(name string, length int) {
	f.items = append(f.items, newCharacterField(name, length))
}

// AddNumericField adds a numeric field to the structure.
func (f *Fields) AddNumericField(name string, length, dec int) {
	f.items = append(f.items, newNumericField(name, length, dec))
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
