package service

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var (
	ErrNoUpdateAvailable = infraerrors.Conflict("ALREADY_UP_TO_DATE", "no update available; current version is latest")
)

const (
	updateCacheKey = "update_check_cache"
	updateCacheTTL = 1200 // 20 minutes

	// Update releases must be resolved from this fork instead of the original
	// upstream repository. The original repository is checked only for released
	// upstream versions, not for raw branch commits, so every upstream commit
	// does not become an application update prompt.
	githubRepo         = "1459547510/sub2api-so-fy"
	upstreamGithubRepo = "Wei-Shaw/sub2api"
	githubBranch       = "main"

	// Security: allowed download domains for updates
	allowedDownloadHost = "github.com"
	allowedAssetHost    = "objects.githubusercontent.com"

	// Security: max download size (500MB)
	maxDownloadSize = 500 * 1024 * 1024
)

// UpdateCache defines cache operations for update service
type UpdateCache interface {
	GetUpdateInfo(ctx context.Context) (string, error)
	SetUpdateInfo(ctx context.Context, data string, ttl time.Duration) error
}

// GitHubReleaseClient 获取 GitHub release 信息的接口
type GitHubReleaseClient interface {
	FetchLatestRelease(ctx context.Context, repo string) (*GitHubRelease, error)
	FetchBranch(ctx context.Context, repo, branch string) (*GitHubBranch, error)
	CompareCommits(ctx context.Context, repo, base, head string) (*GitHubCompare, error)
	DownloadFile(ctx context.Context, url, dest string, maxSize int64) error
	FetchChecksumFile(ctx context.Context, url string) ([]byte, error)
}

// UpdateService handles software updates
type UpdateService struct {
	cache          UpdateCache
	githubClient   GitHubReleaseClient
	currentVersion string
	currentCommit  string
	upstreamCommit string
	buildType      string // "source" for manual builds, "release" for CI builds
}

// NewUpdateService creates a new UpdateService
func NewUpdateService(cache UpdateCache, githubClient GitHubReleaseClient, version, buildType string, commit ...string) *UpdateService {
	currentCommit := "unknown"
	if len(commit) > 0 && strings.TrimSpace(commit[0]) != "" {
		currentCommit = strings.TrimSpace(commit[0])
	}
	if normalizeCommitSHA(currentCommit) == "" {
		currentCommit = resolveBuildCommit()
	}
	upstreamCommit := "unknown"
	if len(commit) > 1 && strings.TrimSpace(commit[1]) != "" {
		upstreamCommit = strings.TrimSpace(commit[1])
	}
	return &UpdateService{
		cache:          cache,
		githubClient:   githubClient,
		currentVersion: version,
		currentCommit:  currentCommit,
		upstreamCommit: upstreamCommit,
		buildType:      buildType,
	}
}

// UpdateInfo contains update information
type UpdateInfo struct {
	CurrentVersion    string        `json:"current_version"`
	LatestVersion     string        `json:"latest_version"`
	ForkLatestVersion string        `json:"fork_latest_version"`
	HasUpdate         bool          `json:"has_update"`
	UpdateReady       bool          `json:"update_ready"`
	ReleaseInfo       *ReleaseInfo  `json:"release_info,omitempty"`
	BranchInfo        *BranchInfo   `json:"branch_info,omitempty"`
	UpstreamInfo      *UpstreamInfo `json:"upstream_info,omitempty"`
	Cached            bool          `json:"cached"`
	Warning           string        `json:"warning,omitempty"`
	BuildType         string        `json:"build_type"` // "source" or "release"
}

// ReleaseInfo contains GitHub release details
type ReleaseInfo struct {
	Name        string  `json:"name"`
	Body        string  `json:"body"`
	PublishedAt string  `json:"published_at"`
	HTMLURL     string  `json:"html_url"`
	Assets      []Asset `json:"assets,omitempty"`
}

// Asset represents a release asset
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	Size        int64  `json:"size"`
}

