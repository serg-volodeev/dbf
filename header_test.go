package xbase

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func currentDate() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestHeaderNew(t *testing.T) {
	h := newHeader()

	require.Equal(t, byte(0x03), h.Id)
	require.Equal(t, uint32(0), h.RecCount)
	require.Equal(t, -1, h.fieldCount())
	require.Equal(t, uint16(0), h.RecSize)
	require.Equal(t, currentDate(), h.modDate())
}

func TestHeaderWrite(t *testing.T) {
	h := newHeader()
	h.RecCount = uint32(3)
	h.RecSize = uint16(39)
	h.setFieldCount(5)
	h.setModDate(time.Date(1930, 2, 20, 0, 0, 0, 0, time.UTC))
	h.setCodePage(866)

	buf := bytes.NewBuffer(nil)
	h.write(buf)

	expected := []byte{0x3, 0x1e, 0x2, 0x14, 0x3, 0, 0, 0, 0xc1, 0, 0x27,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x65, 0, 0}
	require.Equal(t, expected, buf.Bytes())
}

func TestHeaderRead(t *testing.T) {
	b := []byte{0x3, 0x1e, 0x2, 0x14, 0x3, 0, 0, 0, 0xc1, 0, 0x27,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x65, 0, 0}
	r := bytes.NewReader(b)

	h := &header{}
	h.read(r)

	require.Equal(t, byte(0x03), h.Id)
	require.Equal(t, uint32(3), h.RecCount)
	require.Equal(t, 5, h.fieldCount())
	require.Equal(t, uint16(39), h.RecSize)
	require.Equal(t, 866, h.codePage())

	d := time.Date(1930, 2, 20, 0, 0, 0, 0, time.UTC)
	require.Equal(t, d, h.modDate())
}

func TestHeaderReadNotDBF(t *testing.T) {
	b := make([]byte, headerSize)
	b[0] = 0x05 // valid 0x03
	r := bytes.NewReader(b)

	h := &header{}
	err := h.read(r)
	require.Error(t, err)
}

func TestHeaderSetCodePage(t *testing.T) {
	h := &header{}
	h.setCodePage(866)
	require.Equal(t, byte(0x65), h.CP)
	require.Equal(t, 866, h.codePage())
}
