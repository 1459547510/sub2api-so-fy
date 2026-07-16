package service

import (
	"context"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNormalizeLeoBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "docker http", raw: "http://leostudio:8000/v1/", want: "http://leostudio:8000/v1"},
		{name: "private ip", raw: "http://10.0.0.20:8000/v1", want: "http://10.0.0.20:8000/v1"},
		{name: "https", raw: "https://leo.internal.example/v1", want: "https://leo.internal.example/v1"},
		{name: "missing v1", raw: "http://leostudio:8000", wantErr: true},
		{name: "userinfo", raw: "http://user:pass@leostudio:8000/v1", wantErr: true},
		{name: "query", raw: "http://leostudio:8000/v1?x=1", wantErr: true},
		{name: "fragment", raw: "http://leostudio:8000/v1#x", wantErr: true},
		{name: "bad scheme", raw: "file:///tmp/v1", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeLeoBaseURL(tt.raw)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidateLeoAccountCredentials(t *testing.T) {
	valid := map[string]any{
		"base_url": "http://leostudio:8000/v1",
		"api_key":  "leo-secret",
		"model_mapping": map[string]any{
			"seedance-2.0": "seedance-2.0",
		},
	}

	require.NoError(t, ValidateLeoAccountCredentials(PlatformLeo, AccountTypeAPIKey, valid))
	require.Error(t, ValidateLeoAccountCredentials(PlatformLeo, AccountTypeOAuth, valid))
	require.Error(t, ValidateLeoAccountCredentials(PlatformLeo, AccountTypeAPIKey, map[string]any{}))
	require.Error(t, ValidateLeoAccountCredentials(PlatformLeo, AccountTypeAPIKey, map[string]any{
		"base_url":      valid["base_url"],
		"api_key":       valid["api_key"],
		"model_mapping": map[string]any{"seedance-2.0": 2},
	}))
	require.NoError(t, ValidateLeoAccountCredentials(PlatformOpenAI, AccountTypeOAuth, nil))
}

func TestLeoAccountContract(t *testing.T) {
	account := &Account{
		Platform: PlatformLeo,
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"base_url": "http://leostudio:8000/v1/",
			"api_key":  " leo-secret ",
		},
	}

	require.True(t, account.IsLeo())
	require.True(t, account.IsOpenAICompatible())
	require.Equal(t, "leo-secret", account.GetLeoAPIKey())
	require.Equal(t, "http://leostudio:8000/v1", account.GetLeoBaseURL())

	videosURL, err := BuildLeoVideosGenerationsURL(account.GetLeoBaseURL())
	require.NoError(t, err)
	require.Equal(t, "http://leostudio:8000/v1/videos/generations", videosURL)

	healthURL, err := BuildLeoHealthURL(account.GetLeoBaseURL())
	require.NoError(t, err)
	require.Equal(t, "http://leostudio:8000/health", healthURL)
}

func TestAdminServiceCreateAccountRejectsInvalidLeoCredentials(t *testing.T) {
	svc := &adminServiceImpl{}

	account, err := svc.CreateAccount(context.Background(), &CreateAccountInput{
		Platform:             PlatformLeo,
		Type:                 AccountTypeAPIKey,
		Credentials:          map[string]any{},
		SkipDefaultGroupBind: true,
	})

	require.Nil(t, account)
	require.Equal(t, 400, infraerrors.Code(err))
}

func TestDefaultModelsListCandidateIDsLeo(t *testing.T) {
	require.Equal(t, []string{"seedance-2.0", "seedance-2.0-fast"}, defaultModelsListCandidateIDs(PlatformLeo))
}
