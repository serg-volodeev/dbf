package xbase

import (
	"bufio"
	"fmt"
	"io"

	"golang.org/x/text/encoding"
)

type Writer struct {
	header   *header
	fields   []*field
	writer   *bufio.Writer
	buf      []byte
	encoder  *encoding.Encoder
	recCount uint32
}

type FieldInfo struct {
	Name string
	Type string
	Len  int
	Dec  int
}

func NewWriter(ws io.WriteSeeker, fields []FieldInfo, codePage int) (*Writer, error) {
	w := &Writer{
		header: newHeader(),
		writer: bufio.NewWriter(ws),
	}
	for _, f := range fields {
		if err := w.addField(f.Name, f.Type, f.Len, f.Dec); err != nil {
			return nil, err
		}
	}
	cm := charmapByPage(codePage)
	if cm == nil {
		return nil, fmt.Errorf("Unsupported code page %d", codePage)
	}
	w.encoder = cm.NewEncoder()
	w.header.setCodePage(codePage)
	return w, nil
}

func (w *Writer) addField(name string, typ string, length int, dec int) error {
	f, err := newField(name, typ, length, dec)
	if err != nil {
		return err
	}
	w.fields = append(w.fields, f)
	return nil
}
