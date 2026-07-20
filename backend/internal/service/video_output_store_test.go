package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

var minimalMP4 = []byte{0, 0, 0, 16, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm', 0, 0, 0, 0}

func TestVideoOutputStoreDownloadsValidMP4AndRewritesResult(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.Header().Set("Content-Type", "video/mp4")
		_, _ = w.Write(minimalMP4)
	}))
	defer server.Close()
	store := NewVideoOutputStore(t.TempDir())
	result := json.RawMessage(`{"data":[{"mp4_url":"` + server.URL + `/video.mp4"}],"provider":{"duration":8}}`)

	rewritten, err := store.Save(context.Background(), "vidjob_valid", result)

	require.NoError(t, err)
	require.Equal(t, 1, calls)
	require.Equal(t, server.URL+"/video.mp4", gjson.GetBytes(rewritten, "data.0.source_url").String())
	require.Equal(t, "/v1/videos/jobs/vidjob_valid/content", gjson.GetBytes(rewritten, "data.0.mp4_url").String())
	file, err := store.Open("vidjob_valid")
	require.NoError(t, err)
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	require.NoError(t, err)
	require.Equal(t, int64(len(minimalMP4)), info.Size())

	_, err = store.Save(context.Background(), "vidjob_valid", rewritten)
	require.NoError(t, err)
	require.Equal(t, 1, calls)
}

func TestVideoOutputStoreRejectsMissingURLAndInvalidMP4(t *testing.T) {
	store := NewVideoOutputStore(t.TempDir())
	_, err := store.Save(context.Background(), "vidjob_empty", json.RawMessage(`{"data":[]}`))
	require.ErrorIs(t, err, ErrVideoOutputURLMissing)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("not-an-mp4"))
	}))
	defer server.Close()
	result := json.RawMessage(`{"data":[{"url":"` + server.URL + `/video"}]}`)
	_, err = store.Save(context.Background(), "vidjob_invalid", result)
	require.ErrorIs(t, err, ErrVideoOutputInvalid)
	_, err = store.Open("vidjob_invalid")
	require.True(t, errors.Is(err, ErrVideoOutputNotFound))
}
