package dbf

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"golang.org/x/text/encoding"
)

// A Writer writes records in DBF file.
// The writes of individual records are buffered.
// After all data has been written, the client should
// call the Flush method to guarantee all data has been
// forwarded to the underlying io.Writer.
type Writer struct {
	header   *header
	fields   []*field
	writer   *bufio.Writer
	ws       io.WriteSeeker
	buf      []byte
	encoder  *encoding.Encoder
	recCount uint32
}

// NewWriter returns a new Writer that writes to w.
// The function writes the header of the DBF file.
// If you call the Flash method afterwards, an empty file will be created.
//
// Supported code pages:
//     437   - US MS-DOS
//     850   - International MS-DOS
//     1252  - Windows ANSI
//     10000 - Standard Macintosh
//     852   - Easern European MS-DOS
//     866   - Russian MS-DOS
//     865   - Nordic MS-DOS
//     1255  - Hebrew Windows
//     1256  - Arabic Windows
//     10007 - Russian Macintosh
//     1250  - Eastern European Windows
//     1251  - Russian Windows
//     1254  - Turkish Windows
//     1253  - Greek Windows
//
// If the kp parameter is zero, the text fields will not be encoded.
func NewWriter(ws io.WriteSeeker, fields []FieldInfo, codePage int) (*Writer, error) {
	if _, ok := ws.(io.WriteSeeker); !ok {
		return nil, fmt.Errorf("parameter %v is not io.WriteSeeker", ws)
	}
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
	if err := w.initWriter(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Writer) addField(name string, typ string, length int, dec int) error {
	f, err := newField(name, typ, length, dec)
	if err != nil {
		return fmt.Errorf("field %q: %w", name, err)
	}
	w.fields = append(w.fields, f)
	return nil
}

// Write writes a single record to w.
// A record is a slice of interface{} with each value being one field.
// Writes are buffered, so Flush must eventually be called to ensure
// that the record is written to the underlying io.Writer.
func (w *Writer) Write(record []interface{}) error {
	if err := w.writeRecord(record); err != nil {
		return fmt.Errorf("record %d: %w", w.recCount+1, err)
	}
	return nil
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
			return fmt.Errorf("field %q: %w", w.fields[i].name(), err)
		}
	}
	if _, err := w.writer.Write(w.buf); err != nil {
		return err
	}
	w.recCount++
	return nil
}

func (w *Writer) setFieldValue(index int, value interface{}) error {
	var s string
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
		s = padRight(v, int(f.Len))
	case 'L':
		v, err := interfaceToBool(value)
		if err != nil {
			return err
		}
		s = "F"
		if v {
			s = "T"
		}
	case 'D':
		v, err := interfaceToDate(value)
		if err != nil {
			return err
		}
		s = v.Format("20060102")
	case 'N':
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
		s = padLeft(s, int(f.Len))
	default:
		return fmt.Errorf("invalid field type: got %s, want C, N, L, D", string(f.Type))
	}
	copy(w.buf[int(f.Offset):int(f.Offset)+int(f.Len)], s)
	return nil
}

// Flush writes any buffered data to the underlying io.Writer.
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
