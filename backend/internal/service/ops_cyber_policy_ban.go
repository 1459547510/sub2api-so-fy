package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const cyberPolicyRevocationLookback = time.Duration(defaultContentModerationViolationWindowHours) * time.Hour

func (s *OpsService) applyCyberPolicyRevocationBan(ctx context.Context, entry *OpsInsertErrorLogInput) {
	if s == nil || entry == nil || s.opsRepo == nil || s.accountRepo == nil || s.userRepo == nil {
		return
	}

	for _, accountID := range tokenRevokedOpenAIAccountIDs(entry) {
		account, err := s.accountRepo.GetByID(ctx, accountID)
		if err != nil || account == nil {
			if err != nil {
				slog.Warn("cyber_policy_revocation_ban_account_lookup_failed", "account_id", accountID, "error", err)
			}
			continue
		}

		authAccount := account
		if resolved, resolveErr := resolveCredentialAccount(ctx, s.accountRepo, account); resolveErr == nil && resolved != nil {
			authAccount = resolved
		}
		if !authAccount.IsOpenAIProOrProLite() || authAccount.Status != StatusError || !isTokenRevokedAccountError(authAccount.ErrorMessage) {
			continue
		}

		revokedAt := entry.CreatedAt
		if revokedAt.IsZero() {
			revokedAt = time.Now()
		}
		candidate, err := s.opsRepo.FindCyberPolicyBanCandidate(ctx, accountID, revokedAt.Add(-cyberPolicyRevocationLookback), revokedAt)
		if err != nil {
			slog.Warn("cyber_policy_revocation_ban_candidate_lookup_failed", "account_id", accountID, "error", err)
			continue
		}
		if candidate == nil || candidate.UserID <= 0 {
			continue
		}

		user, err := s.userRepo.GetByID(ctx, candidate.UserID)
		if err != nil || user == nil {
			if err != nil {
				slog.Warn("cyber_policy_revocation_ban_user_lookup_failed", "account_id", accountID, "user_id", candidate.UserID, "error", err)
			}
			continue
		}
		if user.IsAdmin() {
			slog.Warn("cyber_policy_revocation_ban_skipped_admin", "account_id", accountID, "user_id", user.ID)
			continue
		}
		if user.Status == StatusDisabled {
			continue
		}

		user.Status = StatusDisabled
		if err := s.userRepo.Update(ctx, user, UserUpdateFields{Status: true}); err != nil {
			slog.Warn("cyber_policy_revocation_ban_update_failed", "account_id", accountID, "user_id", user.ID, "error", err)
			continue
		}
		if s.apiKeyService != nil {
			s.apiKeyService.InvalidateAuthCacheByUserID(ctx, user.ID)
		}
		if err := s.recordCyberPolicyRevocationAudit(ctx, entry, accountID, authAccount, user, candidate, revokedAt); err != nil {
			slog.Warn("cyber_policy_revocation_ban_audit_failed", "account_id", accountID, "user_id", user.ID, "error", err)
		}
		slog.Warn(
			"cyber_policy_revocation_user_disabled",
			"account_id", accountID,
			"credential_account_id", authAccount.ID,
			"user_id", user.ID,
			"hit_count", candidate.HitCount,
			"last_hit_at", candidate.LastHitAt,
		)
	}
}

func (s *OpsService) recordCyberPolicyRevocationAudit(
	ctx context.Context,
	revocation *OpsInsertErrorLogInput,
	accountID int64,
	credentialAccount *Account,
	user *User,
	candidate *CyberPolicyBanCandidate,
	revokedAt time.Time,
) error {
	if s == nil || s.auditLogService == nil || revocation == nil || credentialAccount == nil || user == nil || candidate == nil {
		return nil
	}

	requestBody := ""
	if excerpt := strings.TrimSpace(candidate.InputExcerpt); excerpt != "" {
		if encoded, err := json.Marshal(map[string]string{"input_excerpt": excerpt}); err == nil {
			requestBody = RedactAuditBody(encoded, "application/json")
		}
	}

	extra := map[string]any{
		"rule":                         AuditActionCyberPolicyRevocationBan,
		"outcome":                      "user_disabled",
		"revoked_account_id":           accountID,
		"credential_account_id":        credentialAccount.ID,
		"account_plan_type":            strings.TrimSpace(credentialAccount.GetCredential("plan_type")),
		"hit_count":                    candidate.HitCount,
		"trigger_ops_error_log_id":     candidate.OpsErrorLogID,
		"trigger_request_id":           candidate.RequestID,
		"trigger_client_request_id":    candidate.ClientRequestID,
		"triggered_at":                 candidate.LastHitAt.UTC().Format(time.RFC3339Nano),
		"trigger_model":                candidate.Model,
		"trigger_request_path":         candidate.RequestPath,
		"revocation_request_id":        strings.TrimSpace(revocation.RequestID),
		"revocation_client_request_id": strings.TrimSpace(revocation.ClientRequestID),
		"revoked_at":                   revokedAt.UTC().Format(time.RFC3339Nano),
		"disabled_user_id":             user.ID,
		"disabled_user_email":          user.Email,
	}
	if candidate.APIKeyID != nil {
		extra["trigger_api_key_id"] = *candidate.APIKeyID
	}

	path := strings.TrimSpace(candidate.RequestPath)
	if path == "" {
		path = "/security-audit/cyber-policy-revocation"
	}
	actorUserID := user.ID
	return s.auditLogService.RecordSync(ctx, &AuditLog{
		CreatedAt:        revokedAt.UTC(),
		ActorUserID:      &actorUserID,
		ActorEmail:       user.Email,
		ActorRole:        user.Role,
		AuthMethod:       "api_key",
		CredentialMasked: candidate.APIKeyPrefix,
		Action:           AuditActionCyberPolicyRevocationBan,
		Method:           "SYSTEM",
		Path:             path,
		RequestID:        candidate.RequestID,
		ClientIP:         candidate.ClientIP,
		UserAgent:        candidate.UserAgent,
		RequestBody:      requestBody,
		StatusCode:       http.StatusOK,
		Extra:            extra,
	})
}

func tokenRevokedOpenAIAccountIDs(entry *OpsInsertErrorLogInput) []int64 {
	if entry == nil {
		return nil
	}
	seen := make(map[int64]struct{})
	accountIDs := make([]int64, 0, 1)
	add := func(accountID int64, platform string, statusCode int) {
		if accountID <= 0 || statusCode != http.StatusUnauthorized || !strings.EqualFold(strings.TrimSpace(platform), PlatformOpenAI) {
			return
		}
		if _, exists := seen[accountID]; exists {
			return
		}
		seen[accountID] = struct{}{}
		accountIDs = append(accountIDs, accountID)
	}

	if entry.AccountID != nil && entry.UpstreamStatusCode != nil {
		add(*entry.AccountID, entry.Platform, *entry.UpstreamStatusCode)
	}
	if entry.UpstreamErrorsJSON == nil {
		return accountIDs
	}
	events, err := ParseOpsUpstreamErrors(*entry.UpstreamErrorsJSON)
	if err != nil {
		return accountIDs
	}
	for _, event := range events {
		if event == nil {
			continue
		}
		platform := event.Platform
		if strings.TrimSpace(platform) == "" {
			platform = entry.Platform
		}
		add(event.AccountID, platform, event.UpstreamStatusCode)
	}
	return accountIDs
}

func isTokenRevokedAccountError(message string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(message)), "token revoked (401):")
}
