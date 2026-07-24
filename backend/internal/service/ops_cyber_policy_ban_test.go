package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type cyberPolicyBanAccountRepo struct {
	AccountRepository
	accounts map[int64]*Account
}

func (r *cyberPolicyBanAccountRepo) GetByID(_ context.Context, id int64) (*Account, error) {
	account := r.accounts[id]
	if account == nil {
		return nil, ErrAccountNotFound
	}
	clone := *account
	return &clone, nil
}

type cyberPolicyBanUserRepo struct {
	UserRepository
	users       map[int64]*User
	updateCount int
}

type cyberPolicyBanAuditRepo struct {
	AuditLogRepository
	logs []*AuditLog
}

func (r *cyberPolicyBanAuditRepo) Insert(_ context.Context, log *AuditLog) error {
	clone := *log
	r.logs = append(r.logs, &clone)
	return nil
}

func (r *cyberPolicyBanUserRepo) GetByID(_ context.Context, id int64) (*User, error) {
	user := r.users[id]
	if user == nil {
		return nil, ErrUserNotFound
	}
	clone := *user
	return &clone, nil
}

func (r *cyberPolicyBanUserRepo) Update(_ context.Context, user *User) error {
	clone := *user
	r.users[user.ID] = &clone
	r.updateCount++
	return nil
}

func TestOpsServiceTokenRevokedDisablesTopCyberPolicyUser(t *testing.T) {
	createdAt := time.Date(2026, 7, 24, 7, 51, 58, 0, time.UTC)
	lastHit := createdAt.Add(-22 * time.Minute)
	var (
		candidateCalls int
		gotSince       time.Time
		gotUntil       time.Time
	)
	opsRepo := &opsRepoMock{
		InsertErrorLogFn: func(context.Context, *OpsInsertErrorLogInput) (int64, error) { return 1, nil },
		FindCyberPolicyBanCandidateFn: func(_ context.Context, accountID int64, since, until time.Time) (*CyberPolicyBanCandidate, error) {
			candidateCalls++
			require.Equal(t, int64(1714), accountID)
			gotSince, gotUntil = since, until
			apiKeyID := int64(88)
			return &CyberPolicyBanCandidate{
				UserID:          147,
				HitCount:        3,
				LastHitAt:       lastHit,
				OpsErrorLogID:   9001,
				RequestID:       "req-cyber",
				ClientRequestID: "client-cyber",
				APIKeyID:        &apiKeyID,
				APIKeyPrefix:    "sk-test1",
				ClientIP:        "203.0.113.9",
				UserAgent:       "claude-code/1.0",
				Model:           "gpt-5.5",
				RequestPath:     "/v1/responses",
				InputExcerpt:    "retained prompt excerpt",
			}, nil
		},
	}
	accountRepo := &cyberPolicyBanAccountRepo{accounts: map[int64]*Account{
		1714: {ID: 1714, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusError, ErrorMessage: "Token revoked (401): token has been revoked", Credentials: map[string]any{"plan_type": "pro"}},
	}}
	userRepo := &cyberPolicyBanUserRepo{users: map[int64]*User{
		147: {ID: 147, Email: "user@example.com", Role: RoleUser, Status: StatusActive},
	}}
	auditRepo := &cyberPolicyBanAuditRepo{}
	svc := NewOpsService(opsRepo, nil, nil, accountRepo, userRepo, nil, nil, nil, nil, nil, nil)
	svc.SetAuditLogService(NewAuditLogService(auditRepo, nil))

	entry := revokedOpenAIErrorEntry(createdAt, 1714)
	entry.RequestID = "req-revocation"
	entry.ClientRequestID = "client-revocation"
	err := svc.RecordError(context.Background(), entry)

	require.NoError(t, err)
	require.Equal(t, 1, candidateCalls)
	require.Equal(t, createdAt, gotUntil)
	require.Equal(t, createdAt.Add(-cyberPolicyRevocationLookback), gotSince)
	require.Equal(t, StatusDisabled, userRepo.users[147].Status)
	require.Equal(t, 1, userRepo.updateCount)
	require.Len(t, auditRepo.logs, 1)
	audit := auditRepo.logs[0]
	require.Equal(t, AuditActionCyberPolicyRevocationBan, audit.Action)
	require.Equal(t, "SYSTEM", audit.Method)
	require.Equal(t, "req-cyber", audit.RequestID)
	require.Equal(t, "user@example.com", audit.ActorEmail)
	require.Equal(t, "203.0.113.9", audit.ClientIP)
	require.JSONEq(t, `{"input_excerpt":"retained prompt excerpt"}`, audit.RequestBody)
	require.EqualValues(t, 1714, audit.Extra["revoked_account_id"])
	require.EqualValues(t, 3, audit.Extra["hit_count"])
	require.Equal(t, "req-revocation", audit.Extra["revocation_request_id"])
	require.Equal(t, "req-cyber", audit.Extra["trigger_request_id"])
}

func TestOpsServiceGenericOpenAI401DoesNotDisableCyberPolicyUser(t *testing.T) {
	createdAt := time.Date(2026, 7, 24, 8, 0, 0, 0, time.UTC)
	candidateCalls := 0
	opsRepo := &opsRepoMock{
		InsertErrorLogFn: func(context.Context, *OpsInsertErrorLogInput) (int64, error) { return 1, nil },
		FindCyberPolicyBanCandidateFn: func(context.Context, int64, time.Time, time.Time) (*CyberPolicyBanCandidate, error) {
			candidateCalls++
			return &CyberPolicyBanCandidate{UserID: 147, HitCount: 3}, nil
		},
	}
	accountRepo := &cyberPolicyBanAccountRepo{accounts: map[int64]*Account{
		1714: {ID: 1714, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusError, ErrorMessage: "Unauthorized (401): account authentication failed permanently", Credentials: map[string]any{"plan_type": "pro"}},
	}}
	userRepo := &cyberPolicyBanUserRepo{users: map[int64]*User{
		147: {ID: 147, Role: RoleUser, Status: StatusActive},
	}}
	svc := NewOpsService(opsRepo, nil, nil, accountRepo, userRepo, nil, nil, nil, nil, nil, nil)

	require.NoError(t, svc.RecordError(context.Background(), revokedOpenAIErrorEntry(createdAt, 1714)))
	require.Zero(t, candidateCalls)
	require.Equal(t, StatusActive, userRepo.users[147].Status)
	require.Zero(t, userRepo.updateCount)
}

