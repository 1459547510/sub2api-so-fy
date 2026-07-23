package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeVideoBillingResolutionLeo(t *testing.T) {
	require.Equal(t, VideoBillingResolution480P, NormalizeVideoBillingResolutionOrDefault("RESOLUTION_480"))
	require.Equal(t, VideoBillingResolution720P, NormalizeVideoBillingResolutionOrDefault("RESOLUTION_720"))
	require.Equal(t, VideoBillingResolution1080P, NormalizeVideoBillingResolutionOrDefault("RESOLUTION_1080"))
}

func TestLeoVideoPricingResolutions(t *testing.T) {
	require.Equal(t, []string{"480p", "720p", "1080p"}, LeoVideoPricingResolutions("seedance-2.0"))
	require.Equal(t, []string{"720p"}, LeoVideoPricingResolutions("seedance-2.0-mini"))
	require.Equal(t, "720p", DefaultLeoVideoResolution("seedance-2.0-mini"))
	require.True(t, LeoVideoModelSupportsResolution("seedance-2.0-fast", "480p"))
	require.False(t, LeoVideoModelSupportsResolution("seedance-2.0-fast", "1080p"))
	require.True(t, LeoVideoModelSupportsResolution("seedance-2.0-mini", "720p"))
	require.False(t, LeoVideoModelSupportsResolution("seedance-2.0-mini", "1080p"))
	require.False(t, LeoVideoModelSupportsResolution("seedance-2.0", "invalid"))
}
