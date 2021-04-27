package dbf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldsNewFields(t *testing.T) {
	f := NewFields()
	require.NoError(t, f.Error())
	require.Equal(t, 0, f.Count())
}

func TestFieldsAdd(t *testing.T) {
	f := NewFields()
	f.Add("date", "d")
	require.NoError(t, f.Error())
	require.Equal(t, 1, f.Count())
}

func TestFieldsGet(t *testing.T) {
	f := NewFields()
	f.Add("date", "d")
	name, typ, length, dec := f.Get(0)
	require.Equal(t, "DATE", name)
	require.Equal(t, "D", typ)
	require.Equal(t, 8, length)
	require.Equal(t, 0, dec)
}

func TestFieldsWrite(t *testing.T) {
	f := NewFields()
	f.Add("name", "C", 14)
	f.Add("date", "D")

	buf := bytes.NewBuffer(nil)
	err := f.write(buf)

	require.NoError(t, err)
	b := buf.Bytes()
	require.Equal(t, 64, len(b))
	require.Equal(t, byte('N'), b[0])
	require.Equal(t, byte('A'), b[1])
	require.Equal(t, byte('M'), b[2])
	require.Equal(t, byte('E'), b[3])
	require.Equal(t, byte('C'), b[11])
	require.Equal(t, byte(0), b[12])
	require.Equal(t, byte(14), b[16])
	require.Equal(t, byte('D'), b[32])
	require.Equal(t, byte('A'), b[33])
}

func TestFieldsRead(t *testing.T) {
	b := make([]byte, fieldSize)
	copy(b[:], "NAME")
	b[11] = 'C'
	b[12] = 1
	b[16] = 14
	r := bytes.NewReader(b)

	f := NewFields()
	err := f.read(r, 1)

	require.NoError(t, err)
	require.Equal(t, 1, f.Count())
}
