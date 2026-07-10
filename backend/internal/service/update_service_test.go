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
	recentReleases []*GitHubRelease
	recentErr      error
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

func (s *updateServiceGitHubClientStub) FetchRecentReleases(context.Context, string, int) ([]*GitHubRelease, error) {
	return s.recentReleases, s.recentErr
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
	require.NotContains(t, client.branchRequests, upstreamGithubRepo+":"+githubBranch)
}

func TestCompareVersionsSupportsForkSuffix(t *testing.T) {
	require.Less(t, compareVersions("0.1.138", "0.1.138-fy.1"), 0)
	require.Less(t, compareVersions("0.1.138-fy.1", "0.1.138-fy.2"), 0)
	require.Less(t, compareVersions("v0.1.138-fy.2", "0.1.139"), 0)
	require.Greater(t, compareVersions("0.1.139", "0.1.138-fy.99"), 0)
	require.Equal(t, 0, compareVersions("v0.1.138-fy.2", "0.1.138-fy.2"))
}

func TestUpdateServiceDetectsForkSuffixReleaseUpdate(t *testing.T) {
	currentCommit := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	client := &updateServiceGitHubClientStub{
		releases: map[string]*GitHubRelease{
			githubRepo: {
				TagName: "v0.1.138-fy.2",
				Name:    "v0.1.138-fy.2",
				Assets: []GitHubAsset{
					{Name: "sub2api_0.1.138-fy.2_linux_amd64.tar.gz", BrowserDownloadURL: "https://github.com/1459547510/sub2api-so-fy/releases/download/v0.1.138-fy.2/sub2api_0.1.138-fy.2_linux_amd64.tar.gz"},
				},
			},
			upstreamGithubRepo: {
				TagName: "v0.1.138",
				Name:    "v0.1.138",
			},
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
					SHA: currentCommit,
				},
			},
		},
		compares: map[string]*GitHubCompare{
			githubRepo + ":" + currentCommit + ":" + "v0.1.138-fy.2": {
				Status:       "ahead",
				AheadBy:      1,
				TotalCommits: 1,
			},
		},
	}
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		client,
		"0.1.138",
		"release",
		currentCommit,
		currentCommit,
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.True(t, info.HasUpdate)
	require.True(t, info.UpdateReady)
	require.Equal(t, "0.1.138-fy.2", info.ForkLatestVersion)
	require.Equal(t, "0.1.138-fy.2", info.LatestVersion)
}

func TestUpdateServiceSuppressesForkSuffixReleaseWhenCurrentCommitMatchesTag(t *testing.T) {
	currentCommit := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	upstreamCommit := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	cache := &updateServiceCacheStub{}
	client := &updateServiceGitHubClientStub{
		releases: map[string]*GitHubRelease{
			githubRepo: {
				TagName: "v0.1.138-fy.2",
				Name:    "v0.1.138-fy.2",
			},
			upstreamGithubRepo: {
				TagName: "v0.1.138",
				Name:    "v0.1.138",
			},
		},
		branches: map[string]*GitHubBranch{
			githubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: currentCommit,
				},
			},
			upstreamGithubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: upstreamCommit,
				},
			},
		},
		compares: map[string]*GitHubCompare{
			githubRepo + ":" + currentCommit + ":" + "v0.1.138-fy.2": {
				Status: "identical",
			},
			upstreamGithubRepo + ":" + upstreamCommit + ":" + upstreamCommit: {
				Status: "identical",
			},
		},
	}
	svc := NewUpdateService(
		cache,
		client,
		"0.1.138",
		"release",
		currentCommit,
		upstreamCommit,
	)

	info, err := svc.CheckUpdate(context.Background(), true)
	require.NoError(t, err)
	require.False(t, info.HasUpdate)
	require.False(t, info.UpdateReady)
	require.Equal(t, "0.1.138-fy.2", info.ForkLatestVersion)

	cached, err := svc.CheckUpdate(context.Background(), false)
	require.NoError(t, err)
	require.True(t, cached.Cached)
	require.False(t, cached.HasUpdate)
	require.False(t, cached.UpdateReady)
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

func TestUpdateServiceDetectsUpstreamReleaseVersionButBlocksOneClickUntilForkRelease(t *testing.T) {
	forkHead := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
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
	require.False(t, info.UpstreamInfo.HasNewCommit)
	require.True(t, info.UpstreamInfo.SyncRequired)
	require.Equal(t, "release_checked", info.UpstreamInfo.Status)
	require.Empty(t, info.UpstreamInfo.CompareURL)
}

