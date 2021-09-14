package dbf

import (
	"bufio"
	"fmt"
	"io"
	"time"

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
	err      error
}

// NewWriter returns a new Writer that writes to ws.
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
func NewWriter(ws io.WriteSeeker, fields *Fields, codePage int) (w *Writer, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dbf.NewWriter: %w", err)
		}
	}()
	if ws == nil {
		return nil, fmt.Errorf("parameter is nil")
	}
	if fields.err != nil {
		return nil, fields.err
	}
	if fields.Count() == 0 {
		return nil, fmt.Errorf("no fields defined")
	}
	w = &Writer{
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
	w.header.RecSize = uint16(w.fields.recSize)

	if err = w.header.write(w.writer); err != nil {
		return nil, err
	}
	if err = w.fields.write(w.writer); err != nil {
		return nil, err
	}
	if err = w.writer.WriteByte(headerEnd); err != nil {
		return nil, err
	}
	w.buf = make([]byte, int(w.header.RecSize))
	w.clearBuf()
	return w, nil
}

func (w *Writer) clearBuf() {
	for i := range w.buf {
		w.buf[i] = ' '
	}
}

// Err returns the first error that was encountered by the Writer.
func (w *Writer) Err() error {
	if w.err != nil {
		return fmt.Errorf("dbf.Writer: %w", w.err)
	}
	return nil
}

// Write writes a single record to w.
func (w *Writer) Write() {
	if w.err != nil {
		return
	}
	if _, err := w.writer.Write(w.buf); err != nil {
		w.err = fmt.Errorf("Write: record %d: %w", w.recCount+1, err)
		return
	}
	w.recCount++
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() {
	if w.err != nil {
		return
	}
	if err := w.flush(); err != nil {
		w.err = fmt.Errorf("Flush: %w", err)
	}
}

func (w *Writer) flush() error {
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

// SetDeleted sets record mark is deleted.
func (w *Writer) SetDeteted(deleted bool) {
	if w.err != nil {
		return
	}
	if deleted {
		w.buf[0] = '*'
	} else {
		w.buf[0] = ' '
	}
}

// SetStringFieldValue assigns a value to a field by index.
// Field type must be Character, Logical, Date or Numeric.
func (w *Writer) SetStringFieldValue(index int, value string) {
	if w.err != nil {
		return
	}
	err := w.fields.setStringFieldValue(index, w.buf, value, w.encoder)
	if err != nil {
		w.err = fmt.Errorf("SetStringFieldValue: %w", err)
	}
}

// SetBoolFieldValue assigns a value to a field by index.
// Field type must be Logical.
func (w *Writer) SetBoolFieldValue(index int, value bool) {
	if w.err != nil {
		return
	}
	err := w.fields.setBoolFieldValue(index, w.buf, value)
	if err != nil {
		w.err = fmt.Errorf("SetBoolFieldValue: %w", err)
	}
}

// SetDateFieldValue assigns a value to a field by index.
// Field type must be Date.
func (w *Writer) SetDateFieldValue(index int, value time.Time) {
	if w.err != nil {
		return
	}
	err := w.fields.setDateFieldValue(index, w.buf, value)
	if err != nil {
		w.err = fmt.Errorf("SetDateFieldValue: %w", err)
	}
}

// SetIntFieldValue assigns a value to a field by index.
// Field type must be Numeric.
func (w *Writer) SetIntFieldValue(index int, value int64) {
	if w.err != nil {
		return
	}
	err := w.fields.setIntFieldValue(index, w.buf, value)
	if err != nil {
		w.err = fmt.Errorf("SetIntFieldValue: %w", err)
	}
}

// SetFloatFieldValue assigns a value to a field by index.
// Field type must be Numeric.
func (w *Writer) SetFloatFieldValue(index int, value float64) {
	if w.err != nil {
		return
	}
	err := w.fields.setFloatFieldValue(index, w.buf, value)
	if err != nil {
		w.err = fmt.Errorf("SetFloatFieldValue: %w", err)
	}
}
