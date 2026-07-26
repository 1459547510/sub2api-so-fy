package service

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
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

const (
	VideoInputMaxBytes       = 10 * 1024 * 1024
	VideoReferenceMaxBytes   = 100 * 1024 * 1024
	AudioReferenceMaxBytes   = 15 * 1024 * 1024
	AudioReferenceMinSeconds = 2
	AudioReferenceMaxSeconds = 30
)

type VideoInputKind string

const (
	VideoInputKindImage VideoInputKind = "image"
	VideoInputKindVideo VideoInputKind = "video"
	VideoInputKindAudio VideoInputKind = "audio"
)

// Seven 32-byte tokens plus delimiters fit the existing VARCHAR(255) job field.
const videoInputTokenStorageLimit = 7

var (
	ErrVideoInputTooLarge            = errors.New("video input is too large")
	ErrVideoInputUnsupportedType     = errors.New("video input type is not supported")
	ErrVideoInputUnsupportedDuration = errors.New("audio input duration must be between 2 and 30 seconds")
	ErrVideoInputNotFound            = errors.New("video input not found")
)

type VideoInput struct {
	Token       string
	Kind        VideoInputKind
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
	return s.SaveMedia(reader, VideoInputKindImage, "")
}

func (s *VideoInputStore) SaveMedia(reader io.Reader, kind VideoInputKind, filename string) (*VideoInput, error) {
	if s == nil || reader == nil {
		return nil, ErrVideoInputNotFound
	}
	limit, ok := videoInputMaxBytes(kind)
	if !ok {
		return nil, ErrVideoInputUnsupportedType
	}
	data, err := io.ReadAll(io.LimitReader(reader, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, ErrVideoInputTooLarge
	}
	contentType, err := validateVideoInputData(data, kind, filename)
	if err != nil {
		return nil, err
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
	return &VideoInput{Token: token, Kind: kind, ContentType: contentType, Size: int64(len(data)), URL: s.InternalURL(token), Data: data}, nil
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
	contentType := detectStoredMediaType(data)
	s.mu.RLock()
	if entry, ok := s.entries[token]; ok && entry.contentType != "" {
		contentType = entry.contentType
	}
	s.mu.RUnlock()
	return &VideoInput{Token: token, ContentType: contentType, Size: int64(len(data)), URL: s.InternalURL(token), Data: data}, nil
}

func videoInputMaxBytes(kind VideoInputKind) (int64, bool) {
	switch kind {
	case VideoInputKindImage:
		return VideoInputMaxBytes, true
	case VideoInputKindVideo:
		return VideoReferenceMaxBytes, true
	case VideoInputKindAudio:
		return AudioReferenceMaxBytes, true
	default:
		return 0, false
	}
}

func validateVideoInputData(data []byte, kind VideoInputKind, filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(filename)))
	switch kind {
	case VideoInputKindImage:
		contentType := http.DetectContentType(data)
		if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/webp" {
			return "", ErrVideoInputUnsupportedType
		}
		return contentType, nil
	case VideoInputKindVideo:
		if !isISOBaseMedia(data) || (ext != ".mp4" && ext != ".mov") {
			return "", ErrVideoInputUnsupportedType
		}
		if ext == ".mov" || isQuickTimeBrand(data) {
			return "video/quicktime", nil
		}
		return "video/mp4", nil
	case VideoInputKindAudio:
		switch ext {
		case ".mp3":
			duration, ok := parseMP3Duration(data)
			if !ok {
				return "", ErrVideoInputUnsupportedType
			}
			if duration < AudioReferenceMinSeconds || duration > AudioReferenceMaxSeconds {
				return "", ErrVideoInputUnsupportedDuration
			}
			return "audio/mpeg", nil
		case ".wav":
			duration, ok := parsePCM16Or24WAVDuration(data)
			if !ok {
				return "", ErrVideoInputUnsupportedType
			}
			if duration < AudioReferenceMinSeconds || duration > AudioReferenceMaxSeconds {
				return "", ErrVideoInputUnsupportedDuration
			}
			return "audio/wav", nil
		default:
			return "", ErrVideoInputUnsupportedType
		}
	default:
		return "", ErrVideoInputUnsupportedType
	}
}

func detectStoredMediaType(data []byte) string {
	if contentType := http.DetectContentType(data); strings.HasPrefix(contentType, "image/") {
		return contentType
	}
	if isISOBaseMedia(data) {
		if isQuickTimeBrand(data) {
			return "video/quicktime"
		}
		return "video/mp4"
	}
	if _, ok := parsePCM16Or24WAVDuration(data); ok {
		return "audio/wav"
	}
	if isMP3(data) {
		return "audio/mpeg"
	}
	return "application/octet-stream"
}

