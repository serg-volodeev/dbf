package dbf

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// New field

func TestNewLogicalField(t *testing.T) {
	f := newLogicalField("Flag")
	require.Equal(t, "FLAG", f.name())
	require.Equal(t, byte('L'), f.Type)
	require.Equal(t, byte(1), f.Len)
	require.Equal(t, byte(0), f.Dec)
}

func TestNewDateField(t *testing.T) {
	f := newDateField("Date")
	require.Equal(t, "DATE", f.name())
	require.Equal(t, byte('D'), f.Type)
	require.Equal(t, byte(8), f.Len)
	require.Equal(t, byte(0), f.Dec)
}

func TestNewCharacterField(t *testing.T) {
	f := newCharacterField("name", 25)
	require.Equal(t, "NAME", f.name())
	require.Equal(t, byte('C'), f.Type)
	require.Equal(t, byte(25), f.Len)
	require.Equal(t, byte(0), f.Dec)
}

func TestNewNumericField(t *testing.T) {
	f := newNumericField("price", 12, 2)
	require.Equal(t, "PRICE", f.name())
	require.Equal(t, byte('N'), f.Type)
	require.Equal(t, byte(12), f.Len)
	require.Equal(t, byte(2), f.Dec)
}

// Field name

func TestFieldName(t *testing.T) {
	f := &field{
		Name: [11]byte{'N', 'A', 'M', 'E', 0, 0, 0, 0, 0, 0},
	}
	require.Equal(t, "NAME", f.name())
}

func TestFieldSetName(t *testing.T) {
	f := &field{}
	f.setName("name")
	require.Equal(t, "NAME", f.name())
}

// Field read/write

var fieldBytes = []byte{'N', 'A', 'M', 'E', 0, 0, 0, 0, 0, 0, 0, 'C', 0, 0, 0, 0, 14, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func TestFieldRead(t *testing.T) {
	r := bytes.NewReader(fieldBytes)

	f := &field{}
	err := f.read(r)

	require.NoError(t, err)
	require.Equal(t, "NAME", f.name())
	require.Equal(t, byte('C'), f.Type)
	require.Equal(t, byte(14), f.Len)
	require.Equal(t, byte(0), f.Dec)
}

func TestFieldWrite(t *testing.T) {
	f := newCharacterField("name", 14)

	buf := bytes.NewBuffer(nil)
	err := f.write(buf)

	require.NoError(t, err)
	require.Equal(t, fieldBytes, buf.Bytes())
}

// Value to buffer

func TestCharacterToBuf(t *testing.T) {
	f := newCharacterField("name", 6)

	var v interface{} = "Abc"
	s, err := f.characterToBuf(v, nil)
	require.NoError(t, err)
	require.Equal(t, "Abc   ", s)

	v = true
	_, err = f.characterToBuf(v, nil)
	require.Error(t, err)
}

func TestLogicalToBuf(t *testing.T) {
	f := newLogicalField("name")
	var v interface{} = false
	s, err := f.logicalToBuf(v)
	require.NoError(t, err)
	require.Equal(t, "F", s)

	v = "abc"
	_, err = f.logicalToBuf(v)
	require.Error(t, err)
}

func TestDateToBuf(t *testing.T) {
	f := newDateField("name")
	var v interface{} = time.Date(2021, 7, 26, 0, 0, 0, 0, time.UTC)
	s, err := f.dateToBuf(v)
	require.NoError(t, err)
	require.Equal(t, "20210726", s)

	v = "abc"
	_, err = f.logicalToBuf(v)
	require.Error(t, err)
}

func TestNumericIntToBuf(t *testing.T) {
	f := newNumericField("name", 6, 0)
	var v interface{} = -123
	s, err := f.numericToBuf(v)
	require.NoError(t, err)
	require.Equal(t, "  -123", s)

	v = "abc"
	_, err = f.logicalToBuf(v)
	require.Error(t, err)
}

func TestNumericFloatToBuf(t *testing.T) {
	f := newNumericField("name", 9, 2)
	var v interface{} = -123.4
	s, err := f.numericToBuf(v)
	require.NoError(t, err)
	require.Equal(t, "  -123.40", s)
}

func TestNumericFloatNullToBuf(t *testing.T) {
	f := newNumericField("name", 9, 2)
	var v interface{} = 0
	s, err := f.numericToBuf(v)
	require.NoError(t, err)
	require.Equal(t, "     0.00", s)
}

// Buffer to value

func TestBufToCharacter(t *testing.T) {
	f := newCharacterField("name", 6)
	v, err := f.bufToCharacter([]byte("Abc   "), nil)
	require.NoError(t, err)
	require.Equal(t, "Abc", v.(string))
}

func TestBufToLogical(t *testing.T) {
	f := newLogicalField("name")
	v := f.bufToLogical([]byte("T"))
	require.Equal(t, true, v.(bool))
}

func TestBufToDate(t *testing.T) {
	f := newDateField("name")
	v, err := f.bufToDate([]byte("20210727"))
	require.NoError(t, err)
	require.Equal(t, time.Date(2021, 7, 27, 0, 0, 0, 0, time.UTC), v.(time.Time))
}

func TestBufToNumericInt(t *testing.T) {
	f := newNumericField("name", 5, 0)
	v, err := f.bufToNumeric([]byte(" -123"))
	require.NoError(t, err)
	require.Equal(t, int64(-123), v.(int64))
}

func TestBufToNumericFloat(t *testing.T) {
	f := newNumericField("name", 8, 2)
	v, err := f.bufToNumeric([]byte(" -123.45"))
	require.NoError(t, err)
	require.Equal(t, float64(-123.45), v.(float64))
}