// BranchInfo contains default branch head information for source/fork updates.
type BranchInfo struct {
	Repo          string `json:"repo"`
	Branch        string `json:"branch"`
	CurrentCommit string `json:"current_commit"`
	LatestCommit  string `json:"latest_commit"`
	HasNewCommit  bool   `json:"has_new_commit"`
	CanCompare    bool   `json:"can_compare"`
	Status        string `json:"status"`
	CompareURL    string `json:"compare_url,omitempty"`
	CommitURL     string `json:"commit_url,omitempty"`
}

// UpstreamInfo contains original repository update information.
type UpstreamInfo struct {
	Repo          string       `json:"repo"`
	Branch        string       `json:"branch"`
	LatestVersion string       `json:"latest_version"`
	HasUpdate     bool         `json:"has_update"`
	HasNewVersion bool         `json:"has_new_version"`
	HasNewCommit  bool         `json:"has_new_commit"`
	SyncRequired  bool         `json:"sync_required"`
	CanCompare    bool         `json:"can_compare"`
	Status        string       `json:"status"`
	ReleaseInfo   *ReleaseInfo `json:"release_info,omitempty"`
	CurrentCommit string       `json:"current_commit"`
	LatestCommit  string       `json:"latest_commit"`
	CompareURL    string       `json:"compare_url,omitempty"`
	CommitURL     string       `json:"commit_url,omitempty"`
	Warning       string       `json:"warning,omitempty"`
}

