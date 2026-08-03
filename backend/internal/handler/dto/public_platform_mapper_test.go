package dto

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestPublicPlatformMappersHideVideoProvider(t *testing.T) {
	group := &service.Group{ID: 1, Platform: service.PlatformLeo}

	publicKey := APIKeyFromServicePublic(&service.APIKey{Group: group})
	require.Equal(t, service.PublicVideoPlatform, publicKey.Group.Platform)

	publicSubscription := UserSubscriptionFromServicePublic(&service.UserSubscription{Group: group})
	require.Equal(t, service.PublicVideoPlatform, publicSubscription.Group.Platform)

	adminGroup := GroupFromServiceAdmin(group)
	require.Equal(t, service.PlatformLeo, adminGroup.Platform)
}
