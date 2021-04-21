package dbf

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// String utils

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

// Convert interface{}

func TestInterfaceToString(t *testing.T) {
	var x interface{} = "Abc"
	s, err := interfaceToString(x)
	require.NoError(t, err)
	require.Equal(t, "Abc", s)

	x = true
	_, err = interfaceToString(x)
	require.Error(t, err)
}

func TestInterfaceToBool(t *testing.T) {
	var x interface{} = true
	b, err := interfaceToBool(x)
	require.NoError(t, err)
	require.Equal(t, true, b)

	x = 5
	_, err = interfaceToBool(x)
	require.Error(t, err)
}

func TestInterfaceToInt(t *testing.T) {
	var x interface{} = 123
	n, err := interfaceToInt(x)
	require.NoError(t, err)
	require.Equal(t, int64(123), n)

	x = 5.1
	_, err = interfaceToInt(x)
	require.Error(t, err)
}

func TestInterfaceToFloat(t *testing.T) {
	var x interface{} = 123.45
	n, err := interfaceToFloat(x)
	require.NoError(t, err)
	require.Equal(t, float64(123.45), n)

	x = 5
	n, err = interfaceToFloat(x)
	require.NoError(t, err)
	require.Equal(t, float64(5), n)
}

func TestInterfaceToDate(t *testing.T) {
	d := time.Date(2020, 9, 23, 0, 0, 0, 0, time.UTC)
	var x interface{} = d
	r, err := interfaceToDate(x)
	require.NoError(t, err)
	require.Equal(t, d, r)

	x = 5
	r, err = interfaceToDate(x)
	require.Error(t, err)
}
