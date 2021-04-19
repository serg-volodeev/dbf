package dbf

import (
	"testing"
	"time"

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

func TestXBaseOpenEmptyFile(t *testing.T) {
	db := New()
	err := db.OpenFile("./testdata/rec0.dbf", true)
	require.NoError(t, err)
	require.Equal(t, int64(0), db.RecCount())
	require.Equal(t, 5, db.FieldCount())
	require.Equal(t, true, db.EOF())
	require.Equal(t, true, db.BOF())

	err = db.First()
	require.NoError(t, err)
	require.Equal(t, true, db.EOF())
	require.Equal(t, true, db.BOF())

	err = db.Next()
	require.NoError(t, err)
	require.Equal(t, true, db.EOF())
	require.Equal(t, true, db.BOF())

	err = db.Last()
	require.NoError(t, err)
	require.Equal(t, true, db.EOF())
	require.Equal(t, true, db.BOF())

	err = db.Prev()
	require.NoError(t, err)
	require.Equal(t, true, db.EOF())
	require.Equal(t, true, db.BOF())

	err = db.CloseFile()
	require.NoError(t, err)
}

func TestXBaseReadEmptyRec(t *testing.T) {
	db := New()
	err := db.OpenFile("./testdata/rec1.dbf", true)
	require.NoError(t, err)

	err = db.First()
	require.NoError(t, err)

	v1, err := db.FieldValueAsString(1)
	require.NoError(t, err)
	require.Equal(t, "", v1)

	v2, err := db.FieldValueAsBool(2)
	require.NoError(t, err)
	require.Equal(t, false, v2)

	v3, err := db.FieldValueAsInt(3)
	require.NoError(t, err)
	require.Equal(t, int64(0), v3)

	v4, err := db.FieldValueAsFloat(4)
	require.NoError(t, err)
	require.Equal(t, float64(0), v4)

	var d time.Time
	v5, err := db.FieldValueAsDate(5)
	require.NoError(t, err)
	require.Equal(t, d, v5)

	_, err = db.FieldValueAsDate(6)
	require.Error(t, err)

	err = db.CloseFile()
	require.NoError(t, err)
}

func TestXBaseReadNext(t *testing.T) {
	db := New()
	err := db.OpenFile("./testdata/rec3.dbf", true)
	require.NoError(t, err)

	err = db.First()
	require.NoError(t, err)
	require.Equal(t, int64(1), db.RecNo())
	require.Equal(t, false, db.EOF())
	v1, err := db.FieldValueAsString(1)
	require.NoError(t, err)
	require.Equal(t, "Abc", v1)

	err = db.Next()
	require.NoError(t, err)
	require.Equal(t, int64(2), db.RecNo())
	require.Equal(t, false, db.EOF())
	v1, err = db.FieldValueAsString(1)
	require.NoError(t, err)
	require.Equal(t, "", v1)

	err = db.Next()
	require.NoError(t, err)
	require.Equal(t, int64(3), db.RecNo())
	require.Equal(t, false, db.EOF())
	v1, err = db.FieldValueAsString(1)
	require.NoError(t, err)
	require.Equal(t, "Мышь", v1)

	err = db.Next()
	require.NoError(t, err)
	require.Equal(t, true, db.EOF())

	err = db.Next()
	require.NoError(t, err)
	require.Equal(t, true, db.EOF())

	err = db.CloseFile()
	require.NoError(t, err)
}
