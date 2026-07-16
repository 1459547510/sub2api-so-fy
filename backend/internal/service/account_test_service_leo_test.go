//go:build unit

package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/tlsfingerprint"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type leoAccountTestRepo struct {
	mockAccountRepoForGemini
	setErrorID      int64
	setErrorMessage string
}

func (r *leoAccountTestRepo) SetError(_ context.Context, id int64, message string) error {
	r.setErrorID = id
	r.setErrorMessage = message
	return nil
}

type leoAccountTestHTTPUpstream struct{}

func (leoAccountTestHTTPUpstream) Do(req *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func (u leoAccountTestHTTPUpstream) DoWithTLS(req *http.Request, proxyURL string, accountID int64, concurrency int, _ *tlsfingerprint.Profile) (*http.Response, error) {
	return u.Do(req, proxyURL, accountID, concurrency)
}

func TestAccountTestService_Leo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success uses bearer authenticated health endpoint", func(t *testing.T) {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/health", r.URL.Path)
			require.Equal(t, "Bearer leo-secret", r.Header.Get("Authorization"))
			require.Equal(t, "application/json", r.Header.Get("Accept"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}))
		defer upstream.Close()

		repo, svc, account := newLeoAccountTestService(upstream.URL + "/v1")
		rec, c := newLeoAccountTestContext(account.ID)

		err := svc.TestAccountConnection(c, account.ID, "", "", AccountTestModeDefault)

		require.NoError(t, err)
		require.Contains(t, rec.Body.String(), `"type":"content"`)
		require.Contains(t, rec.Body.String(), `"type":"test_complete"`)
		require.Contains(t, rec.Body.String(), `"success":true`)
		require.Zero(t, repo.setErrorID)
	})

	t.Run("unauthorized marks account error without leaking secrets", func(t *testing.T) {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"leo-secret upstream detail"}`))
		}))
		defer upstream.Close()

		repo, svc, account := newLeoAccountTestService(upstream.URL + "/v1")
		rec, c := newLeoAccountTestContext(account.ID)

		err := svc.TestAccountConnection(c, account.ID, "", "", AccountTestModeDefault)

		require.Error(t, err)
		require.Equal(t, account.ID, repo.setErrorID)
		require.NotContains(t, repo.setErrorMessage, "leo-secret")
		require.NotContains(t, rec.Body.String(), "leo-secret")
		require.Contains(t, rec.Body.String(), "HTTP 401")
	})

	t.Run("invalid base URL fails before transport", func(t *testing.T) {
		_, svc, account := newLeoAccountTestService("http://leostudio:8000")
		rec, c := newLeoAccountTestContext(account.ID)

		err := svc.TestAccountConnection(c, account.ID, "", "", AccountTestModeDefault)

		require.Error(t, err)
		require.Contains(t, rec.Body.String(), "Invalid Leo base URL")
	})

	t.Run("transport failure is redacted", func(t *testing.T) {
		upstream := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		baseURL := upstream.URL + "/v1"
		upstream.Close()

		_, svc, account := newLeoAccountTestService(baseURL)
		rec, c := newLeoAccountTestContext(account.ID)

		err := svc.TestAccountConnection(c, account.ID, "", "", AccountTestModeDefault)

		require.Error(t, err)
		require.Contains(t, rec.Body.String(), "Leo health check request failed")
		require.NotContains(t, rec.Body.String(), "leo-secret")
	})
}

func newLeoAccountTestService(baseURL string) (*leoAccountTestRepo, *AccountTestService, *Account) {
	account := &Account{
		ID:          42,
		Platform:    PlatformLeo,
		Type:        AccountTypeAPIKey,
		Concurrency: 3,
		Credentials: map[string]any{
			"base_url": baseURL,
			"api_key":  "leo-secret",
			"model_mapping": map[string]any{
				"seedance-2.0": "seedance-2.0",
			},
		},
	}
	repo := &leoAccountTestRepo{
		mockAccountRepoForGemini: mockAccountRepoForGemini{
			accountsByID: map[int64]*Account{account.ID: account},
		},
	}
	return repo, &AccountTestService{
		accountRepo:  repo,
		httpUpstream: leoAccountTestHTTPUpstream{},
	}, account
}

func newLeoAccountTestContext(accountID int64) (*httptest.ResponseRecorder, *gin.Context) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/42/test", nil)
	return rec, c
}
