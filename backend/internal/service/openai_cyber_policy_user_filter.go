package service

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
)

// OpenAICyberPolicyUserBlockingExtraKey enables account-level filtering for
// users with a historical upstream cyber_policy event.
const OpenAICyberPolicyUserBlockingExtraKey = "openai_block_cyber_policy_users"

// CyberPolicyUserMarkerRepository is intentionally optional so existing
// UserRepository implementations and test doubles remain source-compatible.
type CyberPolicyUserMarkerRepository interface {
	MarkCyberPolicyUser(ctx context.Context, userID int64, at time.Time) (bool, error)
	HasCyberPolicyUser(ctx context.Context, userID int64) (bool, error)
}

// CyberPolicyUserMarkerCache is an optional cache for the permanent marker.
// found=false represents a cache miss; cache failures are fail-open.
type CyberPolicyUserMarkerCache interface {
	GetCyberPolicyUserMark(ctx context.Context, userID int64) (marked bool, found bool, err error)
	SetCyberPolicyUserMark(ctx context.Context, userID int64, marked bool) error
}

type cyberPolicyUserFilterContextKey struct{}

type cyberPolicyUserFilterState struct {
	once     sync.Once
	marked   bool
	resolved bool
}

func (a *Account) IsOpenAICyberPolicyUserBlockingEnabled() bool {
	return a != nil && a.Platform == PlatformOpenAI && resolveAccountExtraBool(a.Extra, OpenAICyberPolicyUserBlockingExtraKey)
}

func withCyberPolicyUserFilterState(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := ctx.Value(cyberPolicyUserFilterContextKey{}).(*cyberPolicyUserFilterState); ok {
		return ctx
	}
	return context.WithValue(ctx, cyberPolicyUserFilterContextKey{}, &cyberPolicyUserFilterState{})
}

func cyberPolicyUserFilterStateFromContext(ctx context.Context) *cyberPolicyUserFilterState {
	if ctx == nil {
		return nil
	}
	state, _ := ctx.Value(cyberPolicyUserFilterContextKey{}).(*cyberPolicyUserFilterState)
	return state
}

func (s *OpenAIGatewayService) withCyberPolicyUserFilterContext(ctx context.Context) context.Context {
	return withCyberPolicyUserFilterState(ctx)
}

func (s *OpenAIGatewayService) cyberPolicyUserMarked(ctx context.Context) bool {
	if s == nil || ctx == nil {
		return false
	}
	userID, ok := ctx.Value(ctxkey.UserID).(int64)
	if !ok || userID <= 0 {
		return false
	}
	if role, _ := ctx.Value(ctxkey.UserRole).(string); strings.EqualFold(strings.TrimSpace(role), RoleAdmin) {
		return false
	}

	state := cyberPolicyUserFilterStateFromContext(ctx)
	if state == nil {
		state = &cyberPolicyUserFilterState{}
	}
	state.once.Do(func() {
		state.resolved = true
		if cache, ok := s.cache.(CyberPolicyUserMarkerCache); ok && cache != nil {
			marked, found, err := cache.GetCyberPolicyUserMark(ctx, userID)
			if err != nil {
				slog.Warn("openai.cyber_policy_user_marker_cache_read_failed", "user_id", userID, "error", err)
				return
			} else if found {
				state.marked = marked
				return
			}
		}
		repo, ok := s.userRepo.(CyberPolicyUserMarkerRepository)
		if !ok || repo == nil {
			return
		}
		marked, err := repo.HasCyberPolicyUser(ctx, userID)
		if err != nil {
			slog.Warn("openai.cyber_policy_user_marker_db_read_failed", "user_id", userID, "error", err)
			return
		}
		state.marked = marked
		if cache, ok := s.cache.(CyberPolicyUserMarkerCache); ok && cache != nil {
			if err := cache.SetCyberPolicyUserMark(ctx, userID, marked); err != nil {
				slog.Warn("openai.cyber_policy_user_marker_cache_write_failed", "user_id", userID, "error", err)
			}
		}
	})
	return state.resolved && state.marked
}

func (s *OpenAIGatewayService) shouldSkipCyberPolicyUserAccount(ctx context.Context, account *Account) bool {
	return account != nil && account.IsOpenAICyberPolicyUserBlockingEnabled() && s.cyberPolicyUserMarked(ctx)
}

func (s *OpenAIGatewayService) filterCyberPolicyUserAccounts(ctx context.Context, accounts []Account) []Account {
	if len(accounts) == 0 {
		return accounts
	}
	out := accounts[:0]
	for i := range accounts {
		if s.shouldSkipCyberPolicyUserAccount(ctx, &accounts[i]) {
			continue
		}
		out = append(out, accounts[i])
	}
	return out
}

func (s *ContentModerationService) markCyberPolicyUser(ctx context.Context, userID int64, at time.Time) {
	if s == nil || userID <= 0 || s.userRepo == nil {
		return
	}
	repo, ok := s.userRepo.(CyberPolicyUserMarkerRepository)
	if !ok || repo == nil {
		return
	}
	marked, err := repo.MarkCyberPolicyUser(ctx, userID, at)
	if err != nil {
		slog.Warn("content_moderation.cyber_policy_user_marker_write_failed", "user_id", userID, "error", err)
		return
	}
	if !marked {
		return
	}
	if cache, ok := s.hashCache.(CyberPolicyUserMarkerCache); ok && cache != nil {
		if err := cache.SetCyberPolicyUserMark(ctx, userID, true); err != nil {
			slog.Warn("content_moderation.cyber_policy_user_marker_cache_write_failed", "user_id", userID, "error", err)
		}
	}
}
