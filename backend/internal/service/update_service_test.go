//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type updateServiceCacheStub struct {
	data string
}

func (s *updateServiceCacheStub) GetUpdateInfo(context.Context) (string, error) {
	if s.data == "" {
		return "", errors.New("cache miss")
	}
	return s.data, nil
}

func (s *updateServiceCacheStub) SetUpdateInfo(_ context.Context, data string, _ time.Duration) error {
	s.data = data
	return nil
}

type updateServiceGitHubClientStub struct {
	release    *GitHubRelease
	branch     *GitHubBranch
	repo       string
	branchRepo string
	branchName string
}

func (s *updateServiceGitHubClientStub) FetchLatestRelease(_ context.Context, repo string) (*GitHubRelease, error) {
	s.repo = repo
	return s.release, nil
}

func (s *updateServiceGitHubClientStub) FetchBranch(_ context.Context, repo, branch string) (*GitHubBranch, error) {
	s.branchRepo = repo
	s.branchName = branch
	if s.branch != nil {
		return s.branch, nil
	}
	return &GitHubBranch{
		Name: branch,
		Commit: GitHubCommitRef{
			SHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
	}, nil
}

func (s *updateServiceGitHubClientStub) DownloadFile(context.Context, string, string, int64) error {
	panic("DownloadFile should not be called when no update is available")
}

func (s *updateServiceGitHubClientStub) FetchChecksumFile(context.Context, string) ([]byte, error) {
	panic("FetchChecksumFile should not be called when no update is available")
}

func TestUpdateServicePerformUpdateNoUpdateReturnsSentinel(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.132",
				Name:    "v0.1.132",
			},
		},
		"0.1.132",
		"release",
	)

	err := svc.PerformUpdate(context.Background())

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNoUpdateAvailable))
	require.ErrorIs(t, err, ErrNoUpdateAvailable)
}

func TestUpdateServiceUsesForkReleaseRepository(t *testing.T) {
	client := &updateServiceGitHubClientStub{
		release: &GitHubRelease{
			TagName: "v0.1.132",
			Name:    "v0.1.132",
		},
	}
	svc := NewUpdateService(&updateServiceCacheStub{}, client, "0.1.132", "release")

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.False(t, info.HasUpdate)
	require.Equal(t, githubRepo, client.repo)
	require.Equal(t, "1459547510/sub2api-so-fy", client.repo)
	require.Equal(t, githubRepo, client.branchRepo)
	require.Equal(t, githubBranch, client.branchName)
}

func TestUpdateServiceDetectsForkBranchCommitUpdate(t *testing.T) {
	client := &updateServiceGitHubClientStub{
		release: &GitHubRelease{
			TagName: "v0.1.138",
			Name:    "v0.1.138",
		},
		branch: &GitHubBranch{
			Name: githubBranch,
			Commit: GitHubCommitRef{
				SHA: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			},
		},
	}
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		client,
		"0.1.138",
		"release",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.False(t, info.HasUpdate, "release update should remain version/tag based")
	require.NotNil(t, info.BranchInfo)
	require.True(t, info.BranchInfo.HasNewCommit)
	require.True(t, info.BranchInfo.CanCompare)
	require.Equal(t, "behind", info.BranchInfo.Status)
	require.Equal(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", info.BranchInfo.CurrentCommit)
	require.Equal(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", info.BranchInfo.LatestCommit)
	require.Contains(t, info.BranchInfo.CompareURL, "/compare/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa...bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
}

func TestUpdateServiceDoesNotFlagSameCommitPrefix(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{},
		"0.1.138",
		"release",
		"abcdef1",
	)

	info := svc.branchInfoForLatestCommit("abcdef1234567890abcdef1234567890abcdef12", "abcdef1")

	require.False(t, info.HasNewCommit)
	require.Equal(t, "current", info.Status)
	require.True(t, info.CanCompare)
}