// GitHubRelease represents GitHub API response
type GitHubRelease struct {
	TagName     string        `json:"tag_name"`
	Name        string        `json:"name"`
	Body        string        `json:"body"`
	PublishedAt string        `json:"published_at"`
	HTMLURL     string        `json:"html_url"`
	Assets      []GitHubAsset `json:"assets"`
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// GitHubBranch represents GitHub branch API response.
type GitHubBranch struct {
	Name   string          `json:"name"`
	Commit GitHubCommitRef `json:"commit"`
}

// GitHubCommitRef represents the commit object returned by GitHub branch API.
type GitHubCommitRef struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

// GitHubCompare represents GitHub compare API response.
type GitHubCompare struct {
	Status       string `json:"status"`
	AheadBy      int    `json:"ahead_by"`
	BehindBy     int    `json:"behind_by"`
	HTMLURL      string `json:"html_url"`
	Permalink    string `json:"permalink_url"`
	TotalCommits int    `json:"total_commits"`
}

// CheckUpdate checks for available updates
func (s *UpdateService) CheckUpdate(ctx context.Context, force bool) (*UpdateInfo, error) {
	// Try cache first
	if !force {
		if cached, err := s.getFromCache(ctx); err == nil && cached != nil {
			return cached, nil
		}
	}

	// Fetch from GitHub
	info, err := s.fetchLatestRelease(ctx)
	if err != nil {
		branchInfo, branchErr := s.fetchBranchInfo(ctx)
		upstreamInfo, upstreamErr := s.fetchUpstreamInfo(ctx, nil)
		// Return cached on error
		if cached, cacheErr := s.getFromCache(ctx); cacheErr == nil && cached != nil {
			cached.Warning = "Using cached data: " + appendWarnings(err, branchErr, upstreamErr)
			if cached.BranchInfo == nil {
				cached.BranchInfo = branchInfo
			}
			if cached.UpstreamInfo == nil {
				cached.UpstreamInfo = upstreamInfo
			}
			return cached, nil
		}
		return &UpdateInfo{
			CurrentVersion:    s.currentVersion,
			LatestVersion:     upstreamLatestVersionOrCurrent(s.currentVersion, upstreamInfo),
			ForkLatestVersion: s.currentVersion,
			HasUpdate:         upstreamInfo != nil && upstreamInfo.HasUpdate,
			UpdateReady:       false,
			BranchInfo:        branchInfo,
			UpstreamInfo:      upstreamInfo,
			Warning:           appendWarnings(err, branchErr, upstreamErr),
			BuildType:         s.buildType,
		}, nil
	}

	// Cache result
	s.saveToCache(ctx, info)
	return info, nil
}

// PerformUpdate downloads and applies the update
// Uses atomic file replacement pattern for safe in-place updates
func (s *UpdateService) PerformUpdate(ctx context.Context) error {
	info, err := s.CheckUpdate(ctx, true)
	if err != nil {
		return err
	}

	if !info.UpdateReady {
		return ErrNoUpdateAvailable
	}

	// Find matching archive and checksum for current platform
	archiveName := s.getArchiveName()
	var downloadURL string
	var checksumURL string

	for _, asset := range info.ReleaseInfo.Assets {
		if strings.Contains(asset.Name, archiveName) && !strings.HasSuffix(asset.Name, ".txt") {
			downloadURL = asset.DownloadURL
		}
		if asset.Name == "checksums.txt" {
			checksumURL = asset.DownloadURL
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no compatible release found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// SECURITY: Validate download URL is from trusted domain
	if err := validateDownloadURL(downloadURL); err != nil {
		return fmt.Errorf("invalid download URL: %w", err)
	}
	if checksumURL != "" {
		if err := validateDownloadURL(checksumURL); err != nil {
			return fmt.Errorf("invalid checksum URL: %w", err)
		}
	}

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	exeDir := filepath.Dir(exePath)

	// Create temp directory in the SAME directory as executable
	// This ensures os.Rename is atomic (same filesystem)
	tempDir, err := os.MkdirTemp(exeDir, ".sub2api-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Download archive
	archivePath := filepath.Join(tempDir, filepath.Base(downloadURL))
	if err := s.downloadFile(ctx, downloadURL, archivePath); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Verify checksum if available
	if checksumURL != "" {
		if err := s.verifyChecksum(ctx, archivePath, checksumURL); err != nil {
			return fmt.Errorf("checksum verification failed: %w", err)
		}
	}

	// Extract binary from archive
	newBinaryPath := filepath.Join(tempDir, "sub2api")
	if err := s.extractBinary(archivePath, newBinaryPath); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	// Set executable permission before replacement
	if err := os.Chmod(newBinaryPath, 0755); err != nil {
		return fmt.Errorf("chmod failed: %w", err)
	}

	// Atomic replacement using rename pattern:
	// 1. Rename current -> backup (atomic on Unix)
	// 2. Rename new -> current (atomic on Unix, same filesystem)
	// If step 2 fails, restore backup
	backupPath := exePath + ".backup"

	// Remove old backup if exists
	_ = os.Remove(backupPath)

	// Step 1: Move current binary to backup
	if err := os.Rename(exePath, backupPath); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	// Step 2: Move new binary to target location (atomic, same filesystem)
	if err := os.Rename(newBinaryPath, exePath); err != nil {
		// Restore backup on failure
		if restoreErr := os.Rename(backupPath, exePath); restoreErr != nil {
			return fmt.Errorf("replace failed and restore failed: %w (restore error: %v)", err, restoreErr)
		}
		return fmt.Errorf("replace failed (restored backup): %w", err)
	}

	// Success - backup file is kept for rollback capability
	// It will be cleaned up on next successful update
	return nil
}

// Rollback restores the previous version
func (s *UpdateService) Rollback() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("failed to resolve symlinks: %w", err)
	}

	backupFile := exePath + ".backup"
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("no backup found")
	}

	// Replace current with backup
	if err := os.Rename(backupFile, exePath); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	return nil
}

func (s *UpdateService) fetchLatestRelease(ctx context.Context) (*UpdateInfo, error) {
	release, err := s.githubClient.FetchLatestRelease(ctx, githubRepo)
	if err != nil {
		return nil, err
	}

	forkLatestVersion := strings.TrimPrefix(release.TagName, "v")
	forkHasUpdate, forkCompareErr := s.forkReleaseHasUpdate(ctx, release, forkLatestVersion)

	assets := make([]Asset, len(release.Assets))
	for i, a := range release.Assets {
		assets[i] = Asset{
			Name:        a.Name,
			DownloadURL: a.BrowserDownloadURL,
			Size:        a.Size,
		}
	}

	info := &UpdateInfo{
		CurrentVersion:    s.currentVersion,
		LatestVersion:     forkLatestVersion,
		ForkLatestVersion: forkLatestVersion,
		HasUpdate:         forkHasUpdate,
		UpdateReady:       forkHasUpdate,
		ReleaseInfo: &ReleaseInfo{
			Name:        release.Name,
			Body:        release.Body,
			PublishedAt: release.PublishedAt,
			HTMLURL:     release.HTMLURL,
			Assets:      assets,
		},
		Cached:    false,
		BuildType: s.buildType,
	}

	branchInfo, branchErr := s.fetchBranchInfo(ctx)
	info.BranchInfo = branchInfo
	upstreamInfo, upstreamErr := s.fetchUpstreamInfo(ctx, release)
	info.UpstreamInfo = upstreamInfo
	if upstreamInfo != nil {
		info.LatestVersion = displayLatestVersion(s.currentVersion, forkLatestVersion, upstreamInfo)
		if upstreamInfo.HasUpdate {
			info.HasUpdate = true
		}
		if upstreamInfo.HasUpdate && upstreamInfo.SyncRequired {
			// 原仓库已经有更新但当前 fork 尚未同步/发布时，不允许一键更新，
			// 避免用未合并上游变更的 fork Release 覆盖运行实例。
			info.UpdateReady = false
		}
	}
	if branchErr != nil || upstreamErr != nil || forkCompareErr != nil {
		info.Warning = appendWarnings(nil, forkCompareErr, branchErr, upstreamErr)
	}

	return info, nil
}

func (s *UpdateService) forkReleaseHasUpdate(ctx context.Context, release *GitHubRelease, forkLatestVersion string) (bool, error) {
	if release == nil || compareVersions(s.currentVersion, forkLatestVersion) >= 0 {
		return false, nil
	}

	currentCommit := normalizeCommitSHA(s.currentCommit)
	tagRef := strings.TrimSpace(release.TagName)
	if currentCommit == "" || tagRef == "" {
		return true, nil
	}

	compare, err := s.githubClient.CompareCommits(ctx, githubRepo, currentCommit, tagRef)
	if err != nil {
		// Do not hide a version-tag update just because the optional commit
		// comparison failed (for example when GitHub rate-limits comparisons).
		return true, fmt.Errorf("fork release compare: %w", err)
	}
	switch strings.ToLower(strings.TrimSpace(compare.Status)) {
	case "identical", "behind":
		// base=current, head=release tag. "behind" means the release tag is
		// behind the running commit, so applying it would not be an upgrade.
		return false, nil
	case "ahead", "diverged":
		return true, nil
	default:
		if compare.AheadBy > 0 || compare.TotalCommits > 0 {
			return true, nil
		}
		return true, nil
	}
}

func (s *UpdateService) fetchBranchInfo(ctx context.Context) (*BranchInfo, error) {
	branch, err := s.githubClient.FetchBranch(ctx, githubRepo, githubBranch)
	if err != nil {
		return nil, err
	}

	latestCommit := normalizeCommitSHA(branch.Commit.SHA)
	if latestCommit == "" {
		return nil, fmt.Errorf("empty latest commit for %s/%s", githubRepo, githubBranch)
	}

	currentCommit := normalizeCommitSHA(s.currentCommit)
	info := s.branchInfoForLatestCommit(latestCommit, currentCommit)

	return info, nil
}

func (s *UpdateService) branchInfoForLatestCommit(latestCommit, currentCommit string) *BranchInfo {
	latestCommit = normalizeCommitSHA(latestCommit)
	currentCommit = normalizeCommitSHA(currentCommit)
	info := &BranchInfo{
		Repo:          githubRepo,
		Branch:        githubBranch,
		CurrentCommit: unknownIfEmpty(currentCommit),
		LatestCommit:  latestCommit,
		CanCompare:    currentCommit != "" && latestCommit != "",
		Status:        "unknown_current",
	}
	if latestCommit != "" {
		info.CommitURL = fmt.Sprintf("https://github.com/%s/commit/%s", githubRepo, latestCommit)
	}
	if currentCommit == "" {
		info.CompareURL = fmt.Sprintf("https://github.com/%s/commits/%s", githubRepo, githubBranch)
		return info
	}

	if commitMatches(currentCommit, latestCommit) {
		info.Status = "current"
		return info
	}

	info.HasNewCommit = true
	info.Status = "behind"
	info.CompareURL = fmt.Sprintf("https://github.com/%s/compare/%s...%s", githubRepo, currentCommit, latestCommit)
	return info
}

func (s *UpdateService) fetchUpstreamInfo(ctx context.Context, forkRelease *GitHubRelease) (*UpstreamInfo, error) {
	info := &UpstreamInfo{
		Repo:       upstreamGithubRepo,
		Branch:     githubBranch,
		Status:     "release_checked",
		CanCompare: false,
	}

	var errs []error
	forkLatestVersion := ""
	if forkRelease != nil {
		forkLatestVersion = strings.TrimPrefix(forkRelease.TagName, "v")
	}

	if upstreamRelease, err := s.githubClient.FetchLatestRelease(ctx, upstreamGithubRepo); err != nil {
		errs = append(errs, fmt.Errorf("upstream release: %w", err))
		info.Status = "release_unavailable"
	} else {
		latestVersion := strings.TrimPrefix(upstreamRelease.TagName, "v")
		info.LatestVersion = latestVersion
		info.ReleaseInfo = releaseInfoFromGitHub(upstreamRelease)
	}

	normalizeUpstreamInfoForReleaseOnly(info, s.currentVersion, forkLatestVersion)

	warnErr := joinUpdateWarnings(errs...)
	if warnErr != nil {
		info.Warning = warnErr.Error()
	}
	return info, warnErr
}

func (s *UpdateService) downloadFile(ctx context.Context, downloadURL, dest string) error {
	return s.githubClient.DownloadFile(ctx, downloadURL, dest, maxDownloadSize)
}

func (s *UpdateService) getArchiveName() string {
	osName := runtime.GOOS
	arch := runtime.GOARCH
	return fmt.Sprintf("%s_%s", osName, arch)
}

// validateDownloadURL checks if the URL is from an allowed domain
// SECURITY: This prevents SSRF and ensures downloads only come from trusted GitHub domains
func validateDownloadURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Must be HTTPS
	if parsedURL.Scheme != "https" {
		return fmt.Errorf("only HTTPS URLs are allowed")
	}

	// Check against allowed hosts
	host := parsedURL.Host
	// GitHub release URLs can be from github.com or objects.githubusercontent.com
	if host != allowedDownloadHost &&
		!strings.HasSuffix(host, "."+allowedDownloadHost) &&
		host != allowedAssetHost &&
		!strings.HasSuffix(host, "."+allowedAssetHost) {
		return fmt.Errorf("download from untrusted host: %s", host)
	}

	return nil
}

