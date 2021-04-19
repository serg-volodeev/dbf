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

type Reader struct {
	header  *header
	fields  []*field
	reader  *bufio.Reader
	buf     []byte
	recNo   uint32
	decoder *encoding.Decoder
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		header: &header{},
		reader: bufio.NewReader(r),
	}
}

func (r *Reader) Read() ([]interface{}, error) {
	if len(r.fields) == 0 {
		if err := r.readHeader(); err != nil {
			return nil, err
		}
	}
	return r.readRecord()
}

func (r *Reader) readHeader() error {
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

func (r *Reader) readRecord() ([]interface{}, error) {
	r.recNo++
	if _, err := io.ReadFull(r.reader, r.buf); err != nil {
		if err == io.ErrUnexpectedEOF {
			err = io.EOF
		}
		return nil, err
	}
	record := make([]interface{}, len(r.fields))
	var err error
	for i := range r.fields {
		record[i], err = r.fieldValue(i)
		if err != nil {
			return nil, err
		}
	}
	return record, nil
}

func (r *Reader) fieldValue(index int) (interface{}, error) {
	var result interface{}
	var err error
	f := r.fields[index]
	fieldBuf := r.buf[int(f.Offset) : int(f.Offset)+int(f.Len)]
	fieldStr := string(fieldBuf)

	switch f.Type {
	case 'C':
		s := strings.TrimRight(fieldStr, " ")
		if r.decoder != nil && !isASCII(s) {
			s, err = r.decoder.String(s)
			if err != nil {
				return result, err
			}
		}
		result = s
	case 'L':
		b := fieldBuf[0]
		result = (b == 'T' || b == 't' || b == 'Y' || b == 'y')
	case 'D':
		var d time.Time
		s := fieldStr
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
		s := strings.TrimSpace(fieldStr)
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
