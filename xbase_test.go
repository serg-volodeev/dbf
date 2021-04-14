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

func TestXBaseAddField(t *testing.T) {
	db := New()
	err := db.AddField("NAME", "C", 20)
	require.NoError(t, err)
	require.Equal(t, 1, db.FieldCount())
	require.Equal(t, 1, db.FieldNo("name"))

	name, typ, len, dec, err := db.FieldInfo(1)
	require.NoError(t, err)
	require.Equal(t, "NAME", name)
	require.Equal(t, "C", typ)
	require.Equal(t, 20, len)
	require.Equal(t, 0, dec)
}

func TestXBaseCodePage(t *testing.T) {
	db := New()
	db.SetCodePage(866)
	require.Equal(t, 866, db.CodePage())
}
