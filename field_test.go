package dbf

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"golang.org/x/text/encoding"
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

// Set field value

func Test_field_setStringFieldValue(t *testing.T) {
	f, _ := newCharacterField("name", 6)

	tests := []struct {
		value   string
		encoder *encoding.Encoder
		want    string
		isErr   bool
	}{
		{value: "Abc", encoder: nil, want: "Abc   ", isErr: false},
		{value: "", encoder: nil, want: "      ", isErr: false},
		{value: "1234567", encoder: nil, want: "      ", isErr: true},
	}
	for _, tc := range tests {
		buf := []byte("      ")
		err := f.setStringFieldValue(buf, tc.value, tc.encoder)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.setStringFieldValue(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.want != string(buf) {
			t.Errorf("field.setStringFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

func Test_field_setStringFieldValue_L(t *testing.T) {
	f, _ := newLogicalField("name")

	tests := []struct {
		value string
		want  string
	}{
		{value: "F", want: "F"},
		{value: "true", want: "T"},
		{value: "", want: " "},
		{value: "    ", want: " "},
	}
	for _, tc := range tests {
		buf := []byte("x")
		f.setStringFieldValue(buf, tc.value, nil)

		if tc.want != string(buf) {
			t.Errorf("field.setStringFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

func Test_field_setStringFieldValue_D(t *testing.T) {
	f, _ := newDateField("name")

	tests := []struct {
		value string
		want  string
	}{
		{value: "20210830", want: "20210830"},
		{value: "", want: "        "},
		{value: "   ", want: "        "},
		{value: "abc", want: "xxxxxxxx"},
	}
	for _, tc := range tests {
		buf := []byte("xxxxxxxx")
		f.setStringFieldValue(buf, tc.value, nil)

		if tc.want != string(buf) {
			t.Errorf("field.setStringFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

func Test_field_setStringFieldValue_N_int(t *testing.T) {
	f, _ := newNumericField("name", 5, 0)

	tests := []struct {
		value string
		want  string
	}{
		{value: " 123 ", want: "  123"},
		{value: "", want: "    0"},
		{value: "-12", want: "  -12"},
		{value: "12.34", want: "xxxxx"},
	}
	for _, tc := range tests {
		buf := []byte("xxxxx")
		f.setStringFieldValue(buf, tc.value, nil)

		if tc.want != string(buf) {
			t.Errorf("field.setStringFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

func Test_field_setStringFieldValue_N_float(t *testing.T) {
	f, _ := newNumericField("name", 8, 2)

	tests := []struct {
		value string
		want  string
	}{
		{value: " 123.4 ", want: "  123.40"},
		{value: "", want: "    0.00"},
		{value: "-12", want: "  -12.00"},
		{value: "abc", want: "xxxxxxxx"},
	}
	for _, tc := range tests {
		buf := []byte("xxxxxxxx")
		f.setStringFieldValue(buf, tc.value, nil)

		if tc.want != string(buf) {
			t.Errorf("field.setStringFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

func Test_field_setBoolFieldValue(t *testing.T) {
	f, _ := newLogicalField("name")

	tests := []struct {
		value bool
		want  string
		isErr bool
	}{
		{value: false, want: "F", isErr: false},
		{value: true, want: "T", isErr: false},
	}
	for _, tc := range tests {
		buf := []byte(" ")
		err := f.setBoolFieldValue(buf, tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.setBoolFieldValue(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.want != string(buf) {
			t.Errorf("field.setBoolFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

func Test_field_setDateFieldValue(t *testing.T) {
	f, _ := newDateField("name")
	d := time.Date(2021, 7, 26, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		value time.Time
		want  string
		isErr bool
	}{
		{value: d, want: "20210726", isErr: false},
	}
	for _, tc := range tests {
		buf := []byte("        ")
		err := f.setDateFieldValue(buf, tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.setDateFieldValue(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.want != string(buf) {
			t.Errorf("field.setDateFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

func Test_field_setIntFieldValue(t *testing.T) {
	f, _ := newNumericField("name", 6, 0)

	tests := []struct {
		value int64
		want  string
		isErr bool
	}{
		{value: -123, want: "  -123", isErr: false},
		{value: 123, want: "   123", isErr: false},
		{value: 0, want: "     0", isErr: false},
		{value: 1234567, want: "      ", isErr: true},
	}
	for _, tc := range tests {
		buf := []byte("      ")
		err := f.setIntFieldValue(buf, tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.setIntFieldValue(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.want != string(buf) {
			t.Errorf("field.setIntFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

func Test_field_setFloatFieldValue(t *testing.T) {
	f, _ := newNumericField("name", 9, 2)

	tests := []struct {
		value float64
		want  string
		isErr bool
	}{
		{value: -123.45, want: "  -123.45", isErr: false},
		{value: 123, want: "   123.00", isErr: false},
		{value: 0, want: "     0.00", isErr: false},
		{value: 1234567, want: "         ", isErr: true},
	}
	for _, tc := range tests {
		buf := []byte("         ")
		err := f.setFloatFieldValue(buf, tc.value)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.setFloatFieldValue(%#v): want error: %v, got error: %v", tc.value, tc.isErr, gotErr)
		}
		if tc.want != string(buf) {
			t.Errorf("field.setFloatFieldValue(%#v): want: %#v, got: %#v", tc.value, tc.want, string(buf))
		}
	}
}

// Get field value

func Test_field_stringFieldValue(t *testing.T) {
	f, _ := newCharacterField("name", 6)

	tests := []struct {
		buf     []byte
		decoder *encoding.Decoder
		want    string
		isErr   bool
	}{
		{buf: []byte("Abc   "), decoder: nil, want: "Abc", isErr: false},
		{buf: []byte(" Abc  "), decoder: nil, want: " Abc", isErr: false},
		{buf: []byte("      "), decoder: nil, want: "", isErr: false},
	}
	for _, tc := range tests {
		got, err := f.stringFieldValue(tc.buf, tc.decoder)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.stringFieldValue(%#v): want error: %v, got error: %v", string(tc.buf), tc.isErr, gotErr)
		}
		if tc.want != got {
			t.Errorf("field.stringFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_stringFieldValue_L(t *testing.T) {
	f, _ := newLogicalField("name")

	tests := []struct {
		buf  []byte
		want string
	}{
		{buf: []byte("T"), want: "T"},
		{buf: []byte("F"), want: "F"},
		{buf: []byte(" "), want: ""},
	}
	for _, tc := range tests {
		got, _ := f.stringFieldValue(tc.buf, nil)
		if tc.want != got {
			t.Errorf("field.stringFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_stringFieldValue_D(t *testing.T) {
	f, _ := newDateField("name")

	tests := []struct {
		buf  []byte
		want string
	}{
		{buf: []byte("20210830"), want: "20210830"},
		{buf: []byte("        "), want: ""},
	}
	for _, tc := range tests {
		got, _ := f.stringFieldValue(tc.buf, nil)
		if tc.want != got {
			t.Errorf("field.stringFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_stringFieldValue_N(t *testing.T) {
	f, _ := newNumericField("name", 8, 2)

	tests := []struct {
		buf  []byte
		want string
	}{
		{buf: []byte("  123.45"), want: "123.45"},
		{buf: []byte(" -123.00"), want: "-123.00"},
		{buf: []byte("    0.00"), want: "0.00"},
		{buf: []byte("        "), want: ""},
	}
	for _, tc := range tests {
		got, _ := f.stringFieldValue(tc.buf, nil)
		if tc.want != got {
			t.Errorf("field.stringFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_boolFieldValue(t *testing.T) {
	f, _ := newLogicalField("name")

	tests := []struct {
		buf   []byte
		want  bool
		isErr bool
	}{
		{buf: []byte("T"), want: true, isErr: false},
		{buf: []byte("F"), want: false, isErr: false},
		{buf: []byte(" "), want: false, isErr: false},
	}
	for _, tc := range tests {
		got, err := f.boolFieldValue(tc.buf)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.boolFieldValue(%#v): want error: %v, got error: %v", string(tc.buf), tc.isErr, gotErr)
		}
		if tc.want != got {
			t.Errorf("field.boolFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_dateFieldValue(t *testing.T) {
	f, _ := newDateField("name")

	d1 := time.Date(2021, 7, 27, 0, 0, 0, 0, time.UTC)
	d2 := time.Time{}

	tests := []struct {
		buf   []byte
		want  time.Time
		isErr bool
	}{
		{buf: []byte("20210727"), want: d1, isErr: false},
		{buf: []byte("        "), want: d2, isErr: false},
	}
	for _, tc := range tests {
		got, err := f.dateFieldValue(tc.buf)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.dateFieldValue(%#v): want error: %v, got error: %v", string(tc.buf), tc.isErr, gotErr)
		}
		if tc.want != got {
			t.Errorf("field.dateFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_intFieldValue(t *testing.T) {
	f, _ := newNumericField("name", 5, 0)

	tests := []struct {
		buf   []byte
		want  int64
		isErr bool
	}{
		{buf: []byte(" -123"), want: -123, isErr: false},
		{buf: []byte("  123"), want: 123, isErr: false},
		{buf: []byte("     "), want: 0, isErr: false},
		{buf: []byte("abc  "), want: 0, isErr: true},
	}
	for _, tc := range tests {
		got, err := f.intFieldValue(tc.buf)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.intFieldValue(%#v): want error: %v, got error: %v", string(tc.buf), tc.isErr, gotErr)
		}
		if tc.want != got {
			t.Errorf("field.intFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_intFieldValue_dec_not_zero(t *testing.T) {
	f, _ := newNumericField("name", 8, 2)

	tests := []struct {
		buf   []byte
		want  int64
		isErr bool
	}{
		{buf: []byte(" -123.45"), want: -123, isErr: false},
		{buf: []byte("  123.67"), want: 123, isErr: false},
		{buf: []byte("        "), want: 0, isErr: false},
	}
	for _, tc := range tests {
		got, err := f.intFieldValue(tc.buf)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.intFieldValue(%#v): want error: %v, got error: %v", string(tc.buf), tc.isErr, gotErr)
		}
		if tc.want != got {
			t.Errorf("field.intFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_floatFieldValue(t *testing.T) {
	f, _ := newNumericField("name", 8, 2)

	tests := []struct {
		buf   []byte
		want  float64
		isErr bool
	}{
		{buf: []byte(" -123.45"), want: -123.45, isErr: false},
		{buf: []byte("  123.00"), want: 123, isErr: false},
		{buf: []byte("        "), want: 0, isErr: false},
		{buf: []byte("abc     "), want: 0, isErr: true},
	}
	for _, tc := range tests {
		got, err := f.floatFieldValue(tc.buf)
		gotErr := (err != nil)

		if tc.isErr != gotErr {
			t.Errorf("field.floatFieldValue(%#v): want error: %v, got error: %v", string(tc.buf), tc.isErr, gotErr)
		}
		if tc.want != got {
			t.Errorf("field.floatFieldValue(%#v): want: %#v, got: %#v", string(tc.buf), tc.want, got)
		}
	}
}

func Test_field_check_type(t *testing.T) {
	f, _ := newNumericField("name", 8, 2)

	if _, err := f.dateFieldValue(nil); err == nil {
		t.Errorf("field check type: error requered")
	}
}
