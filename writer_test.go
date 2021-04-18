package xbase

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWriterNewWriter(t *testing.T) {
	fields := []FieldInfo{
		{"Name", "C", 20, 0},
		{"Count", "N", 10, 0},
	}
	w, err := NewWriter(nil, fields, 866)
	require.NoError(t, err)
	require.Equal(t, 2, len(w.fields))
	require.NotNil(t, w.encoder)
}

func TestWriterWriteRecords(t *testing.T) {
	f, err := os.Create("./testdata/test.dbf")
	require.NoError(t, err)
	defer f.Close()

	fields := []FieldInfo{
		{"NAME", "C", 20, 0},
		{"FLAG", "L", 0, 0},
		{"COUNT", "N", 5, 0},
		{"PRICE", "N", 9, 2},
		{"DATE", "D", 0, 0},
	}
	w, err := NewWriter(f, fields, 866)
	require.NoError(t, err)

	var d time.Time
	d1 := time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC)

	records := [][]interface{}{
		{"Abc", true, 123, 123.45, d1},
		{"", false, 0, 0.0, d},
		{"Мышь", false, -321, -54.32, d1},
	}

	for i := range records {
		err := w.Write(records[i])
		require.NoError(t, err)
	}

	err = w.Flush()
	require.NoError(t, err)

	testBytes := readFile("./testdata/test.dbf")
	goldBytes := readFile("./testdata/rec3.dbf")
	require.Equal(t, goldBytes, testBytes)
}

func readFile(name string) []byte {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
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
