package dbf

import (
	"bufio"
	"fmt"
	"io"

	"golang.org/x/text/encoding"
)

const fileEnd byte = 0x1A

// A Writer writes records in DBF file.
// The writes of individual records are buffered.
// After all data has been written, the client should
// call the Flush method to guarantee all data has been
// forwarded to the underlying io.Writer.
type Writer struct {
	header   *header
	fields   *Fields
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
// If the codePage parameter is zero, the text fields will not be encoded.
func NewWriter(ws io.WriteSeeker, fields *Fields, codePage int) (*Writer, error) {
	if _, ok := ws.(io.WriteSeeker); !ok {
		return nil, fmt.Errorf("parameter %v is not io.WriteSeeker", ws)
	}
	if fields.Count() == 0 {
		return nil, fmt.Errorf("no fields defined")
	}
	w := &Writer{
		header: newHeader(),
		fields: fields,
		ws:     ws,
		writer: bufio.NewWriter(ws),
	}
	if codePage > 0 {
		cm := charmapByPage(codePage)
		if cm == nil {
			return nil, fmt.Errorf("unsupported code page %d", codePage)
		}
		w.encoder = cm.NewEncoder()
		w.header.setCodePage(codePage)
	}
	w.header.setFieldCount(w.fields.Count())
	w.header.RecSize = w.fields.calcRecSize()
	if err := w.header.write(w.writer); err != nil {
		return nil, err
	}
	if err := w.fields.write(w.writer); err != nil {
		return nil, err
	}
	if err := w.writer.WriteByte(headerEnd); err != nil {
		return nil, err
	}
	w.buf = make([]byte, int(w.header.RecSize))
	w.buf[0] = ' ' // deleted mark
	return w, nil
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

func (w *Writer) writeRecord(record []interface{}) error {
	if len(record) != w.fields.Count() {
		return fmt.Errorf("the record does not match the number of fields")
	}
	for i := range w.fields.items {
		if err := w.setFieldValue(i, record[i]); err != nil {
			return fmt.Errorf("field %q: %w", w.fields.items[i].name(), err)
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
	var err error
	f := w.fields.items[index]

	switch f.Type {
	case 'C':
		s, err = f.characterToBuf(value, w.encoder)
	case 'L':
		s, err = f.logicalToBuf(value)
	case 'D':
		s, err = f.dateToBuf(value)
	case 'N':
		s, err = f.numericToBuf(value)
	default:
		return fmt.Errorf("invalid field type: got %s, want C, N, L, D", string(f.Type))
	}
	if err != nil {
		return err
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
