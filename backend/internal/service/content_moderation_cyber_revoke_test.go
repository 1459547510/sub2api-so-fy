package service

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type contentModerationRevokeGroupRepo struct {
	groupRepoNoop
	group *Group
}

func (r *contentModerationRevokeGroupRepo) GetByIDLite(context.Context, int64) (*Group, error) {
	return r.group, nil
}

func TestContentModerationConfig_CyberPolicyRevokeGroupDefaultsDisabled(t *testing.T) {
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{}},
		nil, nil, nil, nil, nil, nil, nil,
	)

	view, err := svc.GetConfig(context.Background())

	require.NoError(t, err)
	require.Zero(t, view.CyberPolicyRevokeGroupID)
}

func TestContentModerationConfig_SavesValidCyberPolicyRevokeGroup(t *testing.T) {
	settingRepo := &contentModerationTestSettingRepo{values: map[string]string{}}
	groupRepo := &contentModerationRevokeGroupRepo{group: &Group{
		ID:          42,
		Platform:    PlatformOpenAI,
		Status:      StatusActive,
		IsExclusive: true,
	}}
	svc := NewContentModerationService(settingRepo, nil, nil, groupRepo, nil, nil, nil, nil)
	groupID := int64(42)

	view, err := svc.UpdateConfig(context.Background(), UpdateContentModerationConfigInput{
		CyberPolicyRevokeGroupID: &groupID,
	})

	require.NoError(t, err)
	require.Equal(t, groupID, view.CyberPolicyRevokeGroupID)
	require.True(t, strings.Contains(settingRepo.values[SettingKeyContentModerationConfig], `"cyber_policy_revoke_group_id":42`))
	reloaded, err := svc.GetConfig(context.Background())
	require.NoError(t, err)
	require.Equal(t, groupID, reloaded.CyberPolicyRevokeGroupID)
}

func TestRecordCyberPolicyEvent_RemovesConfiguredExclusiveGroupForRegularUser(t *testing.T) {
	repo := &contentModerationTestRepo{}
	userRepo := &contentModerationTestUserRepo{user: &User{ID: 7, Role: RoleUser, AllowedGroups: []int64{42, 99}}}
	invalidator := &contentModerationTestAuthCacheInvalidator{}
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: `{"cyber_policy_revoke_group_id":42,"cyber_policy_exclude_from_ban_count":true}`,
		}},
		repo, nil, nil, userRepo, nil, invalidator, nil,
	)

	svc.RecordCyberPolicyEvent(context.Background(), CyberPolicyRecordInput{UserID: 7, UserEmail: "u@example.com"})

	require.Equal(t, []struct{ userID, groupID int64 }{{userID: 7, groupID: 42}}, userRepo.removed)
	require.Equal(t, []int64{99}, userRepo.user.AllowedGroups)
	require.Equal(t, []int64{7}, invalidator.userIDs)
	require.Len(t, repo.snapshotLogs(), 1)
}

func TestRecordCyberPolicyEvent_AdministratorsKeepConfiguredExclusiveGroup(t *testing.T) {
	repo := &contentModerationTestRepo{}
	userRepo := &contentModerationTestUserRepo{user: &User{ID: 8, Role: RoleAdmin, AllowedGroups: []int64{42}}}
	invalidator := &contentModerationTestAuthCacheInvalidator{}
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: `{"cyber_policy_revoke_group_id":42,"cyber_policy_exclude_from_ban_count":true}`,
		}},
		repo, nil, nil, userRepo, nil, invalidator, nil,
	)

	svc.RecordCyberPolicyEvent(context.Background(), CyberPolicyRecordInput{UserID: 8, UserEmail: "admin@example.com"})

	require.Empty(t, userRepo.removed)
	require.Equal(t, []int64{42}, userRepo.user.AllowedGroups)
	require.Empty(t, invalidator.userIDs)
	require.Len(t, repo.snapshotLogs(), 1)
}