func (s *UpdateService) verifyChecksum(ctx context.Context, filePath, checksumURL string) error {
	// Download checksums file
	checksumData, err := s.githubClient.FetchChecksumFile(ctx, checksumURL)
	if err != nil {
		return fmt.Errorf("failed to download checksums: %w", err)
	}

	// Calculate file hash
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	actualHash := hex.EncodeToString(h.Sum(nil))

	// Find expected hash in checksums file
	fileName := filepath.Base(filePath)
	scanner := bufio.NewScanner(strings.NewReader(string(checksumData)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[1] == fileName {
			if parts[0] == actualHash {
				return nil
			}
			return fmt.Errorf("checksum mismatch: expected %s, got %s", parts[0], actualHash)
		}
	}

	return fmt.Errorf("checksum not found for %s", fileName)
}

func (s *UpdateService) extractBinary(archivePath, destPath string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	var reader io.Reader = f

	// Handle gzip compression
	if strings.HasSuffix(archivePath, ".gz") || strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz") {
		gzr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer func() { _ = gzr.Close() }()
		reader = gzr
	}

	// Handle tar archive
	if strings.Contains(archivePath, ".tar") {
		tr := tar.NewReader(reader)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			// SECURITY: Prevent Zip Slip / Path Traversal attack
			// Only allow files with safe base names, no directory traversal
			baseName := filepath.Base(hdr.Name)

			// Check for path traversal attempts
			if strings.Contains(hdr.Name, "..") {
				return fmt.Errorf("path traversal attempt detected: %s", hdr.Name)
			}

			// Validate the entry is a regular file
			if hdr.Typeflag != tar.TypeReg {
				continue // Skip directories and special files
			}

			// Only extract the specific binary we need
			if baseName == "sub2api" || baseName == "sub2api.exe" {
				// Additional security: limit file size (max 500MB)
				const maxBinarySize = 500 * 1024 * 1024
				if hdr.Size > maxBinarySize {
					return fmt.Errorf("binary too large: %d bytes (max %d)", hdr.Size, maxBinarySize)
				}

				out, err := os.Create(destPath)
				if err != nil {
					return err
				}

				// Use LimitReader to prevent decompression bombs
				limited := io.LimitReader(tr, maxBinarySize)
				if _, err := io.Copy(out, limited); err != nil {
					_ = out.Close()
					return err
				}
				if err := out.Close(); err != nil {
					return err
				}
				return nil
			}
		}
		return fmt.Errorf("binary not found in archive")
	}

	// Direct copy for non-tar files (with size limit)
	const maxBinarySize = 500 * 1024 * 1024
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}

	limited := io.LimitReader(reader, maxBinarySize)
	if _, err := io.Copy(out, limited); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func (s *UpdateService) getFromCache(ctx context.Context) (*UpdateInfo, error) {
	data, err := s.cache.GetUpdateInfo(ctx)
	if err != nil {
		return nil, err
	}

	var cached struct {
		CurrentVersion string        `json:"current_version"`
		CurrentCommit  string        `json:"current_commit"`
		Latest         string        `json:"latest"`
		ForkLatest     string        `json:"fork_latest"`
		HasUpdate      *bool         `json:"has_update"`
		UpdateReady    *bool         `json:"update_ready"`
		ReleaseInfo    *ReleaseInfo  `json:"release_info"`
		BranchInfo     *BranchInfo   `json:"branch_info"`
		UpstreamInfo   *UpstreamInfo `json:"upstream_info"`
		Warning        string        `json:"warning"`
		Timestamp      int64         `json:"timestamp"`
	}
	if err := json.Unmarshal([]byte(data), &cached); err != nil {
		return nil, err
	}

	if time.Now().Unix()-cached.Timestamp > updateCacheTTL {
		return nil, fmt.Errorf("cache expired")
	}
	if cached.CurrentVersion != "" && cached.CurrentVersion != s.currentVersion {
		return nil, fmt.Errorf("cache version mismatch")
	}
	if cached.CurrentCommit != "" && !commitMatches(normalizeCommitSHA(cached.CurrentCommit), normalizeCommitSHA(s.currentCommit)) {
		return nil, fmt.Errorf("cache commit mismatch")
	}

	forkLatest := cached.ForkLatest
	if forkLatest == "" {
		forkLatest = cached.Latest
	}
	normalizeUpstreamInfoForReleaseOnly(cached.UpstreamInfo, s.currentVersion, forkLatest)
	latest := displayLatestVersion(s.currentVersion, forkLatest, cached.UpstreamInfo)
	forkHasUpdate := compareVersions(s.currentVersion, forkLatest) < 0
	upstreamHasUpdate := cached.UpstreamInfo != nil && cached.UpstreamInfo.HasUpdate
	upstreamSyncRequired := cached.UpstreamInfo != nil && cached.UpstreamInfo.HasUpdate && cached.UpstreamInfo.SyncRequired
	hasUpdate := forkHasUpdate || upstreamHasUpdate
	updateReady := forkHasUpdate && !upstreamSyncRequired
	if cached.HasUpdate != nil {
		hasUpdate = *cached.HasUpdate
	}
	if cached.UpdateReady != nil {
		updateReady = *cached.UpdateReady
	}
	return &UpdateInfo{
		CurrentVersion:    s.currentVersion,
		LatestVersion:     latest,
		ForkLatestVersion: forkLatest,
		HasUpdate:         hasUpdate,
		UpdateReady:       updateReady,
		ReleaseInfo:       cached.ReleaseInfo,
		BranchInfo:        cached.BranchInfo,
		UpstreamInfo:      cached.UpstreamInfo,
		Cached:            true,
		Warning:           cached.Warning,
		BuildType:         s.buildType,
	}, nil
}