func isISOBaseMedia(data []byte) bool {
	return len(data) >= 12 && string(data[4:8]) == "ftyp"
}

func isQuickTimeBrand(data []byte) bool {
	return len(data) >= 12 && string(data[8:12]) == "qt  "
}

func isMP3(data []byte) bool {
	if bytes.HasPrefix(data, []byte("ID3")) {
		return true
	}
	limit := len(data)
	if limit > 4096 {
		limit = 4096
	}
	for index := 0; index+1 < limit; index++ {
		if data[index] == 0xff && data[index+1]&0xe0 == 0xe0 && data[index+1]&0x06 != 0 {
			return true
		}
	}
	return false
}

func parseMP3Duration(data []byte) (float64, bool) {
	start := 0
	if bytes.HasPrefix(data, []byte("ID3")) && len(data) >= 10 {
		size := int(data[6]&0x7f)<<21 | int(data[7]&0x7f)<<14 | int(data[8]&0x7f)<<7 | int(data[9]&0x7f)
		start = 10 + size
	}
	if start >= len(data) {
		return 0, false
	}
	frames := 0
	var duration float64
	for offset := start; offset+4 <= len(data); {
		header := binary.BigEndian.Uint32(data[offset : offset+4])
		frameLength, frameSamples, sampleRate, ok := mp3FrameInfo(header)
		if !ok || offset+frameLength > len(data) {
			if frames == 0 {
				offset++
				continue
			}
			break
		}
		frames++
		duration += float64(frameSamples) / float64(sampleRate)
		offset += frameLength
	}
	if frames == 0 {
		return 0, false
	}
	return duration, true
}

func mp3FrameInfo(header uint32) (int, int, int, bool) {
	if header>>21 != 0x7ff {
		return 0, 0, 0, false
	}
	version := (header >> 19) & 0x3
	layer := (header >> 17) & 0x3
	bitrateIndex := (header >> 12) & 0xf
	sampleRateIndex := (header >> 10) & 0x3
	if version == 1 || layer != 1 || bitrateIndex == 0 || bitrateIndex == 15 || sampleRateIndex == 3 {
		return 0, 0, 0, false
	}
	bitratesV1 := [...]int{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0}
	bitratesV2 := [...]int{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0}
	bitrate := bitratesV1[bitrateIndex]
	samplesPerFrame := 1152
	coefficient := 144
	if version != 3 {
		bitrate = bitratesV2[bitrateIndex]
		samplesPerFrame = 576
		coefficient = 72
	}
	sampleRates := [...]int{44100, 48000, 32000}
	sampleRate := sampleRates[sampleRateIndex]
	if version == 2 {
		sampleRate /= 2
	} else if version == 0 {
		sampleRate /= 4
	}
	padding := int((header >> 9) & 1)
	frameLength := coefficient*bitrate*1000/sampleRate + padding
	if frameLength < 4 {
		return 0, 0, 0, false
	}
	return frameLength, samplesPerFrame, sampleRate, true
}

func parsePCM16Or24WAVDuration(data []byte) (float64, bool) {
	if len(data) < 12 || string(data[0:4]) != "RIFF" || string(data[8:12]) != "WAVE" {
		return 0, false
	}
	var byteRate uint32
	var bitsPerSample uint16
	var dataSize uint32
	fmtFound, dataFound := false, false
	for offset := 12; offset+8 <= len(data); {
		chunkSize := int(binary.LittleEndian.Uint32(data[offset+4 : offset+8]))
		chunkStart := offset + 8
		if chunkStart > len(data) || chunkSize > len(data)-chunkStart {
			return 0, false
		}
		switch string(data[offset : offset+4]) {
		case "fmt ":
			if chunkSize < 16 {
				return 0, false
			}
			format := binary.LittleEndian.Uint16(data[chunkStart : chunkStart+2])
			byteRate = binary.LittleEndian.Uint32(data[chunkStart+8 : chunkStart+12])
			bitsPerSample = binary.LittleEndian.Uint16(data[chunkStart+14 : chunkStart+16])
			fmtFound = format == 1 && (bitsPerSample == 16 || bitsPerSample == 24) && byteRate > 0
		case "data":
			dataSize = uint32(chunkSize)
			dataFound = true
		}
		if fmtFound && dataFound {
			return float64(dataSize) / float64(byteRate), true
		}
		offset = chunkStart + chunkSize
		if chunkSize%2 != 0 {
			offset++
		}
	}
	return 0, false
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
