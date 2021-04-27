package dbf

type Fields struct {
	items []field
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