func TestUpdateServiceIgnoresUpstreamBranchCommitWhenReleaseVersionIsCurrent(t *testing.T) {
	forkHead := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	client := &updateServiceGitHubClientStub{
		releases: map[string]*GitHubRelease{
			githubRepo: {
				TagName: "v0.1.138",
				Name:    "v0.1.138",
			},
			upstreamGithubRepo: {
				TagName: "v0.1.138",
				Name:    "v0.1.138",
			},
		},
		branches: map[string]*GitHubBranch{
			githubRepo + ":" + githubBranch: {
				Name: githubBranch,
				Commit: GitHubCommitRef{
					SHA: forkHead,
				},
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
	require.False(t, info.HasUpdate)
	require.False(t, info.UpdateReady)
	require.NotNil(t, info.UpstreamInfo)
	require.False(t, info.UpstreamInfo.HasUpdate)
	require.False(t, info.UpstreamInfo.HasNewVersion)
	require.False(t, info.UpstreamInfo.HasNewCommit)
	require.False(t, info.UpstreamInfo.SyncRequired)
	require.Equal(t, "release_checked", info.UpstreamInfo.Status)
	require.NotContains(t, client.branchRequests, upstreamGithubRepo+":"+githubBranch)
}

func TestUpdateServiceAllowsForkReleaseWhenForkVersionCatchesUpWithUpstreamRelease(t *testing.T) {
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
	require.False(t, info.UpstreamInfo.HasNewCommit)
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

func newRollbackTestService(current string, releases []*GitHubRelease) *UpdateService {
	return NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{recentReleases: releases},
		current,
		"release",
	)
}

func TestUpdateServiceListRollbackVersionsFiltersAndCaps(t *testing.T) {
	releases := []*GitHubRelease{
		{TagName: "v0.1.148", PublishedAt: "2026-07-09T00:00:00Z"},                       // newer than current: excluded
		{TagName: "v0.1.147", PublishedAt: "2026-07-08T00:00:00Z"},                       // current: excluded
		{TagName: "v0.1.146-rc1", PublishedAt: "2026-07-07T12:00:00Z", Prerelease: true}, // prerelease: excluded
		{TagName: "v0.1.146", PublishedAt: "2026-07-07T00:00:00Z"},
		{TagName: "v0.1.145", PublishedAt: "2026-07-06T00:00:00Z", Draft: true}, // draft: excluded
		{TagName: "v0.1.144", PublishedAt: "2026-07-05T00:00:00Z"},
		{TagName: "v0.1.144", PublishedAt: "2026-07-05T00:00:00Z"}, // duplicate: excluded
		{TagName: "v0.1.143", PublishedAt: "2026-07-04T00:00:00Z"},
		{TagName: "v0.1.142", PublishedAt: "2026-07-03T00:00:00Z"}, // beyond cap of 3: excluded
	}
	svc := newRollbackTestService("0.1.147", releases)

	versions, err := svc.ListRollbackVersions(context.Background())

	require.NoError(t, err)
	require.Len(t, versions, 3)
	require.Equal(t, "0.1.146", versions[0].Version)
	require.Equal(t, "0.1.144", versions[1].Version)
	require.Equal(t, "0.1.143", versions[2].Version)
}

func TestUpdateServiceListRollbackVersionsSortsUnorderedInput(t *testing.T) {
	releases := []*GitHubRelease{
		{TagName: "v0.1.144"},
		{TagName: "v0.1.146"},
		{TagName: "v0.1.145"},
	}
	svc := newRollbackTestService("0.1.147", releases)

	versions, err := svc.ListRollbackVersions(context.Background())

	require.NoError(t, err)
	require.Len(t, versions, 3)
	require.Equal(t, "0.1.146", versions[0].Version)
	require.Equal(t, "0.1.145", versions[1].Version)
	require.Equal(t, "0.1.144", versions[2].Version)
}

func TestUpdateServiceListRollbackVersionsEmptyWhenNoneOlder(t *testing.T) {
	releases := []*GitHubRelease{
		{TagName: "v0.1.147"},
		{TagName: "v0.1.148"},
	}
	svc := newRollbackTestService("0.1.147", releases)

	versions, err := svc.ListRollbackVersions(context.Background())

	require.NoError(t, err)
	require.Empty(t, versions)
}

func TestUpdateServiceListRollbackVersionsPropagatesFetchError(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{recentErr: errors.New("github unavailable")},
		"0.1.147",
		"release",
	)

	_, err := svc.ListRollbackVersions(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "github unavailable")
}

func TestUpdateServiceRollbackToVersionRejectsDisallowedTargets(t *testing.T) {
	releases := []*GitHubRelease{
		{TagName: "v0.1.148"},
		{TagName: "v0.1.147"},
		{TagName: "v0.1.146"},
		{TagName: "v0.1.145"},
		{TagName: "v0.1.144"},
		{TagName: "v0.1.143"},
		{TagName: "v0.1.142"},
	}
	svc := newRollbackTestService("0.1.147", releases)

	for _, target := range []string{
		"",         // empty
		"0.1.147",  // current version
		"v0.1.147", // current version with prefix
		"0.1.148",  // newer than current
		"0.1.142",  // older than the 3 most recent
		"9.9.9",    // nonexistent
	} {
		err := svc.RollbackToVersion(context.Background(), target)
		require.ErrorIs(t, err, ErrRollbackVersionNotAllowed, "target %q should be rejected", target)
	}
}

func TestUpdateServiceRollbackToVersionAcceptsVPrefix(t *testing.T) {
	// No platform asset in the release: the target passes the allowlist check
	// and fails later at asset lookup, proving the version itself was accepted.
	releases := []*GitHubRelease{
		{TagName: "v0.1.147"},
		{TagName: "v0.1.146"},
	}
	svc := newRollbackTestService("0.1.147", releases)

	err := svc.RollbackToVersion(context.Background(), "v0.1.146")

	require.Error(t, err)
	require.NotErrorIs(t, err, ErrRollbackVersionNotAllowed)
	require.Contains(t, err.Error(), "no compatible release found")
}
