package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeLayer_withValue(t *testing.T) {
	name, ok := NormalizeLayer(" API ")
	require.True(t, ok)
	require.Equal(t, "API", name)
}

func TestNormalizeLayer_empty(t *testing.T) {
	_, ok := NormalizeLayer("")
	require.False(t, ok)
	_, ok = NormalizeLayer("   ")
	require.False(t, ok)
}

func TestNormalizeLayer_custom(t *testing.T) {
	name, ok := NormalizeLayer("my-custom-layer")
	require.True(t, ok)
	require.Equal(t, "my-custom-layer", name)
}

func TestTestLayers_constants(t *testing.T) {
	require.Equal(t, "API", TestLayers.API)
	require.Equal(t, "E2E", TestLayers.E2E)
}
