package dbf

import (
	"bytes"
	"reflect"
	"testing"
)

// New field

func Test_newLogicalField(t *testing.T) {
	f, _ := newLogicalField("Flag")

	tpl := "newLogicalField('Flag'): %s: want: %v, got: %v"

	if f.name() != "FLAG" {
		t.Errorf(tpl, "f.name()", "FLAG", f.name())
	}
	if f.Type != 'L' {
		t.Errorf(tpl, "f.Type", string('L'), string(f.Type))
	}
	if f.Len != 1 {
		t.Errorf(tpl, "f.Len", 1, f.Len)
	}
	if f.Dec != 0 {
		t.Errorf(tpl, "f.Dec", 0, f.Dec)
	}
}

func Test_newDateField(t *testing.T) {
	f, _ := newDateField("Date")

	tpl := "newDateField('Date'): %s: want: %v, got: %v"

	if f.name() != "DATE" {
		t.Errorf(tpl, "f.name()", "DATE", f.name())
	}
	if f.Type != 'D' {
		t.Errorf(tpl, "f.Type", string('D'), string(f.Type))
	}
	if f.Len != 8 {
		t.Errorf(tpl, "f.Len", 8, f.Len)
	}
	if f.Dec != 0 {
		t.Errorf(tpl, "f.Dec", 0, f.Dec)
	}
}

func Test_newCharacterField(t *testing.T) {
	f, _ := newCharacterField("Name", 25)

	tpl := "newCharacterField('Name', 25): %s: want: %v, got: %v"

	if f.name() != "NAME" {
		t.Errorf(tpl, "f.name()", "NAME", f.name())
	}
	if f.Type != 'C' {
		t.Errorf(tpl, "f.Type", string('C'), string(f.Type))
	}
	if f.Len != 25 {
		t.Errorf(tpl, "f.Len", 25, f.Len)
	}
	if f.Dec != 0 {
		t.Errorf(tpl, "f.Dec", 0, f.Dec)
	}
}

func Test_newNumericField(t *testing.T) {
	f, _ := newNumericField("Price", 12, 2)

	tpl := "newNumericField('Price', 12, 2): %s: want: %v, got: %v"

	if f.name() != "PRICE" {
		t.Errorf(tpl, "f.name()", "PRICE", f.name())
	}
	if f.Type != 'N' {
		t.Errorf(tpl, "f.Type", string('C'), string(f.Type))
	}
	if f.Len != 12 {
		t.Errorf(tpl, "f.Len", 25, f.Len)
	}
	if f.Dec != 2 {
		t.Errorf(tpl, "f.Dec", 0, f.Dec)
	}
}

// Field name

func Test_field_name(t *testing.T) {
	f := &field{
		Name: [11]byte{'N', 'A', 'M', 'E', 0, 0, 0, 0, 0, 0},
	}
	if f.name() != "NAME" {
		t.Errorf("field.name(): want: %#v, got: %#v", "NAME", f.name())
	}
}

func Test_field_setName(t *testing.T) {
	f := &field{}
	f.setName("name")

	if f.name() != "NAME" {
		t.Errorf("field.setName('name'): field.name(): want: %#v, got: %#v", "NAME", f.name())
	}
}

// Field read/write

