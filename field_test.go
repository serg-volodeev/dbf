package dbf

import (
	"bytes"
	"testing"

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
