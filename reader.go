package dbf

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding"
)

// The Reader reads records from a CSV file.
type Reader struct {
	header  *header
	fields  []*field
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
func NewReader(r io.Reader) *Reader {
	return &Reader{
		header: &header{},
		reader: bufio.NewReader(r),
	}
}

// Read reads one record (a slice of fields) from r.
// Read always returns either a non-nil record or a non-nil error, but not both.
// If there is no data left to be read, Read returns nil, io.EOF.
// If ReuseRecord is true, the returned slice may be shared between multiple calls to Read.
func (r *Reader) Read() (record []interface{}, err error) {
	if len(r.fields) == 0 {
		if err := r.initReader(); err != nil {
			return nil, err
		}
	}
	if r.ReuseRecord {
		record, err = r.readRecord(r.lastRecord)
		r.lastRecord = record
	} else {
		record, err = r.readRecord(nil)
	}
	return record, err
}

func (r *Reader) initReader() error {
	if err := r.header.read(r.reader); err != nil {
		return err
	}
	if err := r.readFields(); err != nil {
		return err
	}
	// Skip byte header end
	if _, err := r.reader.Discard(1); err != nil {
		return err
	}
	// Create buffer
	r.buf = make([]byte, int(r.header.RecSize))
	// Code page
	cm := charmapByPage(r.header.codePage())
	if cm != nil {
		r.decoder = cm.NewDecoder()
	}
	return nil
}

func (r *Reader) readFields() error {
	offset := 1 // deleted mark
	count := r.header.fieldCount()

	for i := 0; i < count; i++ {
		f := &field{}
		if err := f.read(r.reader); err != nil {
			return err
		}
		f.Offset = uint32(offset)
		r.fields = append(r.fields, f)
		offset += int(f.Len)
	}
	return nil
}

func (r *Reader) readRecord(dst []interface{}) ([]interface{}, error) {
	r.recNo++
	if _, err := io.ReadFull(r.reader, r.buf); err != nil {
		if err == io.ErrUnexpectedEOF {
			err = io.EOF
		}
		return nil, err
	}
	if len(dst) != len(r.fields) {
		dst = make([]interface{}, len(r.fields))
	}
	var err error
	for i := range r.fields {
		dst[i], err = r.fieldValue(i)
		if err != nil {
			return nil, fmt.Errorf("record %d: field %q: %w", r.recNo, r.fields[i].name(), err)
		}
	}
	return dst, nil
}

func (r *Reader) fieldValue(index int) (interface{}, error) {
	var result interface{}
	var err error
	f := r.fields[index]
	s := string(r.buf[int(f.Offset) : int(f.Offset)+int(f.Len)])

	switch f.Type {
	case 'C':
		s = strings.TrimRight(s, " ")
		if r.decoder != nil && !isASCII(s) {
			s, err = r.decoder.String(s)
			if err != nil {
				return result, err
			}
		}
		result = s
	case 'L':
		b := s[0]
		result = (b == 'T' || b == 't' || b == 'Y' || b == 'y')
	case 'D':
		var d time.Time
		if strings.Trim(s, " ") == "" {
			result = d
		} else {
			d, err = time.Parse("20060102", s)
			if err != nil {
				return result, err
			}
			result = d
		}
	case 'N':
		s = strings.TrimSpace(s)
		if s == "" || s == "." {
			s = "0"
		}
		if f.Dec == 0 {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return result, err
			}
			result = n
		} else {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return result, err
			}
			result = n
		}
	default:
		return result, fmt.Errorf("invalid field type: got %s, want C, N, L, D", string(f.Type))
	}
	return result, nil
}
