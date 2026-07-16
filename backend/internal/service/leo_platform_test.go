package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeoPlatformRegistration(t *testing.T) {
	require.Equal(t, "leo", PlatformLeo)
	require.True(t, IsAllowedQuotaPlatform(PlatformLeo))
	require.Contains(t, AllowedQuotaPlatforms, PlatformLeo)
}
