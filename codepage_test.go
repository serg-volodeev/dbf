package dbf

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/charmap"
)

// New field

func TestCodeByPage(t *testing.T) {
	require.Equal(t, byte(0x65), codeByPage(866))
}

func TestPageByCode(t *testing.T) {
	require.Equal(t, 866, pageByCode(0x65))
}

func TestCharmapByPage(t *testing.T) {
	require.Equal(t, charmap.CodePage866, charmapByPage(866))
}