func (s *UpdateService) saveToCache(ctx context.Context, info *UpdateInfo) {
	cacheData := struct {
		CurrentVersion string        `json:"current_version"`
		CurrentCommit  string        `json:"current_commit"`
		Latest         string        `json:"latest"`
		ForkLatest     string        `json:"fork_latest"`
		HasUpdate      bool          `json:"has_update"`
		UpdateReady    bool          `json:"update_ready"`
		ReleaseInfo    *ReleaseInfo  `json:"release_info"`
		BranchInfo     *BranchInfo   `json:"branch_info"`
		UpstreamInfo   *UpstreamInfo `json:"upstream_info"`
		Warning        string        `json:"warning"`
		Timestamp      int64         `json:"timestamp"`
	}{
		CurrentVersion: s.currentVersion,
		CurrentCommit:  normalizeCommitSHA(s.currentCommit),
		Latest:         info.LatestVersion,
		ForkLatest:     info.ForkLatestVersion,
		HasUpdate:      info.HasUpdate,
		UpdateReady:    info.UpdateReady,
		ReleaseInfo:    info.ReleaseInfo,
		BranchInfo:     info.BranchInfo,
		UpstreamInfo:   info.UpstreamInfo,
		Warning:        info.Warning,
		Timestamp:      time.Now().Unix(),
	}

	data, _ := json.Marshal(cacheData)
	_ = s.cache.SetUpdateInfo(ctx, string(data), time.Duration(updateCacheTTL)*time.Second)
}

