package service

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVideoInputStoreValidatesImageContentAndBuildsOpaqueURL(t *testing.T) {
	store := NewVideoInputStore(t.TempDir(), 8080)
	input, err := store.Save(bytes.NewReader(videoPNGBytes()))
	require.NoError(t, err)
	require.Equal(t, "image/png", input.ContentType)
	require.Len(t, input.Token, 32)
	require.Contains(t, input.URL, "/internal/video-inputs/")
	require.NotContains(t, input.URL, "png")

	opened, err := store.Open(input.Token)
	require.NoError(t, err)
	require.Equal(t, "image/png", opened.ContentType)
	require.Equal(t, videoPNGBytes(), opened.Data)

	_, err = store.Save(bytes.NewReader([]byte("<svg></svg>")))
	require.ErrorIs(t, err, ErrVideoInputUnsupportedType)
}

func TestVideoInputStoreRejectsOversizedAndCleansTerminalAndOrphanFiles(t *testing.T) {
	root := t.TempDir()
	store := NewVideoInputStore(root, 8080)
	tooLarge := bytes.Repeat([]byte("x"), VideoInputMaxBytes+1)
	_, err := store.Save(bytes.NewReader(tooLarge))
	require.ErrorIs(t, err, ErrVideoInputTooLarge)

	input, err := store.Save(bytes.NewReader(videoJPEGBytes()))
	require.NoError(t, err)
	terminalAt := time.Now().Add(-2 * time.Hour)
	require.NoError(t, store.MarkTerminal(input.Token, terminalAt))

	orphan := filepath.Join(root, "video-inputs", "01234567890123456789012345678901")
	require.NoError(t, os.WriteFile(orphan, videoPNGBytes(), 0o600))
	old := time.Now().Add(-25 * time.Hour)
	require.NoError(t, os.Chtimes(orphan, old, old))

	removed, err := store.Cleanup(time.Now())
	require.NoError(t, err)
	require.Equal(t, 2, removed)
	_, err = store.Open(input.Token)
	require.ErrorIs(t, err, ErrVideoInputNotFound)
	_, err = os.Stat(orphan)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestVideoInputStoreTracksAllLocalVideoGuidanceURLs(t *testing.T) {
	store := NewVideoInputStore(t.TempDir(), 8080)
	first, err := store.Save(bytes.NewReader(videoPNGBytes()))
	require.NoError(t, err)
	second, err := store.Save(bytes.NewReader(videoJPEGBytes()))
	require.NoError(t, err)
	third, err := store.Save(bytes.NewReader(videoPNGBytes()))
	require.NoError(t, err)
	body := []byte(`{"start_frame_url":"` + first.URL + `","image_urls":["` + second.URL + `"],"guidances":{"image_reference":[{"image":{"url":"` + third.URL + `"}},{"image":{"url":"` + second.URL + `"}}]}}`)

	tokens := store.TokensFromVideoRequest(body)

	require.Equal(t, []string{first.Token, second.Token, third.Token}, tokens)
	require.NoError(t, MarkVideoInputTerminal(store, strings.Join(tokens, ","), time.Now().Add(-2*time.Hour)))
	removed, err := store.Cleanup(time.Now())
	require.NoError(t, err)
	require.Equal(t, 3, removed)
}

func TestVideoInputStoreTracksMediaAndAudioGuidanceURLs(t *testing.T) {
	store := NewVideoInputStore(t.TempDir(), 8080)
	video, err := store.Save(bytes.NewReader(videoPNGBytes()))
	require.NoError(t, err)
	audio, err := store.Save(bytes.NewReader(videoJPEGBytes()))
	require.NoError(t, err)
	body := []byte(`{"guidances":{"video_reference_base":[{"video":{"url":"` + video.URL + `"}}],"audio_reference":[{"audio":{"url":"` + audio.URL + `"}}]}}`)

	require.ElementsMatch(t, []string{video.Token, audio.Token}, store.TokensFromVideoRequest(body))
}

func TestVideoInputStoreValidatesReferenceVideoAndAudioFormats(t *testing.T) {
	store := NewVideoInputStore(t.TempDir(), 8080)

	video, err := store.SaveMedia(bytes.NewReader(referenceMP4Bytes()), VideoInputKindVideo, "reference.mp4")
	require.NoError(t, err)
	require.Equal(t, VideoInputKindVideo, video.Kind)
	require.Equal(t, "video/mp4", video.ContentType)

	mov, err := store.SaveMedia(bytes.NewReader(referenceMOVBytes()), VideoInputKindVideo, "reference.mov")
	require.NoError(t, err)
	require.Equal(t, "video/quicktime", mov.ContentType)

	_, err = store.SaveMedia(bytes.NewReader(referenceMP4Bytes()), VideoInputKindVideo, "reference.webm")
	require.ErrorIs(t, err, ErrVideoInputUnsupportedType)

	mp3, err := store.SaveMedia(bytes.NewReader(referenceMP3Bytes(2)), VideoInputKindAudio, "reference.mp3")
	require.NoError(t, err)
	require.Equal(t, "audio/mpeg", mp3.ContentType)
	_, err = store.SaveMedia(bytes.NewReader(referenceMP3Bytes(1)), VideoInputKindAudio, "reference.mp3")
	require.ErrorIs(t, err, ErrVideoInputUnsupportedDuration)
	_, err = store.SaveMedia(bytes.NewReader(referenceMP3Bytes(31)), VideoInputKindAudio, "reference.mp3")
	require.ErrorIs(t, err, ErrVideoInputUnsupportedDuration)

	wav, err := store.SaveMedia(bytes.NewReader(referenceWAVBytes(2, 16)), VideoInputKindAudio, "reference.wav")
	require.NoError(t, err)
	require.Equal(t, "audio/wav", wav.ContentType)

	_, err = store.SaveMedia(bytes.NewReader(referenceWAVBytes(2, 8)), VideoInputKindAudio, "reference.wav")
	require.ErrorIs(t, err, ErrVideoInputUnsupportedType)
	_, err = store.SaveMedia(bytes.NewReader(referenceWAVBytes(1, 16)), VideoInputKindAudio, "reference.wav")
	require.ErrorIs(t, err, ErrVideoInputUnsupportedDuration)
	_, err = store.SaveMedia(bytes.NewReader(referenceWAVBytes(31, 16)), VideoInputKindAudio, "reference.wav")
	require.ErrorIs(t, err, ErrVideoInputUnsupportedDuration)
}

func TestVideoInputStoreReferenceAudioSizeAndRestartMIME(t *testing.T) {
	root := t.TempDir()
	store := NewVideoInputStore(root, 8080)
	tooLarge := bytes.Repeat([]byte("ID3"), AudioReferenceMaxBytes/3+1)
	_, err := store.SaveMedia(bytes.NewReader(tooLarge), VideoInputKindAudio, "reference.mp3")
	require.ErrorIs(t, err, ErrVideoInputTooLarge)

	input, err := store.SaveMedia(bytes.NewReader(referenceMP4Bytes()), VideoInputKindVideo, "reference.mp4")
	require.NoError(t, err)
	restarted := NewVideoInputStore(root, 8080)
	opened, err := restarted.Open(input.Token)
	require.NoError(t, err)
	require.Equal(t, "video/mp4", opened.ContentType)
}

func videoPNGBytes() []byte {
	return []byte("\x89PNG\r\n\x1a\nvideo")
}

func videoJPEGBytes() []byte {
	return []byte("\xff\xd8\xff\xe0video\xff\xd9")
}

func referenceMP4Bytes() []byte {
	data := make([]byte, 16)
	binary.BigEndian.PutUint32(data[0:4], uint32(len(data)))
	copy(data[4:8], "ftyp")
	copy(data[8:12], "isom")
	return data
}

func referenceMOVBytes() []byte {
	data := referenceMP4Bytes()
	copy(data[8:12], "qt  ")
	return data
}

func referenceMP3Bytes(seconds int) []byte {
	const frameLength = 417
	frameCount := (seconds*44100 + 1151) / 1152
	data := make([]byte, frameLength*frameCount)
	for offset := 0; offset < len(data); offset += frameLength {
		data[offset] = 0xff
		data[offset+1] = 0xfb
		data[offset+2] = 0x90
		data[offset+3] = 0x64
	}
	return data[:frameCount*frameLength]
}

func referenceWAVBytes(seconds, bits uint16) []byte {
	const sampleRate = uint32(44100)
	channels := uint16(1)
	byteRate := sampleRate * uint32(channels) * uint32(bits/8)
	dataSize := byteRate * uint32(seconds)
	data := make([]byte, 44+dataSize)
	copy(data[0:4], "RIFF")
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(data)-8))
	copy(data[8:12], "WAVE")
	copy(data[12:16], "fmt ")
	binary.LittleEndian.PutUint32(data[16:20], 16)
	binary.LittleEndian.PutUint16(data[20:22], 1)
	binary.LittleEndian.PutUint16(data[22:24], channels)
	binary.LittleEndian.PutUint32(data[24:28], sampleRate)
	binary.LittleEndian.PutUint32(data[28:32], byteRate)
	binary.LittleEndian.PutUint16(data[32:34], channels*bits/8)
	binary.LittleEndian.PutUint16(data[34:36], bits)
	copy(data[36:40], "data")
	binary.LittleEndian.PutUint32(data[40:44], dataSize)
	return data
}
