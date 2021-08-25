package dbf

import (
	"bytes"
	"strings"
	"testing"
)

func Test_NewFields(t *testing.T) {
	f := NewFields()

	if f.Count() != 0 {
		t.Errorf("NewFields(): f.Count(): want: %v, got: %v", 0, f.Count())
	}
	if f.recSize != 1 {
		t.Errorf("NewFields(): f.recSize: want: %v, got: %v", 1, f.recSize)
	}
}

func Test_Fields_add_fields(t *testing.T) {
	f := NewFields()
	f.AddNumericField("price", 12, 2)
	f.AddLogicalField("flag")
	f.AddDateField("date")
	f.AddCharacterField("name", 25)

	if f.Count() != 4 {
		t.Errorf("add fields: f.Count(): want: %v, got: %v", 4, f.Count())
	}
	if f.recSize != 47 {
		t.Errorf("NewFields(): f.recSize: want: %v, got: %v", 47, f.recSize)
	}
}

func Test_Fields_FieldInfo(t *testing.T) {
	f := NewFields()
	f.AddDateField("date")
	name, typ, length, dec := f.FieldInfo(0)

	if name != "DATE" {
		t.Errorf("Fields.FieldInfo(): name: want: %#v, got: %#v", "DATE", name)
	}
	if typ != "D" {
		t.Errorf("Fields.FieldInfo(): typ: want: %#v, got: %#v", "D", typ)
	}
	if length != 8 {
		t.Errorf("Fields.FieldInfo(): length: want: %#v, got: %#v", 8, length)
	}
	if dec != 0 {
		t.Errorf("Fields.FieldInfo(): dec: want: %#v, got: %#v", 0, dec)
	}
}

func Test_Fields_write(t *testing.T) {
	f := NewFields()
	f.AddCharacterField("name", 14)
	f.AddDateField("date")

	buf := bytes.NewBuffer(nil)
	err := f.write(buf)

	if err != nil {
		t.Errorf("Fields.write(): %v", err)
	}
	if len(buf.Bytes()) != 64 {
		t.Errorf("Fields.write(): len(buf): want: %v, got: %v", 64, len(buf.Bytes()))
	}
}

func Test_Fields_read(t *testing.T) {
	b := make([]byte, fieldSize)
	copy(b[:], "NAME")
	b[11] = 'C'
	b[12] = 1
	b[16] = 14
	r := bytes.NewReader(b)

	f := NewFields()
	err := f.read(r, 1)

	if err != nil {
		t.Errorf("Fields.read(): %v", err)
	}
	if f.Count() != 1 {
		t.Errorf("Fields.read(): f.Count(): want: %v, got: %v", 1, f.Count())
	}
}

func Test_Fields_set_record_buffer(t *testing.T) {
	f := NewFields()
	f.AddCharacterField("name", 6)
	f.AddLogicalField("flag")
	f.AddNumericField("count", 4, 0)

	buf := []byte(strings.Repeat(" ", f.recSize))

	f.setStringFieldValue(0, buf, "Abc", nil)
	f.setBoolFieldValue(1, buf, true)
	f.setIntFieldValue(2, buf, 34)

	want := " Abc   T  34"

	if string(buf) != want {
		t.Errorf("Record buffer: want: %#v, got: %#v", want, string(buf))
	}

}

func Test_Fields_get_field_values(t *testing.T) {
	f := NewFields()
	f.AddCharacterField("name", 6)
	f.AddLogicalField("flag")
	f.AddNumericField("count", 4, 0)

	buf := []byte(" Abc   T  34")

	name, _ := f.stringFieldValue(0, buf, nil)
	flag, _ := f.boolFieldValue(1, buf)
	count, _ := f.intFieldValue(2, buf)

	if name != "Abc" {
		t.Errorf("Field value: want: %#v, got: %#v", "Abc", name)
	}
	if flag != true {
		t.Errorf("Field value: want: %#v, got: %#v", true, flag)
	}
	if count != 34 {
		t.Errorf("Field value: want: %#v, got: %#v", 34, count)
	}
}

func Test_Fields_check_duplicate(t *testing.T) {
	f := NewFields()
	f.AddCharacterField("flag", 6)
	f.AddLogicalField("flag")
	if f.err == nil {
		t.Errorf("Fields: add field duplicate: not error")
	}
}
