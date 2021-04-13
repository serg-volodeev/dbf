package xbase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldName(t *testing.T) {
	f := &field{
		Name: [11]byte{'N', 'A', 'M', 'E', 0, 0, 0, 0, 0, 0},
	}
	require.Equal(t, "NAME", f.name())
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

func TestNewField(t *testing.T) {
	f, err := newField("Price", "N", 12, 2)
	require.NoError(t, err)
	require.Equal(t, "PRICE", f.name())
	require.Equal(t, byte('N'), f.Type)
	require.Equal(t, byte(12), f.Len)
	require.Equal(t, byte(2), f.Dec)
}

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
