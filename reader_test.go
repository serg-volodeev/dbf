package dbf

import (
	"os"
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

	ok := r.Read()

	if ok {
		t.Errorf("Read(): want: %v, got: %v", false, ok)
	}
	if r.Err() != nil {
		t.Errorf("Read(): Err(): want: %v, got: %v", nil, r.Err())
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

	ok := r.Read()

	if !ok {
		t.Errorf("Read(): want: %v, got: %v", true, ok)
	}
	if r.StringFieldValue(0) != "" {
		t.Errorf("Read(): r.StringFieldValue(0): want: %#v, got: %#v", "", r.StringFieldValue(0))
	}
	if r.BoolFieldValue(1) != false {
		t.Errorf("Read(): r.BoolFieldValue(1): want: %#v, got: %#v", false, r.BoolFieldValue(1))
	}
	if r.IntFieldValue(2) != 0 {
		t.Errorf("Read(): r.IntFieldValue(2): want: %#v, got: %#v", 0, r.IntFieldValue(2))
	}
	if r.FloatFieldValue(3) != 0 {
		t.Errorf("Read(): r.FloatFieldValue(3): want: %#v, got: %#v", 0, r.FloatFieldValue(3))
	}
	var d time.Time
	if r.DateFieldValue(4) != d {
		t.Errorf("Read(): r.DateFieldValue(4): want: %#v, got: %#v", d, r.DateFieldValue(4))
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

	testRecords := []struct {
		name  string
		flag  bool
		count int64
		price float64
		date  time.Time
	}{
		{"Abc", true, 123, 123.45, d1},
		{"", false, 0, 0, d},
		{"Мышь", false, -321, -54.32, d1},
	}

	i := 0
	for r.Read() {
		want := testRecords[i]
		i++

		if r.StringFieldValue(0) != want.name {
			t.Errorf("r.StringFieldValue(0): want: %#v, got: %#v", want.name, r.StringFieldValue(0))
		}
		if r.BoolFieldValue(1) != want.flag {
			t.Errorf("r.BoolFieldValue(1): want: %#v, got: %#v", want.flag, r.BoolFieldValue(1))
		}
		if r.IntFieldValue(2) != want.count {
			t.Errorf("r.IntFieldValue(2): want: %#v, got: %#v", want.count, r.IntFieldValue(2))
		}
		if r.FloatFieldValue(3) != want.price {
			t.Errorf("r.FloatFieldValue(3): want: %#v, got: %#v", want.price, r.FloatFieldValue(3))
		}
		if r.DateFieldValue(4) != want.date {
			t.Errorf("r.DateFieldValue(4): want: %#v, got: %#v", want.date, r.DateFieldValue(4))
		}
	}
}