func TestOpsServiceTokenRevokedBanSkipsAdminAndIsIdempotent(t *testing.T) {
	createdAt := time.Date(2026, 7, 24, 8, 0, 0, 0, time.UTC)
	opsRepo := &opsRepoMock{
		InsertErrorLogFn: func(context.Context, *OpsInsertErrorLogInput) (int64, error) { return 1, nil },
		FindCyberPolicyBanCandidateFn: func(context.Context, int64, time.Time, time.Time) (*CyberPolicyBanCandidate, error) {
			return &CyberPolicyBanCandidate{UserID: 1, HitCount: 10, LastHitAt: createdAt.Add(-time.Minute)}, nil
		},
	}
	accountRepo := &cyberPolicyBanAccountRepo{accounts: map[int64]*Account{
		452: {ID: 452, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusError, ErrorMessage: "Token revoked (401): revoked", Credentials: map[string]any{"plan_type": "pro"}},
	}}
	userRepo := &cyberPolicyBanUserRepo{users: map[int64]*User{
		1: {ID: 1, Role: RoleAdmin, Status: StatusActive},
	}}
	svc := NewOpsService(opsRepo, nil, nil, accountRepo, userRepo, nil, nil, nil, nil, nil, nil)

	require.NoError(t, svc.RecordError(context.Background(), revokedOpenAIErrorEntry(createdAt, 452)))
	require.Equal(t, StatusActive, userRepo.users[1].Status)
	require.Zero(t, userRepo.updateCount)

	userRepo.users[1] = &User{ID: 1, Role: RoleUser, Status: StatusActive}
	require.NoError(t, svc.RecordError(context.Background(), revokedOpenAIErrorEntry(createdAt, 452)))
	require.NoError(t, svc.RecordError(context.Background(), revokedOpenAIErrorEntry(createdAt, 452)))
	require.Equal(t, StatusDisabled, userRepo.users[1].Status)
	require.Equal(t, 1, userRepo.updateCount)
}

func TestOpsServiceTokenRevokedPlusAccountDoesNotDisableCyberPolicyUser(t *testing.T) {
	createdAt := time.Date(2026, 7, 24, 8, 0, 0, 0, time.UTC)
	candidateCalls := 0
	opsRepo := &opsRepoMock{
		InsertErrorLogFn: func(context.Context, *OpsInsertErrorLogInput) (int64, error) { return 1, nil },
		FindCyberPolicyBanCandidateFn: func(context.Context, int64, time.Time, time.Time) (*CyberPolicyBanCandidate, error) {
			candidateCalls++
			return &CyberPolicyBanCandidate{UserID: 147, HitCount: 3}, nil
		},
	}
	accountRepo := &cyberPolicyBanAccountRepo{accounts: map[int64]*Account{
		1714: {ID: 1714, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusError, ErrorMessage: "Token revoked (401): revoked", Credentials: map[string]any{"plan_type": "plus"}},
	}}
	userRepo := &cyberPolicyBanUserRepo{users: map[int64]*User{
		147: {ID: 147, Role: RoleUser, Status: StatusActive},
	}}
	svc := NewOpsService(opsRepo, nil, nil, accountRepo, userRepo, nil, nil, nil, nil, nil, nil)

	require.NoError(t, svc.RecordError(context.Background(), revokedOpenAIErrorEntry(createdAt, 1714)))
	require.Zero(t, candidateCalls)
	require.Equal(t, StatusActive, userRepo.users[147].Status)
	require.Zero(t, userRepo.updateCount)
}

func TestAccountIsOpenAIProOrProLite(t *testing.T) {
	tests := []struct {
		name string
		kind string
		plan string
		want bool
	}{
		{name: "pro", kind: AccountTypeOAuth, plan: "pro", want: true},
		{name: "chatgpt pro", kind: AccountTypeOAuth, plan: " ChatGPT Pro ", want: true},
		{name: "pro lite", kind: AccountTypeOAuth, plan: "pro_lite", want: true},
		{name: "chatgpt pro lite", kind: AccountTypeOAuth, plan: "chatgpt-pro-lite", want: true},
		{name: "plus", kind: AccountTypeOAuth, plan: "plus", want: false},
		{name: "api key", kind: AccountTypeAPIKey, plan: "pro", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &Account{
				Platform:    PlatformOpenAI,
				Type:        tt.kind,
				Credentials: map[string]any{"plan_type": tt.plan},
			}
			require.Equal(t, tt.want, account.IsOpenAIProOrProLite())
		})
	}
}

func revokedOpenAIErrorEntry(createdAt time.Time, accountID int64) *OpsInsertErrorLogInput {
	return &OpsInsertErrorLogInput{
		Platform:  PlatformOpenAI,
		CreatedAt: createdAt,
		UpstreamErrors: []*OpsUpstreamErrorEvent{{
			Platform:           PlatformOpenAI,
			AccountID:          accountID,
			UpstreamStatusCode: 401,
			Message:            "token has been revoked",
		}},
	}
}
