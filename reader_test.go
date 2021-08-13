package dbf

import (
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_NewReader_nil(t *testing.T) {
	r, err := NewReader(nil)
	if err == nil {
		t.Errorf("NewReader(nil): require error")
	}
	if r != nil {
		t.Errorf("NewReader(nil): require nil reader")
	}
}

func Test_NewReader(t *testing.T) {
	fname := "./testdata/rec3.dbf"
	f, err := os.Open(fname)
	if err != nil {
		t.Errorf("os.Open(%q): %v", fname, err)
	}
	defer f.Close()

	r, err := NewReader(f)

	if err != nil {
		t.Errorf("NewReader(): %v", err)
	}
	if r.RecordCount() != 3 {
		t.Errorf("NewReader(): r.RecordCount(): want: %v, got: %v", 3, r.RecordCount())
	}
	if r.CodePage() != 866 {
		t.Errorf("NewReader(): r.CodePage(): want: %v, got: %v", 866, r.CodePage())
	}

	testFields := []struct {
		Name, Type string
		Len, Dec   int
	}{
		{"NAME", "C", 20, 0},
		{"FLAG", "L", 1, 0},
		{"COUNT", "N", 5, 0},
		{"PRICE", "N", 9, 2},
		{"DATE", "D", 8, 0},
	}
	fields := r.Fields()

	if fields.Count() != len(testFields) {
		t.Errorf("NewReader(): fields.Count(): want %v, got %v", len(testFields), fields.Count())
	}

	for i, f := range testFields {
		name, typ, length, dec := fields.FieldInfo(i)

		if f.Name != name {
			t.Errorf("NewReader(): fields.FieldInfo(%d): name: want: %#v, got: %#v", i, f.Name, name)
		}
		if f.Type != typ {
			t.Errorf("NewReader(): fields.FieldInfo(%d): typ: want: %#v, got: %#v", i, f.Type, typ)
		}
		if f.Len != length {
			t.Errorf("NewReader(): fields.FieldInfo(%d): length: want: %#v, got: %#v", i, f.Len, length)
		}
		if f.Dec != dec {
			t.Errorf("NewReader(): fields.FieldInfo(%d): dec: want: %#v, got: %#v", i, f.Dec, dec)
		}
	}
}

func Test_Reader_Read_file_empty(t *testing.T) {
	fname := "./testdata/rec0.dbf"
	f, err := os.Open(fname)
	if err != nil {
		t.Errorf("os.Open(%q): %v", fname, err)
	}
	defer f.Close()

	r, err := NewReader(f)
	if err != nil {
		t.Errorf("NewReader(): %v", err)
	}
	if r.RecordCount() != 0 {
		t.Errorf("NewReader(): r.RecordCount(): want: %v, got: %v", 0, r.RecordCount())
	}

	rec, err := r.Read()

	if err != io.EOF {
		t.Errorf("Read(): err: want: %v, got: %v", io.EOF, err)
	}
	if len(rec) != 0 {
		t.Errorf("Read(): len(rec): want: %v, got: %v", 0, len(rec))
	}
}

func Test_Reader_Read_record_blank(t *testing.T) {
	fname := "./testdata/rec1.dbf"
	f, err := os.Open(fname)
	if err != nil {
		t.Errorf("os.Open(%q): %v", fname, err)
	}
	defer f.Close()

	r, err := NewReader(f)
	if err != nil {
		t.Errorf("NewReader(): %v", err)
	}
	if r.RecordCount() != 1 {
		t.Errorf("NewReader(): r.RecordCount(): want: %v, got: %v", 1, r.RecordCount())
	}

	rec, err := r.Read()

	if err != nil {
		t.Errorf("Read(): %v", err)
	}
	if len(rec) != 5 {
		t.Errorf("Read(): len(rec): want: %v, got: %v", 0, len(rec))
	}
	if rec[0].(string) != "" {
		t.Errorf("Read(): rec[0]: want: %#v, got: %#v", "", rec[0])
	}
	if rec[1].(bool) != false {
		t.Errorf("Read(): rec[1]: want: %#v, got: %#v", false, rec[1])
	}
	if rec[2].(int64) != 0 {
		t.Errorf("Read(): rec[2]: want: %#v, got: %#v", 0, rec[2])
	}
	if rec[3].(float64) != 0 {
		t.Errorf("Read(): rec[3]: want: %#v, got: %#v", 0, rec[3])
	}
	var d time.Time
	if rec[4].(time.Time) != d {
		t.Errorf("Read(): rec[4]: want: %v, got: %v", d, rec[4])
	}
}

func Test_Reader_Read_records(t *testing.T) {
	fname := "./testdata/rec3.dbf"
	f, err := os.Open(fname)
	if err != nil {
		t.Errorf("os.Open(%q): %v", fname, err)
	}
	defer f.Close()

	r, err := NewReader(f)
	if err != nil {
		t.Errorf("NewReader(): %v", err)
	}
	if r.RecordCount() != 3 {
		t.Errorf("NewReader(): r.RecordCount(): want: %v, got: %v", 3, r.RecordCount())
	}

	var d time.Time
	d1 := time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC)

	testRecords := [][]interface{}{
		{"Abc", true, int64(123), float64(123.45), d1},
		{"", false, int64(0), float64(0), d},
		{"Мышь", false, int64(-321), float64(-54.32), d1},
	}

	for i := uint32(0); i < r.RecordCount(); i++ {
		want := testRecords[int(i)]
		got, err := r.Read()

		if err != nil {
			t.Errorf("Read(): %v", err)
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Read():\nwant: %v\ngot : %v", want, got)
		}
	}
}
