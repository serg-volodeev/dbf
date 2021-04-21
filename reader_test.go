package dbf

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReaderNewReaderEmptyReader(t *testing.T) {
	r, err := NewReader(nil)
	require.Error(t, err)
	require.Nil(t, r)
}

func TestReaderNewReader(t *testing.T) {
	f, err := os.Open("./testdata/rec3.dbf")
	require.NoError(t, err)
	defer f.Close()

	r, err := NewReader(f)
	require.NoError(t, err)
	require.Equal(t, uint32(3), r.RecordCount())
	require.Equal(t, 866, r.CodePage())

	testFields := []FieldInfo{
		{"NAME", "C", 20, 0},
		{"FLAG", "L", 1, 0},
		{"COUNT", "N", 5, 0},
		{"PRICE", "N", 9, 2},
		{"DATE", "D", 8, 0},
	}
	fields := r.Fields()
	require.Equal(t, len(testFields), len(fields))
	for i, f := range testFields {
		require.Equal(t, f.Name, fields[i].Name)
		require.Equal(t, f.Type, fields[i].Type)
		require.Equal(t, f.Len, fields[i].Len)
		require.Equal(t, f.Dec, fields[i].Dec)
	}
}

func TestReaderReadFileEmpty(t *testing.T) {
	f, err := os.Open("./testdata/rec0.dbf")
	require.NoError(t, err)
	defer f.Close()

	r, err := NewReader(f)
	require.NoError(t, err)

	rec, err := r.Read()
	require.Error(t, err)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, len(rec))
	require.Equal(t, uint32(0), r.header.RecCount)
	require.Equal(t, 5, len(r.fields))
}

func TestReaderReadRecordEmpty(t *testing.T) {
	f, err := os.Open("./testdata/rec1.dbf")
	require.NoError(t, err)
	defer f.Close()

	r, err := NewReader(f)
	require.NoError(t, err)

	rec, err := r.Read()
	require.NoError(t, err)
	require.Equal(t, 5, len(rec))
	require.Equal(t, "", rec[0].(string))
	require.Equal(t, false, rec[1].(bool))
	require.Equal(t, int64(0), rec[2].(int64))
	require.Equal(t, float64(0), rec[3].(float64))
	var d time.Time
	require.Equal(t, d, rec[4].(time.Time))
}

func TestReaderReadRecords(t *testing.T) {
	f, err := os.Open("./testdata/rec3.dbf")
	require.NoError(t, err)
	defer f.Close()

	r, err := NewReader(f)
	require.NoError(t, err)

	rec, err := r.Read()
	require.NoError(t, err)
	require.Equal(t, "Abc", rec[0].(string))
	require.Equal(t, true, rec[1].(bool))
	require.Equal(t, int64(123), rec[2].(int64))
	require.Equal(t, float64(123.45), rec[3].(float64))
	d1 := time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC)
	require.Equal(t, d1, rec[4].(time.Time))

	rec, err = r.Read()
	require.NoError(t, err)
	require.Equal(t, "", rec[0].(string))
	require.Equal(t, false, rec[1].(bool))
	require.Equal(t, int64(0), rec[2].(int64))
	require.Equal(t, float64(0), rec[3].(float64))
	var d time.Time
	require.Equal(t, d, rec[4].(time.Time))

	rec, err = r.Read()
	require.NoError(t, err)
	require.Equal(t, "Мышь", rec[0].(string))
	require.Equal(t, false, rec[1].(bool))
	require.Equal(t, int64(-321), rec[2].(int64))
	require.Equal(t, float64(-54.32), rec[3].(float64))
	d1 = time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC)
	require.Equal(t, d1, rec[4].(time.Time))

	rec, err = r.Read()
	require.Error(t, err)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, len(rec))
}
