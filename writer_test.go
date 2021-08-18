package dbf

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_Writer_prefix_error(t *testing.T) {
	_, err := NewWriter(nil, nil, 866)
	if !strings.HasPrefix(err.Error(), "dbf.NewWriter:") {
		t.Errorf("NewWriter(): require error prefix")
	}
}

func Test_Writer_write_records(t *testing.T) {
	fname := "./testdata/test.dbf"
	f, err := os.Create(fname)
	if err != nil {
		t.Errorf("os.Create(%s): %v", fname, err)
	}
	defer f.Close()

	fields := NewFields()
	fields.AddCharacterField("NAME", 20)
	fields.AddLogicalField("FLAG")
	fields.AddNumericField("COUNT", 5, 0)
	fields.AddNumericField("PRICE", 9, 2)
	fields.AddDateField("DATE")

	w, err := NewWriter(f, fields, 866)
	if err != nil {
		t.Errorf("NewWriter(): %v", err)
	}

	var d time.Time
	d1 := time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC)

	records := [][]interface{}{
		{"Abc", true, 123, 123.45, d1},
		{"", false, 0, 0, d},
		{"Мышь", false, -321, -54.32, d1},
	}

	for i := range records {
		err := w.Write(records[i])
		if err != nil {
			t.Errorf("Write(%v): %v", records[i], err)
		}
	}

	err = w.Flush()
	if err != nil {
		t.Errorf("Flush(): %v", err)
	}

	got := readFile("./testdata/test.dbf")
	want := readFile("./testdata/rec4.dbf")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("dbf file bytes:\nwant: %#v\ngot : %#v", want, got)
	}
}

func readFile(name string) []byte {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	// ModDate
	b[1] = 0
	b[2] = 0
	b[3] = 0
	// CodePage
	// b[29] = 0
	return b
}
