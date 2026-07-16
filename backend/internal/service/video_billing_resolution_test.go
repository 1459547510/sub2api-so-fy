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
