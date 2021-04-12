package xbase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHeader(t *testing.T) {
	h := newHeader()

	require.Equal(t, byte(0x03), h.Id)
	// require.Equal(t, uint32(0), h.RecCount)
	// require.Equal(t, -1, h.fieldCount())
	// require.Equal(t, uint16(0), h.RecSize)
	// require.Equal(t, currentDate(), h.modDate())
}
