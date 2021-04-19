package dbf

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFieldName(t *testing.T) {
	f := &field{
		Name: [11]byte{'N', 'A', 'M', 'E', 0, 0, 0, 0, 0, 0},
	}
	require.Equal(t, "NAME", f.name())
}

// New field

func TestNewField(t *testing.T) {
	f, err := newField("Price", "N", 12, 2)
	require.NoError(t, err)
	require.Equal(t, "PRICE", f.name())
	require.Equal(t, byte('N'), f.Type)
	require.Equal(t, byte(12), f.Len)
	require.Equal(t, byte(2), f.Dec)
}

func TestFieldSetName(t *testing.T) {
	f := &field{}

	err := f.setName("name")
	require.NoError(t, err)
	require.Equal(t, "NAME", f.name())

	err = f.setName("")
	require.Error(t, err)

	err = f.setName("longlongname")
	require.Error(t, err)
}

func TestFieldSetType(t *testing.T) {
	f := &field{}

	err := f.setType("numeric")
	require.NoError(t, err)
	require.Equal(t, byte('N'), f.Type)

	err = f.setType("")
	require.Error(t, err)

	err = f.setType("X")
	require.Error(t, err)
}

func TestFieldSetLen(t *testing.T) {
	f := &field{Type: 'L'}

	f.setType("L")
	err := f.setLen(0)
	require.NoError(t, err)
	require.Equal(t, byte(1), f.Len)

	f.setType("D")
	err = f.setLen(12)
	require.NoError(t, err)
	require.Equal(t, byte(8), f.Len)

	f.setType("C")
	err = f.setLen(120)
	require.NoError(t, err)
	require.Equal(t, byte(120), f.Len)

	err = f.setLen(255)
	require.Error(t, err)

	f.setType("N")
	err = f.setLen(14)
	require.NoError(t, err)
	require.Equal(t, byte(14), f.Len)

	err = f.setLen(20)
	require.Error(t, err)
}

func TestFieldSetDec(t *testing.T) {
	f := &field{}
	f.setType("N")
	f.setLen(5)

	err := f.setDec(2)
	require.NoError(t, err)
	require.Equal(t, byte(2), f.Dec)

	err = f.setDec(4)
	require.Error(t, err)
}

// String utils

func TestPadRight(t *testing.T) {
	require.Equal(t, "Abc   ", padRight("Abc", 6))
	require.Equal(t, "Abc", padRight("Abc", 2))
}

func TestPadLeft(t *testing.T) {
	require.Equal(t, "   Abc", padLeft("Abc", 6))
	require.Equal(t, "Abc", padLeft("Abc", 2))
}

func TestIsASCII(t *testing.T) {
	require.Equal(t, true, isASCII("Abc"))
	require.Equal(t, false, isASCII("AbЖc"))
}

// Field read/write

func TestFieldRead(t *testing.T) {
	b := make([]byte, fieldSize)
	copy(b[:], "NAME")
	b[11] = 'C'
	b[12] = 1
	b[16] = 14
	r := bytes.NewReader(b)

	f := &field{}
	err := f.read(r)

	require.NoError(t, err)
	require.Equal(t, "NAME", f.name())
	require.Equal(t, byte('C'), f.Type)
	require.Equal(t, uint32(1), f.Offset)
	require.Equal(t, byte(14), f.Len)
	require.Equal(t, byte(0), f.Dec)
}

func TestFieldWrite(t *testing.T) {
	f := &field{}
	copy(f.Name[:], "NAME")
	f.Type = 'C'
	f.Offset = 1
	f.Len = 14
	f.Dec = 0

	buf := bytes.NewBuffer(nil)
	err := f.write(buf)

	require.NoError(t, err)
	b := buf.Bytes()
	require.Equal(t, byte('N'), b[0])
	require.Equal(t, byte('A'), b[1])
	require.Equal(t, byte('M'), b[2])
	require.Equal(t, byte('E'), b[3])
	require.Equal(t, byte('C'), b[11])
	require.Equal(t, byte(0), b[12])
	require.Equal(t, byte(14), b[16])
}