// compareVersions compares upstream semantic versions plus fork release suffixes.
//
// Fork releases use tags like v0.1.138-fy.2 while the embedded upstream base
// version may still be 0.1.138. For the fork update channel, suffix builds are
// intentionally treated as newer than their base upstream version:
// 0.1.138 < 0.1.138-fy.1 < 0.1.138-fy.2 < 0.1.139.
func compareVersions(current, latest string) int {
	currentParts := parseVersion(current)
	latestParts := parseVersion(latest)

	for i := 0; i < 3; i++ {
		if currentParts.core[i] < latestParts.core[i] {
			return -1
		}
		if currentParts.core[i] > latestParts.core[i] {
			return 1
		}
	}

	if currentParts.hasSuffix != latestParts.hasSuffix {
		if currentParts.hasSuffix {
			return 1
		}
		return -1
	}
	if currentParts.hasSuffix {
		if cmp := strings.Compare(currentParts.suffixLabel, latestParts.suffixLabel); cmp != 0 {
			return cmp
		}
		if currentParts.suffixNumber != latestParts.suffixNumber {
			if currentParts.suffixNumber < latestParts.suffixNumber {
				return -1
			}
			return 1
		}
		if cmp := strings.Compare(currentParts.suffix, latestParts.suffix); cmp != 0 {
			return cmp
		}
	}
	return 0
}

