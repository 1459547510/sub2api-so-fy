package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/stretchr/testify/require"
)

type cyberPolicyMarkerUserRepoStub struct {
	UserRepository
	marked      bool
	markAllowed bool
	markCalls   int
	hasCalls    int
	markErr     error
	hasErr      error
}

func (r *cyberPolicyMarkerUserRepoStub) MarkCyberPolicyUser(context.Context, int64, time.Time) (bool, error) {
	r.markCalls++
	return r.markAllowed, r.markErr
}

func (r *cyberPolicyMarkerUserRepoStub) HasCyberPolicyUser(context.Context, int64) (bool, error) {
	r.hasCalls++
	return r.marked, r.hasErr
}

type cyberPolicyMarkerCacheStub struct {
	GatewayCache
	marked    bool
	found     bool
	getErr    error
	setCalls  int
	setMarked bool
}

func (c *cyberPolicyMarkerCacheStub) GetCyberPolicyUserMark(context.Context, int64) (bool, bool, error) {
	return c.marked, c.found, c.getErr
}

func (c *cyberPolicyMarkerCacheStub) SetCyberPolicyUserMark(context.Context, int64, bool) error {
	c.setCalls++
	c.setMarked = true
	return nil
}

type cyberPolicyMarkerHashCacheStub struct {
	ContentModerationHashCache
	setCalls int
}

func (c *cyberPolicyMarkerHashCacheStub) GetCyberPolicyUserMark(context.Context, int64) (bool, bool, error) {
	return false, false, nil
}

func (c *cyberPolicyMarkerHashCacheStub) SetCyberPolicyUserMark(context.Context, int64, bool) error {
	c.setCalls++
	return nil
}

func TestOpenAICyberPolicyUserFilter_SkipsMarkedUserOnlyForProtectedAccounts(t *testing.T) {
	repo := &cyberPolicyMarkerUserRepoStub{marked: true, markAllowed: true}
	svc := &OpenAIGatewayService{userRepo: repo}
	ctx := withCyberPolicyUserFilterState(context.WithValue(context.Background(), ctxkey.UserID, int64(42)))
	accounts := []Account{
		{ID: 1, Platform: PlatformOpenAI, Extra: map[string]any{OpenAICyberPolicyUserBlockingExtraKey: true}},
		{ID: 2, Platform: PlatformOpenAI},
	}

	filtered := svc.filterCyberPolicyUserAccounts(ctx, accounts)
	require.Equal(t, []int64{2}, []int64{filtered[0].ID})
	require.Equal(t, 1, repo.hasCalls)
}

func TestOpenAICyberPolicyUserFilter_DefaultOffAndAdminExempt(t *testing.T) {
	repo := &cyberPolicyMarkerUserRepoStub{marked: true, markAllowed: true}
	svc := &OpenAIGatewayService{userRepo: repo}
	account := Account{ID: 1, Platform: PlatformOpenAI, Extra: map[string]any{OpenAICyberPolicyUserBlockingExtraKey: true}}

	regular := withCyberPolicyUserFilterState(context.WithValue(context.Background(), ctxkey.UserID, int64(42)))
	admin := context.WithValue(regular, ctxkey.UserRole, RoleAdmin)
	noToggle := Account{ID: 2, Platform: PlatformOpenAI}
	filtered := svc.filterCyberPolicyUserAccounts(regular, []Account{noToggle})
	require.Len(t, filtered, 1)
	require.False(t, svc.shouldSkipCyberPolicyUserAccount(admin, &account))
	require.Equal(t, 0, repo.hasCalls)
}