// Buffer field in buffer record

func TestFieldBuffer(t *testing.T) {
	f, _ := newField("Log", "L", 1, 0)
	f.Offset = 6
	recBuf := []byte(" Abc  T 12")

	require.Equal(t, []byte("T"), f.buffer(recBuf))

	f.setBuffer(recBuf, "F")
	require.Equal(t, []byte("F"), f.buffer(recBuf))
}

// Get value

func TestFieldValueC(t *testing.T) {
	f, _ := newField("Name", "C", 5, 0)
	f.Offset = 3
	recBuf := []byte("   Abc    ")
	v, err := f.value(recBuf, nil)
	require.NoError(t, err)
	require.Equal(t, "Abc", v.(string))
}

func TestFieldValueL(t *testing.T) {
	f, _ := newField("Name", "L", 1, 0)
	f.Offset = 3
	recBuf := []byte("   T    ")
	v, err := f.value(recBuf, nil)
	require.NoError(t, err)
	require.Equal(t, true, v.(bool))
}

func TestFieldValueD(t *testing.T) {
	f, _ := newField("Name", "D", 8, 0)
	f.Offset = 3
	recBuf := []byte("   20200923    ")

	d := time.Date(2020, 9, 23, 0, 0, 0, 0, time.UTC)
	v, err := f.value(recBuf, nil)
	require.NoError(t, err)
	require.Equal(t, d, v.(time.Time))
}

func TestFieldValueN1(t *testing.T) {
	f, _ := newField("Name", "N", 8, 0)
	f.Offset = 3
	recBuf := []byte("      -2020    ")
	v, err := f.value(recBuf, nil)
	require.NoError(t, err)
	require.Equal(t, int64(-2020), v.(int64))
}

func TestFieldValueN2(t *testing.T) {
	f, _ := newField("Name", "N", 8, 2)
	f.Offset = 3
	recBuf := []byte("     -20.21    ")
	v, err := f.value(recBuf, nil)
	require.NoError(t, err)
	require.Equal(t, float64(-20.21), v.(float64))
}

// Set value

func TestFieldSetValueC(t *testing.T) {
	recBuf := make([]byte, 20)
	f, _ := newField("NAME", "C", 5, 0)
	f.Offset = 5
	err := f.setValue(recBuf, " Abc", nil)
	require.NoError(t, err)
	require.Equal(t, []byte(" Abc "), recBuf[5:10])
}

func TestFieldSetValueL(t *testing.T) {
	recBuf := make([]byte, 20)
	f, _ := newField("NAME", "L", 1, 0)
	f.Offset = 5
	err := f.setValue(recBuf, true, nil)
	require.NoError(t, err)
	require.Equal(t, []byte("T"), recBuf[5:6])
}

func TestFieldSetValueD(t *testing.T) {
	recBuf := make([]byte, 20)
	f, _ := newField("NAME", "D", 8, 0)
	f.Offset = 5
	d := time.Date(2020, 9, 23, 0, 0, 0, 0, time.UTC)
	err := f.setValue(recBuf, d, nil)
	require.NoError(t, err)
	require.Equal(t, []byte("20200923"), recBuf[5:13])
}

func TestFieldSetValueN1(t *testing.T) {
	recBuf := make([]byte, 20)
	f, _ := newField("NAME", "N", 5, 0)
	f.Offset = 5
	err := f.setValue(recBuf, 123, nil)
	require.NoError(t, err)
	require.Equal(t, []byte("  123"), recBuf[5:10])
}

func TestFieldSetValueN2(t *testing.T) {
	recBuf := make([]byte, 20)
	f, _ := newField("NAME", "N", 8, 2)
	f.Offset = 5
	err := f.setValue(recBuf, 123.45, nil)
	require.NoError(t, err)
	require.Equal(t, []byte("  123.45"), recBuf[5:13])
}
