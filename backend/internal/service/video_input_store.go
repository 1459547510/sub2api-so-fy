package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"net/http"
)

const VideoInputMaxBytes = 10 * 1024 * 1024

// Seven 32-byte tokens plus delimiters fit the existing VARCHAR(255) job field.
const videoInputTokenStorageLimit = 7

var (
	ErrVideoInputTooLarge        = errors.New("video input is too large")
	ErrVideoInputUnsupportedType = errors.New("video input type is not supported")
	ErrVideoInputNotFound        = errors.New("video input not found")
)

type VideoInput struct {
	Token       string
	ContentType string
	Size        int64
	URL         string
	Data        []byte
}

type videoInputEntry struct {
	contentType string
	createdAt   time.Time
	terminalAt  *time.Time
}

type VideoInputStore struct {
	root string
	port int

	mu      sync.RWMutex
	entries map[string]videoInputEntry
}

func MarkVideoInputTerminal(store *VideoInputStore, token string, at time.Time) error {
	if store == nil || strings.TrimSpace(token) == "" {
		return nil
	}
	for _, item := range strings.Split(token, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if err := store.MarkTerminal(item, at); err != nil && !errors.Is(err, ErrVideoInputNotFound) {
			return err
		}
	}
	return nil
}

func NewVideoInputStore(dataDir string, port int) *VideoInputStore {
	return &VideoInputStore{root: filepath.Join(dataDir, "video-inputs"), port: port, entries: make(map[string]videoInputEntry)}
}

func (s *VideoInputStore) Save(reader io.Reader) (*VideoInput, error) {
	if s == nil || reader == nil {
		return nil, ErrVideoInputNotFound
	}
	data, err := io.ReadAll(io.LimitReader(reader, VideoInputMaxBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > VideoInputMaxBytes {
		return nil, ErrVideoInputTooLarge
	}
	contentType := http.DetectContentType(data)
	if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/webp" {
		return nil, ErrVideoInputUnsupportedType
	}
	if err := os.MkdirAll(s.root, 0o700); err != nil {
		return nil, err
	}
	token, err := newVideoInputToken()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(s.root, token)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return nil, err
	}
	now := time.Now()
	s.mu.Lock()
	s.entries[token] = videoInputEntry{contentType: contentType, createdAt: now}
	s.mu.Unlock()
	return &VideoInput{Token: token, ContentType: contentType, Size: int64(len(data)), URL: s.InternalURL(token), Data: data}, nil
}

func (s *VideoInputStore) Open(token string) (*VideoInput, error) {
	if s == nil || !validVideoInputToken(token) {
		return nil, ErrVideoInputNotFound
	}
	path := filepath.Join(s.root, token)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrVideoInputNotFound
	}
	if err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(data)
	s.mu.RLock()
	if entry, ok := s.entries[token]; ok && entry.contentType != "" {
		contentType = entry.contentType
	}
	s.mu.RUnlock()
	return &VideoInput{Token: token, ContentType: contentType, Size: int64(len(data)), URL: s.InternalURL(token), Data: data}, nil
}

func (s *VideoInputStore) MarkTerminal(token string, at time.Time) error {
	if s == nil || !validVideoInputToken(token) {
		return ErrVideoInputNotFound
	}
	if _, err := os.Stat(filepath.Join(s.root, token)); errors.Is(err, os.ErrNotExist) {
		return ErrVideoInputNotFound
	} else if err != nil {
		return err
	}
	s.mu.Lock()
	entry := s.entries[token]
	entry.terminalAt = &at
	if entry.createdAt.IsZero() {
		entry.createdAt = at
	}
	s.entries[token] = entry
	s.mu.Unlock()
	return nil
}

func (s *VideoInputStore) Cleanup(now time.Time) (int, error) {
	if s == nil {
		return 0, nil
	}
	if err := os.MkdirAll(s.root, 0o700); err != nil {
		return 0, err
	}
	removed := 0
	s.mu.Lock()
	for token, entry := range s.entries {
		expired := entry.terminalAt != nil && !now.Before(entry.terminalAt.Add(time.Hour))
		orphan := entry.terminalAt == nil && !now.Before(entry.createdAt.Add(24*time.Hour))
		if !expired && !orphan {
			continue
		}
		if err := os.Remove(filepath.Join(s.root, token)); err != nil && !errors.Is(err, os.ErrNotExist) {
			s.mu.Unlock()
			return removed, err
		}
		delete(s.entries, token)
		removed++
	}
	s.mu.Unlock()

	entries, err := os.ReadDir(s.root)
	if err != nil {
		return removed, err
	}
	for _, item := range entries {
		if item.IsDir() || !validVideoInputToken(item.Name()) {
			continue
		}
		info, statErr := item.Info()
		if statErr != nil || now.Before(info.ModTime().Add(24*time.Hour)) {
			continue
		}
		if err := os.Remove(filepath.Join(s.root, item.Name())); err != nil && !errors.Is(err, os.ErrNotExist) {
			return removed, err
		}
		removed++
	}
	return removed, nil
}

func (s *VideoInputStore) InternalURL(token string) string {
	return fmt.Sprintf("http://127.0.0.1:%d/internal/video-inputs/%s", s.port, url.PathEscape(token))
}

func (s *VideoInputStore) TokenFromURL(raw string) string {
	if s == nil {
		return ""
	}
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host != fmt.Sprintf("127.0.0.1:%d", s.port) || !strings.HasPrefix(u.Path, "/internal/video-inputs/") {
		return ""
	}
	token := strings.TrimPrefix(u.Path, "/internal/video-inputs/")
	if !validVideoInputToken(token) {
		return ""
	}
	return token
}

func (s *VideoInputStore) TokensFromVideoRequest(body []byte) []string {
	if s == nil {
		return nil
	}
	tokens := make([]string, 0, 6)
	seen := make(map[string]struct{})
	for _, raw := range LeoVideoReferenceURLs(body) {
		if len(tokens) >= videoInputTokenStorageLimit {
			break
		}
		token := s.TokenFromURL(raw)
		if token == "" {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		tokens = append(tokens, token)
	}
	return tokens
}

func newVideoInputToken() (string, error) {
	var raw [24]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func validVideoInputToken(token string) bool {
	if len(token) != 32 {
		return false
	}
	for _, r := range token {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_' {
			return false
		}
	}
	return true
}
