package dbf

import "io"

type Fields struct {
	items []*field
	err   error
}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) Count() int {
	return len(f.items)
}

func (f *Fields) Error() error {
	return f.err
}

func (f *Fields) Add(name string, typ string, opts ...int) {
	if f.err != nil {
		return
	}
	length := 0
	dec := 0
	if len(opts) > 0 {
		length = opts[0]
	}
	if len(opts) > 1 {
		dec = opts[1]
	}
	item, err := newField(name, typ, length, dec)
	if err != nil {
		f.err = err
		return
	}
	f.items = append(f.items, item)
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

func (f *Fields) Get(i int) (name, typ string, length, dec int) {
	item := f.items[i]
	name = item.name()
	typ = string(item.Type)
	length = int(item.Len)
	dec = int(item.Dec)
	return
}
