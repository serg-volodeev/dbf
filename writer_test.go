package xbase

import (
	"testing"

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
