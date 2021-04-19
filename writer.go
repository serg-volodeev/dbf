package dbf

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"golang.org/x/text/encoding"
)

type Writer struct {
	header   *header
	fields   []*field
	writer   *bufio.Writer
	ws       io.WriteSeeker
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
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields defined")
	}
	w := &Writer{
		header: newHeader(),
		ws:     ws,
		writer: bufio.NewWriter(ws),
	}
	for _, f := range fields {
		if err := w.addField(f.Name, f.Type, f.Len, f.Dec); err != nil {
			return nil, err
		}
	}
	if codePage > 0 {
		cm := charmapByPage(codePage)
		if cm == nil {
			return nil, fmt.Errorf("unsupported code page %d", codePage)
		}
		w.encoder = cm.NewEncoder()
		w.header.setCodePage(codePage)
	}
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

func (w *Writer) Write(record []interface{}) error {
	if w.recCount == 0 {
		if err := w.initWriter(); err != nil {
			return err
		}
	}
	return w.writeRecord(record)
}

func (w *Writer) initWriter() error {
	w.header.setFieldCount(len(w.fields))
	w.header.RecSize = w.calcRecSize()
	if err := w.header.write(w.writer); err != nil {
		return err
	}
	if err := w.writeFields(); err != nil {
		return err
	}
	if err := w.writer.WriteByte(headerEnd); err != nil {
		return err
	}
	w.buf = make([]byte, int(w.header.RecSize))
	w.buf[0] = ' ' // deleted mark
	return nil
}

func (w *Writer) calcRecSize() uint16 {
	size := 1 // deleted mark
	for _, f := range w.fields {
		size += int(f.Len)
	}
	return uint16(size)
}

func (w *Writer) writeFields() error {
	offset := 1 // deleted mark
	for _, f := range w.fields {
		f.Offset = uint32(offset)
		if err := f.write(w.writer); err != nil {
			return err
		}
		offset += int(f.Len)
	}
	return nil
}

func (w *Writer) writeRecord(record []interface{}) error {
	if len(record) != len(w.fields) {
		return fmt.Errorf("the record does not match the number of fields")
	}
	for i := range w.fields {
		if err := w.setFieldValue(i, record[i]); err != nil {
			return err
		}
	}
	if _, err := w.writer.Write(w.buf); err != nil {
		return err
	}
	w.recCount++
	return nil
}

func (w *Writer) setFieldValue(index int, value interface{}) error {
	var fieldStr string
	f := w.fields[index]

	switch f.Type {
	case 'C':
		v, err := interfaceToString(value)
		if err != nil {
			return err
		}
		if w.encoder != nil && !isASCII(v) {
			v, err = w.encoder.String(v)
			if err != nil {
				return err
			}
		}
		if err := f.checkLen(v); err != nil {
			return err
		}
		fieldStr = padRight(v, int(f.Len))
	case 'L':
		v, err := interfaceToBool(value)
		if err != nil {
			return err
		}
		fieldStr = "F"
		if v {
			fieldStr = "T"
		}
	case 'D':
		v, err := interfaceToDate(value)
		if err != nil {
			return err
		}
		fieldStr = v.Format("20060102")
	case 'N':
		var s string
		if f.Dec == 0 {
			v, err := interfaceToInt(value)
			if err != nil {
				return err
			}
			s = strconv.FormatInt(v, 10)
		} else {
			v, err := interfaceToFloat(value)
			if err != nil {
				return err
			}
			s = strconv.FormatFloat(v, 'f', int(f.Dec), 64)
		}
		if err := f.checkLen(s); err != nil {
			return err
		}
		fieldStr = padLeft(s, int(f.Len))
	default:
		return fmt.Errorf("invalid field type: got %s, want C, N, L, D", string(f.Type))
	}
	copy(w.buf[int(f.Offset):int(f.Offset)+int(f.Len)], fieldStr)
	return nil
}

func (w *Writer) Flush() error {
	if err := w.writer.WriteByte(fileEnd); err != nil {
		return err
	}
	if err := w.writer.Flush(); err != nil {
		return err
	}
	// modify record count in header
	if _, err := w.ws.Seek(0, 0); err != nil {
		return err
	}
	w.header.RecCount = w.recCount
	if err := w.header.write(w.ws); err != nil {
		return err
	}
	return nil
}
