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
	release        *GitHubRelease
	releases       map[string]*GitHubRelease
	branch         *GitHubBranch
	branches       map[string]*GitHubBranch
	compare        *GitHubCompare
	compares       map[string]*GitHubCompare
	repo           string
	releaseRepos   []string
	branchRepo     string
	branchName     string
	branchRequests []string
}

func (s *updateServiceGitHubClientStub) FetchLatestRelease(_ context.Context, repo string) (*GitHubRelease, error) {
	s.repo = repo
	s.releaseRepos = append(s.releaseRepos, repo)
	if s.releases != nil {
		if release, ok := s.releases[repo]; ok {
			return release, nil
		}
	}
	return s.release, nil
}

func (s *updateServiceGitHubClientStub) FetchBranch(_ context.Context, repo, branch string) (*GitHubBranch, error) {
	s.branchRepo = repo
	s.branchName = branch
	s.branchRequests = append(s.branchRequests, repo+":"+branch)
	if s.branches != nil {
		if branchInfo, ok := s.branches[repo+":"+branch]; ok {
			return branchInfo, nil
		}
	}
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

func (s *updateServiceGitHubClientStub) CompareCommits(_ context.Context, repo, base, head string) (*GitHubCompare, error) {
	if s.compares != nil {
		if compare, ok := s.compares[repo+":"+base+":"+head]; ok {
			return compare, nil
		}
	}
	if s.compare != nil {
		return s.compare, nil
	}
	return &GitHubCompare{Status: "ahead", AheadBy: 1, TotalCommits: 1}, nil
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
	require.Contains(t, client.releaseRepos, githubRepo)
	require.Contains(t, client.releaseRepos, upstreamGithubRepo)
	require.Equal(t, "1459547510/sub2api-so-fy", githubRepo)
	require.Contains(t, client.branchRequests, githubRepo+":"+githubBranch)
	require.Contains(t, client.branchRequests, upstreamGithubRepo+":"+githubBranch)
}

func TestUpdateServiceDetectsForkBranchCommitUpdate(t *testing.T) {
	client := &updateServiceGitHubClientStub{
		release: &GitHubRelease{
			TagName: "v0.1.138",
			Name:    "v0.1.138",
		},
		branches: map[string]*GitHubBranch{
			githubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
				},
			},
			upstreamGithubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				},
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

func TestUpdateServiceDetectsUpstreamUpdateButBlocksOneClickUntilForkRelease(t *testing.T) {
	forkHead := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	upstreamHead := "cccccccccccccccccccccccccccccccccccccccc"
	client := &updateServiceGitHubClientStub{
		releases: map[string]*GitHubRelease{
			githubRepo: {
				TagName: "v0.1.138",
				Name:    "v0.1.138",
			},
			upstreamGithubRepo: {
				TagName: "v0.1.139",
				Name:    "v0.1.139",
			},
		},
		branches: map[string]*GitHubBranch{
			githubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: forkHead,
				},
			},
			upstreamGithubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: upstreamHead,
				},
			},
		},
		compares: map[string]*GitHubCompare{
			upstreamGithubRepo + ":" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + ":" + upstreamHead: {
				Status:       "ahead",
				AheadBy:      1,
				TotalCommits: 1,
				HTMLURL:      "https://github.com/Wei-Shaw/sub2api/compare/base...head",
			},
			githubRepo + ":" + upstreamHead + ":" + forkHead: {
				Status:   "behind",
				BehindBy: 1,
			},
		},
	}
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		client,
		"0.1.138",
		"release",
		forkHead,
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.True(t, info.HasUpdate)
	require.False(t, info.UpdateReady, "upstream-only update must not be installed from an unsynced fork release")
	require.Equal(t, "0.1.139", info.LatestVersion)
	require.Equal(t, "0.1.138", info.ForkLatestVersion)
	require.NotNil(t, info.UpstreamInfo)
	require.True(t, info.UpstreamInfo.HasUpdate)
	require.True(t, info.UpstreamInfo.HasNewVersion)
	require.True(t, info.UpstreamInfo.SyncRequired)
}

func TestUpdateServiceAllowsForkReleaseWhenItContainsUpstreamHead(t *testing.T) {
	upstreamHead := "cccccccccccccccccccccccccccccccccccccccc"
	forkHead := "dddddddddddddddddddddddddddddddddddddddd"
	client := &updateServiceGitHubClientStub{
		releases: map[string]*GitHubRelease{
			githubRepo: {
				TagName: "v0.1.139",
				Name:    "v0.1.139",
			},
			upstreamGithubRepo: {
				TagName: "v0.1.139",
				Name:    "v0.1.139",
			},
		},
		branches: map[string]*GitHubBranch{
			githubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: forkHead,
				},
			},
			upstreamGithubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: upstreamHead,
				},
			},
		},
		compares: map[string]*GitHubCompare{
			upstreamGithubRepo + ":" + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" + ":" + upstreamHead: {
				Status:       "ahead",
				AheadBy:      1,
				TotalCommits: 1,
			},
			githubRepo + ":" + upstreamHead + ":" + forkHead: {
				Status:  "ahead",
				AheadBy: 2,
			},
		},
	}
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		client,
		"0.1.138",
		"release",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.True(t, info.HasUpdate)
	require.True(t, info.UpdateReady)
	require.NotNil(t, info.UpstreamInfo)
	require.False(t, info.UpstreamInfo.SyncRequired)
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