type parsedVersion struct {
	core         [3]int
	suffix       string
	suffixLabel  string
	suffixNumber int
	hasSuffix    bool
}

func parseVersion(v string) parsedVersion {
	v = strings.TrimSpace(strings.TrimPrefix(v, "v"))
	if plus := strings.Index(v, "+"); plus >= 0 {
		v = v[:plus]
	}

	result := parsedVersion{core: [3]int{0, 0, 0}, suffixNumber: -1}
	if dash := strings.Index(v, "-"); dash >= 0 {
		result.suffix = strings.ToLower(strings.TrimSpace(v[dash+1:]))
		result.hasSuffix = result.suffix != ""
		v = v[:dash]
	}

	parts := strings.Split(v, ".")
	for i := 0; i < len(parts) && i < 3; i++ {
		part := leadingDigits(parts[i])
		if parsed, err := strconv.Atoi(part); err == nil {
			result.core[i] = parsed
		}
	}

	if result.hasSuffix {
		suffixParts := strings.FieldsFunc(result.suffix, func(r rune) bool {
			return r == '.' || r == '-' || r == '_'
		})
		if len(suffixParts) > 0 {
			result.suffixLabel = suffixParts[0]
		}
		for i := len(suffixParts) - 1; i >= 0; i-- {
			if parsed, err := strconv.Atoi(suffixParts[i]); err == nil {
				result.suffixNumber = parsed
				break
			}
		}
	}
	return result
}

