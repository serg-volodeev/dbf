// xbase package implements functions for working with DBF files.
package xbase

import (
	"fmt"
	"os"
	"strings"
	"time"

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

// FieldNo returns the number of the field by name.
// If name is not found returns 0.
// Fields are numbered starting from 1.
func (db *XBase) FieldNo(name string) int {
	name = strings.ToUpper(strings.TrimSpace(name))
	for i, f := range db.fields {
		if f.name() == name {
			return i + 1
		}
	}
	return 0
}

// FieldInfo returns field attributes by number.
// Fields are numbered starting from 1.
func (db *XBase) FieldInfo(fieldNo int) (name, typ string, length, dec int, err error) {
	f, err := db.fieldByNo(fieldNo)
	if err != nil {
		err = db.wrapFieldError("FieldInfo", fieldNo, err)
		return
	}
	name = f.name()
	typ = string([]byte{f.Type})
	length = int(f.Len)
	dec = int(f.Dec)
	return
}

func (db *XBase) fieldByNo(fieldNo int) (*field, error) {
	if fieldNo < 1 || fieldNo > len(db.fields) {
		return nil, fmt.Errorf("field number out of range")
	}
	return db.fields[fieldNo-1], nil
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

func (db *XBase) wrapFieldError(s string, fieldNo int, err error) error {
	prefix := fmt.Sprintf("xbase: %s: field %d", s, fieldNo)
	if fieldNo < 1 || fieldNo > len(db.fields) {
		return fmt.Errorf("%s: %w", prefix, err)
	}
	return fmt.Errorf("%s %q: %w", prefix, db.fields[fieldNo-1].name(), err)
}

// SetCodePage sets the encoding mode for reading and writing string field values.
// The default code page is 0.
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
func (db *XBase) SetCodePage(cp int) {
	cm := charmapByPage(cp)
	if cm == nil {
		return
	}
	db.encoder = cm.NewEncoder()
	db.decoder = cm.NewDecoder()
	db.header.setCodePage(cp)
}

// CodePage returns the code page of a DBF file.
// Returns 0 if no code page is specified.
func (db *XBase) CodePage() int {
	return db.header.codePage()
}

// OpenFile opens an existing DBF file.
func (db *XBase) OpenFile(name string, readOnly bool) error {
	var err error

	// Open file
	if readOnly {
		db.file, err = os.Open(name)
	} else {
		db.file, err = os.OpenFile(name, os.O_RDWR, 0666)
	}
	if err != nil {
		return wrapError("OpenFile", err)
	}

	// Read header
	err = db.header.read(db.file)
	if err != nil {
		return wrapError("OpenFile", err)
	}

	// Read fields
	err = db.readFields()
	if err != nil {
		return wrapError("OpenFile", err)
	}

	db.buf = make([]byte, int(db.header.RecSize))
	db.SetCodePage(db.header.codePage())
	return nil
}

func (db *XBase) readFields() error {
	offset := 1 // deleted mark
	count := db.header.fieldCount()

	for i := 0; i < count; i++ {
		f := &field{}
		if err := f.read(db.file); err != nil {
			return err
		}
		f.Offset = uint32(offset)
		db.fields = append(db.fields, f)
		offset += int(f.Len)
	}
	return nil
}

// CloseFile closes a previously opened or created DBF file.
func (db *XBase) CloseFile() error {
	if err := db.file.Close(); err != nil {
		return wrapError("CloseFile", err)
	}
	return nil

	// if db.err != nil {
	// 	return
	// }
	// defer db.wrapError("CloseFile")
	// if db.isMod {
	// 	db.header.setModDate(time.Now())
	// 	db.writeHeader()
	// 	db.writeFileEnd()
	// }
	// db.fileClose()
}

func (db *XBase) goTo(recNo int64) error {
	if recNo < 1 {
		db.recNo = 0
		return nil
	}
	if recNo > db.RecCount() {
		db.recNo = db.RecCount() + 1
		return nil
	}
	db.recNo = recNo
	// Seek
	offset := int64(db.header.DataOffset) + int64(db.header.RecSize)*(db.recNo-1)
	if _, err := db.file.Seek(offset, 0); err != nil {
		return err
	}
	// Read
	if _, err := db.file.Read(db.buf); err != nil {
		return err
	}
	return nil
}

// First positions the object to the first record.
func (db *XBase) First() error {
	if err := db.goTo(1); err != nil {
		return wrapError("First", err)
	}
	return nil
}

// Last positions the object to the last record.
func (db *XBase) Last() error {
	if err := db.goTo(db.RecCount()); err != nil {
		return wrapError("Last", err)
	}
	return nil
}

// Next positions the object to the next record.
func (db *XBase) Next() error {
	if err := db.goTo(db.recNo + 1); err != nil {
		return wrapError("Next", err)
	}
	return nil
}

// Prev positions the object to the previous record.
func (db *XBase) Prev() error {
	if err := db.goTo(db.recNo - 1); err != nil {
		return wrapError("Prev", err)
	}
	return nil
}

// FieldValueAsString returns the string value of the field of the current record.
// Fields are numbered starting from 1.
func (db *XBase) FieldValueAsString(fieldNo int) (string, error) {
	f, err := db.fieldByNo(fieldNo)
	if err != nil {
		return "", db.wrapFieldError("FieldValueAsString", fieldNo, err)
	}
	v, err := f.value(db.buf, db.decoder)
	if err != nil {
		return "", db.wrapFieldError("FieldValueAsString", fieldNo, err)
	}
	return v.(string), nil
}

// FieldValueAsInt returns the integer value of the field of the current record.
// Field type must be numeric ("N"). Fields are numbered starting from 1.
func (db *XBase) FieldValueAsInt(fieldNo int) (int64, error) {
	f, err := db.fieldByNo(fieldNo)
	if err != nil {
		return 0, db.wrapFieldError("FieldValueAsInt", fieldNo, err)
	}
	v, err := f.value(db.buf, nil)
	if err != nil {
		return 0, db.wrapFieldError("FieldValueAsInt", fieldNo, err)
	}
	return v.(int64), nil
}

// FieldValueAsFloat returns the float value of the field of the current record.
// Field type must be numeric ("N"). Fields are numbered starting from 1.
func (db *XBase) FieldValueAsFloat(fieldNo int) (float64, error) {
	f, err := db.fieldByNo(fieldNo)
	if err != nil {
		return 0, db.wrapFieldError("FieldValueAsFloat", fieldNo, err)
	}
	v, err := f.value(db.buf, nil)
	if err != nil {
		return 0, db.wrapFieldError("FieldValueAsFloat", fieldNo, err)
	}
	return v.(float64), nil
}

// FieldValueAsBool returns the boolean value of the field of the current record.
// Field type must be logical ("L"). Fields are numbered starting from 1.
func (db *XBase) FieldValueAsBool(fieldNo int) (bool, error) {
	f, err := db.fieldByNo(fieldNo)
	if err != nil {
		return false, db.wrapFieldError("FieldValueAsBool", fieldNo, err)
	}
	v, err := f.value(db.buf, nil)
	if err != nil {
		return false, db.wrapFieldError("FieldValueAsBool", fieldNo, err)
	}
	return v.(bool), nil
}

// FieldValueAsDate returns the date value of the field of the current record.
// Field type must be date ("D"). Fields are numbered starting from 1.
func (db *XBase) FieldValueAsDate(fieldNo int) (time.Time, error) {
	var d time.Time
	f, err := db.fieldByNo(fieldNo)
	if err != nil {
		return d, db.wrapFieldError("FieldValueAsDate", fieldNo, err)
	}
	v, err := f.value(db.buf, nil)
	if err != nil {
		return d, db.wrapFieldError("FieldValueAsDate", fieldNo, err)
	}
	return v.(time.Time), nil
}
