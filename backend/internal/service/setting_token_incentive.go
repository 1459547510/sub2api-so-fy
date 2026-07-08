package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// IsTokenIncentiveEnabled 检查是否启用每周 Token 激励计划。
func (s *SettingService) IsTokenIncentiveEnabled(ctx context.Context) bool {
	value, err := s.settingRepo.GetValue(ctx, SettingKeyTokenIncentiveEnabled)
	if err != nil {
		return false
	}
	return value == "true"
}

func (s *SettingService) GetTokenIncentiveRules(ctx context.Context) []TokenIncentiveRule {
	if s == nil || s.settingRepo == nil {
		return DefaultTokenIncentiveRules()
	}
	raw, err := s.settingRepo.GetValue(ctx, SettingKeyTokenIncentiveRules)
	if err != nil {
		return DefaultTokenIncentiveRules()
	}
	rules, err := ParseTokenIncentiveRules(raw)
	if err != nil {
		slog.Warn("invalid token incentive rules setting, fallback to default", "error", err)
		return DefaultTokenIncentiveRules()
	}
	return rules
}

func ParseTokenIncentiveRules(raw string) ([]TokenIncentiveRule, error) {
	if strings.TrimSpace(raw) == "" {
		return DefaultTokenIncentiveRules(), nil
	}
	var rules []TokenIncentiveRule
	if err := json.Unmarshal([]byte(raw), &rules); err != nil {
		return nil, infraerrors.BadRequest("INVALID_TOKEN_INCENTIVE_RULES", "token incentive rules must be a JSON array").WithCause(err)
	}
	return NormalizeTokenIncentiveRules(rules)
}

func mustMarshalDefaultTokenIncentiveRules() string {
	data, err := json.Marshal(DefaultTokenIncentiveRules())
	if err != nil {
		panic(fmt.Sprintf("marshal default token incentive rules: %v", err))
	}
	return string(data)
}
