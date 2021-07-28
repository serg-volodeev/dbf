package dbf

import (
	"bufio"
	"fmt"
	"io"

	"golang.org/x/text/encoding"
)

// The Reader reads records from a CSV file.
type Reader struct {
	header  *header
	fields  *Fields
	reader  *bufio.Reader
	buf     []byte
	recNo   uint32
	decoder *encoding.Decoder

	// ReuseRecord controls whether calls to Read may return a slice sharing
	// the backing array of the previous call's returned slice for performance.
	// By default, each call to Read returns newly allocated memory owned by the caller.
	ReuseRecord bool

	// lastRecord is a record cache and only used when ReuseRecord == true.
	lastRecord []interface{}
}

// NewReader returns a new Reader that reads from r.
func NewReader(rd io.Reader) (*Reader, error) {
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
func (r *Reader) SetCodePage(cp int) error {
	cm := charmapByPage(cp)
	if cm == nil {
		return fmt.Errorf("unsupported code page %d", cp)
	}
	r.decoder = cm.NewDecoder()
	r.header.setCodePage(cp)
	return nil
}

// CodePage returns the code page set in the file header.
func (r *Reader) CodePage() int {
	return r.header.codePage()
}

// RecordCount returns the number of records in the DBF file.
func (r *Reader) RecordCount() uint32 {
	return r.header.RecCount
}

// Fields returns the file structure.
func (r *Reader) Fields() *Fields {
	return r.fields
}

// Read reads one record (a slice of fields) from r.
// Read always returns either a non-nil record or a non-nil error, but not both.
// If there is no data left to be read, Read returns nil, io.EOF.
// If ReuseRecord is true, the returned slice may be shared between multiple calls to Read.
func (r *Reader) Read() (record []interface{}, err error) {
	if r.ReuseRecord {
		record, err = r.readRecord(r.lastRecord)
		r.lastRecord = record
	} else {
		record, err = r.readRecord(nil)
	}
	return record, err
}

func (r *Reader) readRecord(dst []interface{}) ([]interface{}, error) {
	r.recNo++
	if _, err := io.ReadFull(r.reader, r.buf); err != nil {
		if err == io.ErrUnexpectedEOF {
			err = io.EOF
		}
		return nil, err
	}
	if len(dst) != r.fields.Count() {
		dst = make([]interface{}, r.fields.Count())
	}
	var err error
	for i := range r.fields.items {
		dst[i], err = r.fieldValue(i)
		if err != nil {
			return nil, fmt.Errorf("record %d: field %q: %w", r.recNo, r.fields.items[i].name(), err)
		}
	}
	return dst, nil
}

func (r *Reader) fieldValue(index int) (interface{}, error) {
	var result interface{}
	var err error

	f := r.fields.items[index]
	buf := r.buf[int(f.Offset) : int(f.Offset)+int(f.Len)]

	switch f.Type {
	case 'C':
		result, err = f.bytesToCharacter(buf, r.decoder)
	case 'L':
		result = f.bytesToLogical(buf)
	case 'D':
		result, err = f.bytesToDate(buf)
	case 'N':
		result, err = f.bytesToNumeric(buf)
	default:
		return result, fmt.Errorf("invalid field type: got %s, want C, N, L, D", string(f.Type))
	}
	if err != nil {
		return result, err
	}
	return result, nil
}