var fieldBytes = []byte{'N', 'A', 'M', 'E', 0, 0, 0, 0, 0, 0, 0, 'C', 0, 0, 0, 0, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func Test_field_read(t *testing.T) {
	f := &field{}
	err := f.read(bytes.NewReader(fieldBytes))

	if err != nil {
		t.Errorf("field.read(): %v", err)
	}

	tpl := "field.read(): %s: want: %v, got: %v"

	if f.name() != "NAME" {
		t.Errorf(tpl, "f.name()", "NAME", f.name())
	}
	if f.Type != 'C' {
		t.Errorf(tpl, "f.Type", string('C'), string(f.Type))
	}
	if f.Len != 14 {
		t.Errorf(tpl, "f.Len", 14, f.Len)
	}
	if f.Dec != 0 {
		t.Errorf(tpl, "f.Dec", 0, f.Dec)
	}
}

func Test_field_write(t *testing.T) {
	f, _ := newCharacterField("name", 14)

	buf := bytes.NewBuffer(nil)
	err := f.write(buf)

	if err != nil {
		t.Errorf("field.write(): %v", err)
	}
	if !reflect.DeepEqual(buf.Bytes(), fieldBytes) {
		t.Errorf("field.write():\nwant: %#v\ngot : %#v", fieldBytes, buf.Bytes())
	}
}

// Value to string
/*
func Test_field_characterToString(t *testing.T) {
	f, _ := newCharacterField("name", 6)

	tests := []struct {
		value   interface{}
		encoder *encoding.Encoder
		wantRes string
		isErr   bool
	}{
		{value: "Abc", encoder: nil, wantRes: "Abc   ", isErr: false},
		{value: true, encoder: nil, wantRes: "", isErr: true},
	}
	for _, tc := range tests {
		gotRes, err := f.characterToString(tc.value, tc.encoder)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.characterToString(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.wantRes != gotRes {
			t.Errorf("field.characterToString(%#v): want: %#v, got: %#v", tc.value, tc.wantRes, gotRes)
		}
	}
}

func Test_logicalToString(t *testing.T) {
	f, _ := newLogicalField("name")

	tests := []struct {
		value   interface{}
		wantRes string
		isErr   bool
	}{
		{value: false, wantRes: "F", isErr: false},
		{value: true, wantRes: "T", isErr: false},
		{value: "abc", wantRes: "", isErr: true},
	}
	for _, tc := range tests {
		gotRes, err := f.logicalToString(tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.logicalToString(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.wantRes != gotRes {
			t.Errorf("field.logicalToString(%#v): want: %#v, got: %#v", tc.value, tc.wantRes, gotRes)
		}
	}
}

func Test_dateToString(t *testing.T) {
	f, _ := newDateField("name")
	d := time.Date(2021, 7, 26, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		value   interface{}
		wantRes string
		isErr   bool
	}{
		{value: d, wantRes: "20210726", isErr: false},
		{value: "abc", wantRes: "", isErr: true},
	}
	for _, tc := range tests {
		gotRes, err := f.dateToString(tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.dateToString(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.wantRes != gotRes {
			t.Errorf("field.dateToString(%#v): want: %#v, got: %#v", tc.value, tc.wantRes, gotRes)
		}
	}
}

func Test_numericToString(t *testing.T) {
	field1, _ := newNumericField("name", 6, 0)
	field2, _ := newNumericField("name", 9, 2)

	tests := []struct {
		field   *field
		value   interface{}
		wantRes string
		isErr   bool
	}{
		{field: field1, value: -123, wantRes: "  -123", isErr: false},
		{field: field1, value: "abc", wantRes: "", isErr: true},
		{field: field2, value: -123.4, wantRes: "  -123.40", isErr: false},
		{field: field2, value: 0, wantRes: "     0.00", isErr: false},
	}
	for _, tc := range tests {
		gotRes, err := tc.field.numericToString(tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.numericToString(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.wantRes != gotRes {
			t.Errorf("field.numericToString(%#v): want: %#v, got: %#v", tc.value, tc.wantRes, gotRes)
		}
	}
}
*/
// Bytes to value
/*
func Test_bytesToCharacter(t *testing.T) {
	f, _ := newCharacterField("name", 6)

	tests := []struct {
		value   []byte
		decoder *encoding.Decoder
		wantRes interface{}
		isErr   bool
	}{
		{value: []byte("Abc"), decoder: nil, wantRes: "Abc", isErr: false},
	}
	for _, tc := range tests {
		gotRes, err := f.bytesToCharacter(tc.value, tc.decoder)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.bytesToCharacter(%#v): want error: %v, got error: %v", string(tc.value), tc.isErr, gotErr)
		}
		if tc.wantRes != gotRes {
			t.Errorf("field.bytesToCharacter(%#v): want: %#v, got: %#v", string(tc.value), tc.wantRes, gotRes)
		}
	}
}

func Test_bytesToLogical(t *testing.T) {
	f, _ := newLogicalField("name")

	tests := []struct {
		value   []byte
		wantRes interface{}
	}{
		{value: []byte("T"), wantRes: true},
		{value: []byte("F"), wantRes: false},
	}
	for _, tc := range tests {
		gotRes := f.bytesToLogical(tc.value)

		if tc.wantRes != gotRes {
			t.Errorf("field.bytesToLogical(%#v): want: %#v, got: %#v", string(tc.value), tc.wantRes, gotRes)
		}
	}
}

func Test_bytesToDate(t *testing.T) {
	f, _ := newDateField("name")
	d := time.Date(2021, 7, 27, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		value   []byte
		wantRes interface{}
		isErr   bool
	}{
		{value: []byte("20210727"), wantRes: d, isErr: false},
	}
	for _, tc := range tests {
		gotRes, err := f.bytesToDate(tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.bytesToDate(%#v): want error: %v, got error: %v", string(tc.value), tc.isErr, gotErr)
		}
		if tc.wantRes != gotRes {
			t.Errorf("field.bytesToDate(%#v): want: %v, got: %v", string(tc.value), tc.wantRes, gotRes)
		}
	}
}

func TestBytesToNumericInt(t *testing.T) {
	field1, _ := newNumericField("name", 5, 0)
	field2, _ := newNumericField("name", 8, 2)

	tests := []struct {
		field   *field
		value   []byte
		wantRes interface{}
		isErr   bool
	}{
		{field: field1, value: []byte(" -123"), wantRes: int64(-123), isErr: false},
		{field: field2, value: []byte(" -123.45"), wantRes: float64(-123.45), isErr: false},
	}
	for _, tc := range tests {
		gotRes, err := tc.field.bytesToNumeric(tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.bytesToNumeric(%#v): want error: %v, got error: %v", string(tc.value), tc.isErr, gotErr)
		}
		if tc.wantRes != gotRes {
			t.Errorf("field.bytesToNumeric(%#v): want: %v, got: %v", string(tc.value), tc.wantRes, gotRes)
		}
	}
}
*/
