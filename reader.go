package dbf

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"golang.org/x/text/encoding"
)

// The Reader reads records from a DBF file.
type Reader struct {
	header  *header
	fields  *Fields
	reader  *bufio.Reader
	buf     []byte
	recNo   uint32
	decoder *encoding.Decoder
	err     error
}

// NewReader returns a new Reader that reads from r.
func NewReader(rd io.Reader) (*Reader, error) {
	r, err := newReader(rd)
	if err != nil {
		return nil, fmt.Errorf("dbf.NewReader: %w", err)
	}
	return r, nil
}

func newReader(rd io.Reader) (*Reader, error) {
	if _, ok := rd.(io.Reader); !ok {
		return nil, fmt.Errorf("parameter %v is not io.Reader", rd)
	}
	r := &Reader{
		header: &header{},
		fields: NewFields(),
		reader: bufio.NewReader(rd),
	}
	if err := r.header.read(r.reader); err != nil {
		return nil, err
	}
	if err := r.fields.read(r.reader, r.header.fieldCount()); err != nil {
		return nil, err
	}
	// Skip byte header end
	if _, err := r.reader.Discard(1); err != nil {
		return nil, err
	}
	// Create buffer
	r.buf = make([]byte, int(r.header.RecSize))
	// Code page
	if cp := r.header.codePage(); cp != 0 {
		cm := charmapByPage(cp)
		r.decoder = cm.NewDecoder()
	}
	return r, nil
}

// Err returns the first error that was encountered by the Reader.
func (r *Reader) Err() error {
	if r.err != nil {
		return fmt.Errorf("dbf.Reader: %w", r.err)
	}
	return nil
}

// SetCodePage sets the code page if no code page is set in the file header.
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
func (r *Reader) SetCodePage(cp int) {
	if r.err != nil {
		return
	}
	cm := charmapByPage(cp)
	if cm == nil {
		r.err = fmt.Errorf("SetCodePage: unsupported code page %d", cp)
		return
	}
	r.decoder = cm.NewDecoder()
	r.header.setCodePage(cp)
}

// CodePage returns the code page set in the file header.
func (r *Reader) CodePage() int {
	if r.err != nil {
		return 0
	}
	return r.header.codePage()
}

// ModDate returns the modified date in the file header.
func (r *Reader) ModDate() time.Time {
	if r.err != nil {
		return time.Time{}
	}
	return r.header.modDate()
}

// RecordCount returns the number of records in the DBF file.
func (r *Reader) RecordCount() uint32 {
	if r.err != nil {
		return 0
	}
	return r.header.RecCount
}

// Fields returns the file structure.
func (r *Reader) Fields() *Fields {
	if r.err != nil {
		return nil
	}
	return r.fields
}

// Read reads one record from r.
// Returns false if end of file is reached or an error occurs.
func (r *Reader) Read() bool {
	if r.err != nil {
		return false
	}
	r.recNo++
	if _, err := io.ReadFull(r.reader, r.buf); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			r.err = fmt.Errorf("Read: record %d: %w", r.recNo, err)
		}
		return false
	}
	return true
}

// Deleted returns deleted record flag.
func (r *Reader) Deleted() bool {
	if r.err != nil {
		return false
	}
	return r.buf[0] == '*'
}

// StringFieldValue returns the value of the field by index.
// Field type must be Character, Date, Logical or Numeric.
func (r *Reader) StringFieldValue(index int) string {
	if r.err != nil {
		return ""
	}
	value, err := r.fields.stringFieldValue(index, r.buf, r.decoder)
	if err != nil {
		r.err = fmt.Errorf("StringFieldValue: %w", err)
	}
	return value
}

// BoolFieldValue returns the value of the field by index.
// Field type must be Logical.
func (r *Reader) BoolFieldValue(index int) bool {
	if r.err != nil {
		return false
	}
	value, err := r.fields.boolFieldValue(index, r.buf)
	if err != nil {
		r.err = fmt.Errorf("BoolFieldValue: %w", err)
	}
	return value
}

// DateFieldValue returns the value of the field by index.
// Field type must be Date.
func (r *Reader) DateFieldValue(index int) time.Time {
	if r.err != nil {
		return time.Time{}
	}
	value, err := r.fields.dateFieldValue(index, r.buf)
	if err != nil {
		r.err = fmt.Errorf("DateFieldValue: %w", err)
	}
	return value
}

// IntFieldValue returns the value of the field by index.
// Field type must be Numeric.
// If field decimal places is not zero,
// then it returns the integer part of the number.
func (r *Reader) IntFieldValue(index int) int64 {
	if r.err != nil {
		return 0
	}
	value, err := r.fields.intFieldValue(index, r.buf)
	if err != nil {
		r.err = fmt.Errorf("IntFieldValue: %w", err)
	}
	return value
}

// FloatFieldValue returns the value of the field by index.
// Field type must be Numeric.
func (r *Reader) FloatFieldValue(index int) float64 {
	if r.err != nil {
		return 0
	}
	value, err := r.fields.floatFieldValue(index, r.buf)
	if err != nil {
		r.err = fmt.Errorf("FloatFieldValue: %w", err)
	}
	return value
}
