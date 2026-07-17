package service

import (
	"bytes"
	"os"
	"path/filepath"
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

func videoPNGBytes() []byte {
	return []byte("\x89PNG\r\n\x1a\nvideo")
}

func videoJPEGBytes() []byte {
	return []byte("\xff\xd8\xff\xe0video\xff\xd9")
}
