package dbf

import (
	"io"
	"os"
	"strings"
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
	require.Equal(t, len(testFields), fields.Count())
	for i, f := range testFields {
		name, typ, length, dec := fields.Get(i)
		require.Equal(t, f.Name, name)
		require.Equal(t, f.Type, typ)
		require.Equal(t, f.Len, length)
		require.Equal(t, f.Dec, dec)
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
	require.Equal(t, 5, r.fields.Count())
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
func BenchmarkBoolValue1(b *testing.B) {
	buf := []byte{'T'}
	for i := 0; i < b.N; i++ {
		s := string(buf)
		if s[0] != 'T' {
			b.Fatalf("Fail bool value 1")
		}
	}
}

func BenchmarkBoolValue2(b *testing.B) {
	buf := []byte{'T'}
	for i := 0; i < b.N; i++ {
		// s := string(buf)
		if buf[0] != 'T' {
			b.Fatalf("Fail bool value 2")
		}
	}
}

func BenchmarkStringValue1(b *testing.B) {
	buf := []byte("Abc   ")
	for i := 0; i < b.N; i++ {
		s := string(buf)
		s = strings.TrimRight(s, " ")
		if s != "Abc" {
			b.Fatalf("Fail string value 1")
		}
	}
}

func BenchmarkStringValue2(b *testing.B) {
	buf := []byte("Abc   ")
	for i := 0; i < b.N; i++ {
		// s := string(bytes.TrimRight(buf, " "))
		s := trimRight(buf)
		if s != "Abc" {
			b.Fatalf("Fail string value %q", s)
		}
	}
}

func BenchmarkStringEmpty1(b *testing.B) {
	buf := []byte("        ")
	for i := 0; i < b.N; i++ {
		s := string(buf)
		if strings.Trim(s, " ") != "" {
			b.Fatalf("Fail")
		}
	}
}

func BenchmarkStringEmpty2(b *testing.B) {
	buf := []byte("        ")
	for i := 0; i < b.N; i++ {
		if !isEmpty(buf) {
			b.Fatalf("Fail")
		}
	}
}
