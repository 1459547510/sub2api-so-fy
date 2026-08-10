package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func getCyberPolicyUserMark(ctx context.Context, rdb *redis.Client, userID int64) (bool, bool, error) {
	if rdb == nil || userID <= 0 {
		return false, false, errors.New("cyber policy user marker cache unavailable")
	}
	value, err := rdb.Get(ctx, serviceCyberPolicyUserMarkerKey(userID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	switch value {
	case "1":
		return true, true, nil
	case "0":
		return false, true, nil
	default:
		return false, false, errors.New("invalid cyber policy user marker cache value")
	}
}

func setCyberPolicyUserMark(ctx context.Context, rdb *redis.Client, userID int64, marked bool) error {
	if rdb == nil || userID <= 0 {
		return errors.New("cyber policy user marker cache unavailable")
	}
	value := "0"
	ttl := serviceCyberPolicyUserMarkerNegativeTTL
	if marked {
		value = "1"
		ttl = 0
	}
	return rdb.Set(ctx, serviceCyberPolicyUserMarkerKey(userID), value, ttl).Err()
}

func serviceCyberPolicyUserMarkerKey(userID int64) string {
	return fmt.Sprintf("cyber_policy_user:%d", userID)
}

const serviceCyberPolicyUserMarkerNegativeTTL = 5 * time.Minute

func (c *gatewayCache) GetCyberPolicyUserMark(ctx context.Context, userID int64) (bool, bool, error) {
	if c == nil {
		return false, false, errors.New("cyber policy user marker cache unavailable")
	}
	return getCyberPolicyUserMark(ctx, c.rdb, userID)
}

func (c *gatewayCache) SetCyberPolicyUserMark(ctx context.Context, userID int64, marked bool) error {
	if c == nil {
		return errors.New("cyber policy user marker cache unavailable")
	}
	return setCyberPolicyUserMark(ctx, c.rdb, userID, marked)
}

func (c *contentModerationHashCache) GetCyberPolicyUserMark(ctx context.Context, userID int64) (bool, bool, error) {
	if c == nil {
		return false, false, errors.New("cyber policy user marker cache unavailable")
	}
	return getCyberPolicyUserMark(ctx, c.rdb, userID)
}

func (c *contentModerationHashCache) SetCyberPolicyUserMark(ctx context.Context, userID int64, marked bool) error {
	if c == nil {
		return errors.New("cyber policy user marker cache unavailable")
	}
	return setCyberPolicyUserMark(ctx, c.rdb, userID, marked)
}
