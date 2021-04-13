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
