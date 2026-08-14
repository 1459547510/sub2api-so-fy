package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeVideoBillingResolutionLeo(t *testing.T) {
	require.Equal(t, VideoBillingResolution480P, NormalizeVideoBillingResolutionOrDefault("RESOLUTION_480"))
	require.Equal(t, VideoBillingResolution480P, NormalizeVideoBillingResolutionOrDefault("400p"))
	require.Equal(t, VideoBillingResolution480P, NormalizeVideoBillingResolutionOrDefault("544p"))
	require.Equal(t, VideoBillingResolution720P, NormalizeLeoVideoBillingResolutionOrDefault("happy-horse-1.1", "720p"))
	require.Equal(t, VideoBillingResolution1080P, NormalizeLeoVideoBillingResolutionOrDefault("happy-horse-1.1", "1080p"))
	require.Equal(t, VideoBillingResolution400P, NormalizeLeoVideoBillingResolutionOrDefault("grok-imagine-1.5", "400p"))
	require.Equal(t, VideoBillingResolution544P, NormalizeLeoVideoBillingResolutionOrDefault("grok-imagine-1.5", "544p"))
	require.Equal(t, VideoBillingResolution720P, NormalizeLeoVideoBillingResolutionOrDefault("grok-imagine-1.5", "720p"))
	require.Equal(t, VideoBillingResolution960P, NormalizeLeoVideoBillingResolutionOrDefault("grok-imagine-1.5", "960p"))
	require.Equal(t, VideoBillingResolution720P, NormalizeVideoBillingResolutionOrDefault("RESOLUTION_720"))
	require.Equal(t, VideoBillingResolution1080P, NormalizeVideoBillingResolutionOrDefault("RESOLUTION_1080"))
	require.Equal(t, VideoBillingResolution1080P, NormalizeVideoBillingResolutionOrDefault("960p"))
	require.Equal(t, VideoBillingResolution1080P, NormalizeVideoBillingResolutionOrDefault("1440p"))
	require.Equal(t, VideoBillingResolution1080P, NormalizeVideoBillingResolutionOrDefault("2160p"))
	require.Equal(t, VideoBillingResolution1080P, NormalizeLeoVideoBillingResolutionOrDefault("ltx-2.3-pro", "RESOLUTION_1080"))
	require.Equal(t, VideoBillingResolution1440P, NormalizeLeoVideoBillingResolutionOrDefault("ltx-2.3-fast", "RESOLUTION_1440"))
	require.Equal(t, VideoBillingResolution2160P, NormalizeLeoVideoBillingResolutionOrDefault("ltx-2.3-fast", "RESOLUTION_2160"))
	require.Equal(t, VideoBillingResolution2160P, NormalizeLeoVideoBillingResolutionOrDefault("ltxv-2.3-fast", "RESOLUTION_2160"))
	require.Equal(t, 15, NormalizeVideoBillingDurationSecondsOrDefault(20))
	require.Equal(t, 20, NormalizeLeoVideoBillingDurationSecondsOrDefault("ltx-2.3-fast", 20))
	require.Equal(t, 15, NormalizeLeoVideoBillingDurationSecondsOrDefault("grok-imagine-1.5", 20))
}

func TestLeoVideoPricingResolutions(t *testing.T) {
	require.Equal(t, []string{"480p", "720p", "1080p"}, LeoVideoPricingResolutions("seedance-2.0"))
	require.Equal(t, []string{"480p", "720p"}, LeoVideoPricingResolutions("bytedance/seedance-2.5"))
	require.Equal(t, []string{"480p", "720p"}, LeoVideoPricingResolutions("seedance-2.5"))
	require.Equal(t, []string{"720p"}, LeoVideoPricingResolutions("seedance-2.0-mini"))
	require.Equal(t, "720p", DefaultLeoVideoResolution("seedance-2.0-mini"))
	require.Equal(t, "1080p", DefaultLeoVideoResolution("happy-horse-1.1"))
	require.Equal(t, []string{"720p", "1080p"}, LeoVideoPricingResolutions("happy-horse-1.1"))
	require.Equal(t, []string{"400p", "544p", "720p", "960p"}, LeoVideoPricingResolutions("grok-imagine-1.5"))
	require.Equal(t, []string{"1080p", "1440p", "2160p"}, LeoVideoPricingResolutions("ltx-2.3-pro"))
	require.Equal(t, []string{"1080p", "1440p", "2160p"}, LeoVideoPricingResolutions("ltx-2.3-fast"))
	require.Equal(t, []string{"1080p", "1440p", "2160p"}, LeoVideoPricingResolutions("ltxv-2.3-fast"))
	require.Equal(t, "1080p", DefaultLeoVideoResolution("ltx-2.3-fast"))
	require.True(t, LeoVideoModelSupportsResolution("ltx-2.3-pro", "2160p"))
	require.True(t, LeoVideoModelSupportsResolution("seedance-2.0-fast", "480p"))
	require.False(t, LeoVideoModelSupportsResolution("seedance-2.0-fast", "1080p"))
	require.True(t, LeoVideoModelSupportsResolution("seedance-2.0-mini", "720p"))
	require.False(t, LeoVideoModelSupportsResolution("seedance-2.0-mini", "1080p"))
	require.False(t, LeoVideoModelSupportsResolution("seedance-2.0", "invalid"))
}
