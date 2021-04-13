package xbase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXBaseNew(t *testing.T) {
	db := New()
	require.Equal(t, 0, db.FieldCount())
	require.Equal(t, int64(0), db.RecCount())
	require.Equal(t, int64(0), db.RecNo())
	require.Equal(t, true, db.BOF())
	require.Equal(t, true, db.EOF())
}