func leadingDigits(value string) string {
	value = strings.TrimSpace(value)
	for i, r := range value {
		if r < '0' || r > '9' {
			return value[:i]
		}
	}
	return value
}

func releaseInfoFromGitHub(release *GitHubRelease) *ReleaseInfo {
	if release == nil {
		return nil
	}
	assets := make([]Asset, len(release.Assets))
	for i, a := range release.Assets {
		assets[i] = Asset{
			Name:        a.Name,
			DownloadURL: a.BrowserDownloadURL,
			Size:        a.Size,
		}
	}
	return &ReleaseInfo{
		Name:        release.Name,
		Body:        release.Body,
		PublishedAt: release.PublishedAt,
		HTMLURL:     release.HTMLURL,
		Assets:      assets,
	}
}

func normalizeUpstreamInfoForReleaseOnly(info *UpstreamInfo, currentVersion, forkLatestVersion string) {
	if info == nil {
		return
	}
	if info.LatestVersion == "" {
		info.LatestVersion = currentVersion
	}
	info.HasNewVersion = compareVersions(currentVersion, info.LatestVersion) < 0
	info.HasUpdate = info.HasNewVersion
	info.HasNewCommit = false
	info.CurrentCommit = ""
	info.LatestCommit = ""
	info.CompareURL = ""
	info.CommitURL = ""
	info.CanCompare = false
	if info.Status == "" || info.Status == "unknown" || info.Status == "unknown_current" || info.Status == "unknown_compare" || info.Status == "ahead" || info.Status == "diverged" || info.Status == "behind" || info.Status == "current" {
		info.Status = "release_checked"
	}
	if !info.HasNewVersion {
		info.SyncRequired = false
		return
	}
	if strings.TrimSpace(forkLatestVersion) == "" {
		info.SyncRequired = true
		return
	}
	info.SyncRequired = compareVersions(forkLatestVersion, info.LatestVersion) < 0
}

func normalizeCommitSHA(sha string) string {
	sha = strings.ToLower(strings.TrimSpace(sha))
	if sha == "" || sha == "unknown" {
		return ""
	}
	for _, r := range sha {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return ""
		}
	}
	return sha
}

func unknownIfEmpty(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

func commitMatches(current, latest string) bool {
	if current == "" || latest == "" {
		return false
	}
	if len(current) > len(latest) {
		return strings.HasPrefix(current, latest)
	}
	return strings.HasPrefix(latest, current)
}

func appendWarnings(errs ...error) string {
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}
		message := strings.TrimSpace(err.Error())
		if message != "" {
			messages = append(messages, message)
		}
	}
	return strings.Join(messages, "; ")
}

func joinUpdateWarnings(errs ...error) error {
	message := appendWarnings(errs...)
	if message == "" {
		return nil
	}
	return fmt.Errorf("%s", message)
}

func displayLatestVersion(current, forkLatest string, upstreamInfo *UpstreamInfo) string {
	latest := strings.TrimSpace(forkLatest)
	if latest == "" {
		latest = current
	}
	if upstreamInfo != nil && strings.TrimSpace(upstreamInfo.LatestVersion) != "" && compareVersions(latest, upstreamInfo.LatestVersion) < 0 {
		latest = upstreamInfo.LatestVersion
	}
	return latest
}

func upstreamLatestVersionOrCurrent(current string, upstreamInfo *UpstreamInfo) string {
	if upstreamInfo != nil && upstreamInfo.LatestVersion != "" {
		return upstreamInfo.LatestVersion
	}
	return current
}

func resolveBuildCommit() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" && normalizeCommitSHA(setting.Value) != "" {
			return setting.Value
		}
	}
	return "unknown"
}
