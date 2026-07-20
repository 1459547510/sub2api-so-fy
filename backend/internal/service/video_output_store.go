package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const VideoOutputMaxBytes int64 = 512 * 1024 * 1024

var (
	ErrVideoOutputURLMissing = errors.New("completed video result has no video URL")
	ErrVideoOutputInvalid    = errors.New("video output is not a valid MP4")
	ErrVideoOutputNotFound   = errors.New("video output not found")
)

type VideoOutputStore struct {
	root     string
	client   *http.Client
	maxBytes int64
}

func NewVideoOutputStore(dataDir string) *VideoOutputStore {
	return &VideoOutputStore{
		root:     filepath.Join(dataDir, "video-outputs"),
		client:   &http.Client{Timeout: 10 * time.Minute},
		maxBytes: VideoOutputMaxBytes,
	}
}

func (s *VideoOutputStore) Save(ctx context.Context, jobID string, result json.RawMessage) (json.RawMessage, error) {
	if s == nil || !validVideoOutputJobID(jobID) {
		return nil, ErrVideoOutputNotFound
	}
	payload, item, sourceURL, err := parseVideoOutputResult(result)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(s.root, 0o700); err != nil {
		return nil, err
	}
	path := s.path(jobID)
	if err := validateVideoOutputFile(path, s.limit()); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			_ = os.Remove(path)
		}
		if err := s.download(ctx, sourceURL, path, jobID); err != nil {
			return nil, err
		}
	}

	localURL := VideoOutputURL(jobID)
	item["source_url"] = sourceURL
	item["mp4_url"] = localURL
	item["url"] = localURL
	item["local_url"] = localURL
	removeVideoOutputSaveAttempts(payload)
	rewritten, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("encode saved video result: %w", err)
	}
	return rewritten, nil
}

func (s *VideoOutputStore) Open(jobID string) (*os.File, error) {
	if s == nil || !validVideoOutputJobID(jobID) {
		return nil, ErrVideoOutputNotFound
	}
	path := s.path(jobID)
	if err := validateVideoOutputFile(path, s.limit()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrVideoOutputNotFound
		}
		return nil, err
	}
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrVideoOutputNotFound
	}
	return file, err
}

func (s *VideoOutputStore) download(ctx context.Context, sourceURL, destination, jobID string) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return fmt.Errorf("create video output request: %w", err)
	}
	client := s.client
	if client == nil {
		client = http.DefaultClient
	}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("download video output: %w", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("download video output: upstream returned HTTP %d", response.StatusCode)
	}
	limit := s.limit()
	if response.ContentLength > limit {
		return fmt.Errorf("download video output: file exceeds %d bytes", limit)
	}
	temporary, err := os.CreateTemp(s.root, jobID+"-*.tmp")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() { _ = os.Remove(temporaryPath) }()
	written, copyErr := io.Copy(temporary, io.LimitReader(response.Body, limit+1))
	closeErr := temporary.Close()
	if copyErr != nil {
		return fmt.Errorf("save video output: %w", copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close video output: %w", closeErr)
	}
	if written > limit {
		return fmt.Errorf("download video output: file exceeds %d bytes", limit)
	}
	if err := validateVideoOutputFile(temporaryPath, limit); err != nil {
		return err
	}
	if err := os.Rename(temporaryPath, destination); err != nil {
		return fmt.Errorf("publish video output: %w", err)
	}
	return nil
}

func (s *VideoOutputStore) path(jobID string) string {
	return filepath.Join(s.root, jobID+".mp4")
}

func (s *VideoOutputStore) limit() int64 {
	if s != nil && s.maxBytes > 0 {
		return s.maxBytes
	}
	return VideoOutputMaxBytes
}

func VideoOutputURL(jobID string) string {
	return "/v1/videos/jobs/" + url.PathEscape(strings.TrimSpace(jobID)) + "/content"
}

func parseVideoOutputResult(result json.RawMessage) (map[string]any, map[string]any, string, error) {
	var payload map[string]any
	if len(result) == 0 || json.Unmarshal(result, &payload) != nil {
		return nil, nil, "", ErrVideoOutputURLMissing
	}
	data, ok := payload["data"].([]any)
	if !ok || len(data) == 0 {
		return nil, nil, "", ErrVideoOutputURLMissing
	}
	item, ok := data[0].(map[string]any)
	if !ok {
		return nil, nil, "", ErrVideoOutputURLMissing
	}
	for _, key := range []string{"source_url", "mp4_url", "video_url", "url"} {
		raw, _ := item[key].(string)
		raw = strings.TrimSpace(raw)
		parsed, err := url.Parse(raw)
		if err == nil && parsed.Host != "" && (parsed.Scheme == "http" || parsed.Scheme == "https") {
			return payload, item, raw, nil
		}
	}
	return nil, nil, "", ErrVideoOutputURLMissing
}

func validateVideoOutputFile(path string, maxBytes int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() < 8 || info.Size() > maxBytes {
		return ErrVideoOutputInvalid
	}
	header := make([]byte, 8)
	if _, err := io.ReadFull(file, header); err != nil {
		return ErrVideoOutputInvalid
	}
	if string(header[4:8]) != "ftyp" {
		return ErrVideoOutputInvalid
	}
	return nil
}

func validVideoOutputJobID(jobID string) bool {
	trimmed := strings.TrimSpace(jobID)
	return jobID == trimmed && jobID != "" && jobID != "." && filepath.Base(jobID) == jobID && !strings.ContainsAny(jobID, `/\\`)
}

func removeVideoOutputSaveAttempts(payload map[string]any) {
	provider, ok := payload["provider"].(map[string]any)
	if !ok {
		return
	}
	delete(provider, "local_save_attempts")
}