func TestOpenAICyberPolicyUserFilter_CacheHitAndFailureFailOpen(t *testing.T) {
	repo := &cyberPolicyMarkerUserRepoStub{marked: true, markAllowed: true}
	cache := &cyberPolicyMarkerCacheStub{marked: true, found: true}
	svc := &OpenAIGatewayService{userRepo: repo, cache: cache}
	ctx := withCyberPolicyUserFilterState(context.WithValue(context.Background(), ctxkey.UserID, int64(42)))
	account := Account{ID: 1, Platform: PlatformOpenAI, Extra: map[string]any{OpenAICyberPolicyUserBlockingExtraKey: true}}
	require.True(t, svc.shouldSkipCyberPolicyUserAccount(ctx, &account))
	require.Equal(t, 0, repo.hasCalls)

	cacheFailure := &cyberPolicyMarkerCacheStub{marked: true, getErr: errors.New("redis down")}
	cacheFailureSvc := &OpenAIGatewayService{userRepo: repo, cache: cacheFailure}
	cacheFailureCtx := withCyberPolicyUserFilterState(context.WithValue(context.Background(), ctxkey.UserID, int64(42)))
	require.False(t, cacheFailureSvc.shouldSkipCyberPolicyUserAccount(cacheFailureCtx, &account))

	failingRepo := &cyberPolicyMarkerUserRepoStub{hasErr: errors.New("db down")}
	failingSvc := &OpenAIGatewayService{userRepo: failingRepo}
	failingCtx := withCyberPolicyUserFilterState(context.WithValue(context.Background(), ctxkey.UserID, int64(42)))
	require.False(t, failingSvc.shouldSkipCyberPolicyUserAccount(failingCtx, &account))
}

func TestOpenAICyberPolicyUserFilter_RecheckRejectsProtectedAccount(t *testing.T) {
	svc := &OpenAIGatewayService{userRepo: &cyberPolicyMarkerUserRepoStub{marked: true, markAllowed: true}}
	ctx := withCyberPolicyUserFilterState(context.WithValue(context.Background(), ctxkey.UserID, int64(42)))
	account := &Account{ID: 1, Platform: PlatformOpenAI, Status: StatusActive, Schedulable: true, Extra: map[string]any{OpenAICyberPolicyUserBlockingExtraKey: true}}

	require.Nil(t, svc.recheckSelectedOpenAIAccountFromDB(ctx, account, nil, PlatformOpenAI, "", false, ""))
}

func TestContentModerationService_MarkCyberPolicyUser(t *testing.T) {
	repo := &cyberPolicyMarkerUserRepoStub{marked: true, markAllowed: true}
	cache := &cyberPolicyMarkerHashCacheStub{}
	svc := &ContentModerationService{userRepo: repo, hashCache: cache}

	svc.markCyberPolicyUser(context.Background(), 7, time.Unix(100, 0))
	require.Equal(t, 1, repo.markCalls)
	require.Equal(t, 1, cache.setCalls)

	repo.marked = false
	repo.markAllowed = false
	svc.markCyberPolicyUser(context.Background(), 8, time.Unix(101, 0))
	require.Equal(t, 1, cache.setCalls)

	repo.markErr = errors.New("write failed")
	svc.markCyberPolicyUser(context.Background(), 9, time.Unix(102, 0))
	require.Equal(t, 3, repo.markCalls)
}

func TestRecordCyberPolicyEvent_MarksRegularUser(t *testing.T) {
	repo := &contentModerationTestRepo{}
	userRepo := &cyberPolicyMarkerUserRepoStub{marked: true, markAllowed: true}
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{SettingKeyRiskControlEnabled: "true"}},
		repo, nil, nil, userRepo, nil, nil, nil,
	)

	svc.RecordCyberPolicyEvent(context.Background(), CyberPolicyRecordInput{UserID: 7, UserEmail: "u@example.com"})

	require.Equal(t, 1, userRepo.markCalls)
}

func TestRecordCyberPolicyEvent_DoesNotMarkAdmin(t *testing.T) {
	repo := &contentModerationTestRepo{}
	userRepo := &cyberPolicyMarkerUserRepoStub{markAllowed: false}
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{SettingKeyRiskControlEnabled: "true"}},
		repo, nil, nil, userRepo, nil, nil, nil,
	)

	svc.RecordCyberPolicyEvent(context.Background(), CyberPolicyRecordInput{UserID: 8, UserEmail: "admin@example.com"})

	require.Equal(t, 1, userRepo.markCalls)
}
