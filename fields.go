package dbf

import (
	"fmt"
	"io"
	"time"

	"golang.org/x/text/encoding"
)

// Fields for creating file structure.
type Fields struct {
	items   []*field
	recSize int
	err     error
}

// NewFields returns a pointer to a structure Fields.
func NewFields() *Fields {
	return &Fields{recSize: 1}
}

// Err returns the first error that was encountered by the Fields.
func (f *Fields) Err() error {
	if f.err != nil {
		return fmt.Errorf("dbf.Fields: %w", f.err)
	}
	return nil
}

// Count returns the number of fields.
func (f *Fields) Count() int {
	return len(f.items)
}

func (f *Fields) addItem(item *field) error {
	if f.nameExists(item.name()) {
		return fmt.Errorf("duplicate field name %q", item.name())
	}
	item.Offset = uint32(f.recSize)
	f.recSize += int(item.Len)
	f.items = append(f.items, item)
	return nil
}

func (f *Fields) nameExists(name string) bool {
	for _, item := range f.items {
		if item.name() == name {
			return true
		}
	}
	return false
}

// AddLogicalField adds a logical field to the structure.
func (f *Fields) AddLogicalField(name string) {
	if f.err != nil {
		return
	}
	item, err := newLogicalField(name)
	if err != nil {
		f.err = fmt.Errorf("AddLogicalField: %w", err)
		return
	}
	if err := f.addItem(item); err != nil {
		f.err = fmt.Errorf("AddLogicalField: %w", err)
		return
	}
}

// AddDateField adds a date field to the structure.
func (f *Fields) AddDateField(name string) {
	if f.err != nil {
		return
	}
	item, err := newDateField(name)
	if err != nil {
		f.err = fmt.Errorf("AddDateField: %w", err)
		return
	}
	if err := f.addItem(item); err != nil {
		f.err = fmt.Errorf("AddDateField: %w", err)
		return
	}
}

// AddCharacterField adds a character field to the structure.
func (f *Fields) AddCharacterField(name string, length int) {
	if f.err != nil {
		return
	}
	item, err := newCharacterField(name, length)
	if err != nil {
		f.err = fmt.Errorf("AddCharacterField: %w", err)
		return
	}
	if err := f.addItem(item); err != nil {
		f.err = fmt.Errorf("AddCharacterField: %w", err)
		return
	}
}

// AddNumericField adds a numeric field to the structure.
func (f *Fields) AddNumericField(name string, length, dec int) {
	if f.err != nil {
		return
	}
	item, err := newNumericField(name, length, dec)
	if err != nil {
		f.err = fmt.Errorf("AddNumericField: %w", err)
		return
	}
	if err := f.addItem(item); err != nil {
		f.err = fmt.Errorf("AddNumericField: %w", err)
		return
	}
}

// FieldInfo returns field information by index.
func (f *Fields) FieldInfo(index int) (name, typ string, length, dec int) {
	if f.err != nil {
		return
	}
	if err := f.checkFieldIndex(index); err != nil {
		f.err = fmt.Errorf("FieldInfo: %w", err)
		return
	}
	item := f.items[index]
	name = item.name()
	typ = string(item.Type)
	length = int(item.Len)
	dec = int(item.Dec)
	return
}

func (f *Fields) write(w io.Writer) error {
	for _, item := range f.items {
		if err := item.write(w); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fields) read(r io.Reader, count int) error {
	for i := 0; i < count; i++ {
		item := &field{}
		if err := item.read(r); err != nil {
			return err
		}
		f.addItem(item)
	}
	return nil
}

func (f *Fields) checkFieldIndex(index int) error {
	if index < 0 || index >= f.Count() {
		return fmt.Errorf("field index out of range [%d] with field count %d", index, f.Count())
	}
	return nil
}

// Get value

func (f *Fields) stringFieldValue(index int, recordBuf []byte, decoder *encoding.Decoder) (string, error) {
	if err := f.checkFieldIndex(index); err != nil {
		return "", err
	}
	return f.items[index].stringFieldValue(recordBuf, decoder)
}

func (f *Fields) boolFieldValue(index int, recordBuf []byte) (bool, error) {
	if err := f.checkFieldIndex(index); err != nil {
		return false, err
	}
	return f.items[index].boolFieldValue(recordBuf)
}

func (f *Fields) dateFieldValue(index int, recordBuf []byte) (time.Time, error) {
	if err := f.checkFieldIndex(index); err != nil {
		return time.Time{}, err
	}
	return f.items[index].dateFieldValue(recordBuf)
}

func (f *Fields) intFieldValue(index int, recordBuf []byte) (int64, error) {
	if err := f.checkFieldIndex(index); err != nil {
		return 0, err
	}
	return f.items[index].intFieldValue(recordBuf)
}

func (f *Fields) floatFieldValue(index int, recordBuf []byte) (float64, error) {
	if err := f.checkFieldIndex(index); err != nil {
		return 0, err
	}
	return f.items[index].floatFieldValue(recordBuf)
}

// Set value

func (f *Fields) setStringFieldValue(index int, recordBuf []byte, value string, encoder *encoding.Encoder) error {
	if err := f.checkFieldIndex(index); err != nil {
		return err
	}
	return f.items[index].setStringFieldValue(recordBuf, value, encoder)
}

func (f *Fields) setBoolFieldValue(index int, recordBuf []byte, value bool) error {
	if err := f.checkFieldIndex(index); err != nil {
		return err
	}
	return f.items[index].setBoolFieldValue(recordBuf, value)
}

func (f *Fields) setDateFieldValue(index int, recordBuf []byte, value time.Time) error {
	if err := f.checkFieldIndex(index); err != nil {
		return err
	}
	return f.items[index].setDateFieldValue(recordBuf, value)
}

func (f *Fields) setIntFieldValue(index int, recordBuf []byte, value int64) error {
	if err := f.checkFieldIndex(index); err != nil {
		return err
	}
	return f.items[index].setIntFieldValue(recordBuf, value)
}

func (f *Fields) setFloatFieldValue(index int, recordBuf []byte, value float64) error {
	if err := f.checkFieldIndex(index); err != nil {
		return err
	}
	return f.items[index].setFloatFieldValue(recordBuf, value)
}
