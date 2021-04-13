// xbase package implements functions for working with DBF files.
package xbase

import (
	"fmt"
	"os"

	"golang.org/x/text/encoding"
)

const (
	dbfId     byte = 0x03
	headerEnd byte = 0x0D
	fileEnd   byte = 0x1A
)

const (
	fieldSize  = 32
	headerSize = 32
)

type XBase struct {
	header  *header
	fields  []*field
	file    *os.File
	buf     []byte
	recNo   int64
	isAdd   bool
	isMod   bool
	encoder *encoding.Encoder
	decoder *encoding.Decoder
}

// Public

// New creates a XBase object to work with a DBF file.
func New() *XBase {
	return &XBase{header: newHeader()}
}

// RecNo returns the sequence number of the current record.
// Numbering starts from 1.
func (db *XBase) RecNo() int64 {
	return db.recNo
}

// RecCount returns the number of records in the DBF file.
func (db *XBase) RecCount() int64 {
	return int64(db.header.RecCount)
}

// FieldCount returns the number of fields in the DBF file.
func (db *XBase) FieldCount() int {
	return len(db.fields)
}

// EOF returns true if end of file is reached.
func (db *XBase) EOF() bool {
	return db.recNo > db.RecCount() || db.RecCount() == 0
}

// BOF returns true if the beginning of the file is reached.
func (db *XBase) BOF() bool {
	return db.recNo == 0 || db.RecCount() == 0
}

// AddField adds a field to the structure of the DBF file.
// This method can only be used before creating a new file.
//
// The following field types are supported: "C", "N", "L", "D".
//
// The opts parameter contains optional parameters: field length and number of decimal places.
//
// Examples:
//     db.AddField("NAME", "C", 24)
//     db.AddField("COUNT", "N", 8)
//     db.AddField("PRICE", "N", 12, 2)
//     db.AddField("FLAG", "L")
//     db.AddField("DATE", "D")
func (db *XBase) AddField(name string, typ string, opts ...int) error {
	length := 0
	dec := 0
	if len(opts) > 0 {
		length = opts[0]
	}
	if len(opts) > 1 {
		dec = opts[1]
	}
	f, err := newField(name, typ, length, dec)
	if err != nil {
		return wrapError(fmt.Sprintf("AddField: field %q", name), err)
	}
	db.fields = append(db.fields, f)
	return nil
}

func wrapError(s string, err error) error {
	return fmt.Errorf("xbase: %s: %w", s, err)
}
