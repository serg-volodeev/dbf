package dbf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFields(t *testing.T) {
	f := NewFields()
	require.NoError(t, f.Error())
	require.Equal(t, 0, f.Count())
}
