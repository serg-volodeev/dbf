package dbf

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestTrimRight(t *testing.T) {
	require.Equal(t, "Abc", trimRight([]byte("Abc")))
	require.Equal(t, "Abc", trimRight([]byte("Abc   ")))
	require.Equal(t, "", trimRight([]byte("   ")))
}

func TestTrimLeft(t *testing.T) {
	require.Equal(t, "Abc", trimLeft([]byte("Abc")))
	require.Equal(t, "Abc", trimLeft([]byte("    Abc")))
	require.Equal(t, "", trimLeft([]byte("   ")))
}

func TestParse(t *testing.T) {
	n, err := strconv.ParseFloat("-.1", 64)
	require.NoError(t, err)
	require.Equal(t, float64(-0.1), n)
}
