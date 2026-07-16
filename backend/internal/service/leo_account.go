package service

import (
	"net/url"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

var LeoDefaultVideoModelIDs = []string{"seedance-2.0", "seedance-2.0-fast"}

func (a *Account) IsLeo() bool {
	return a != nil && a.Platform == PlatformLeo
}

func (a *Account) GetLeoAPIKey() string {
	if !a.IsLeo() || a.Type != AccountTypeAPIKey {
		return ""
	}
	return strings.TrimSpace(a.GetCredential("api_key"))
}

func (a *Account) GetLeoBaseURL() string {
	if !a.IsLeo() {
		return ""
	}
	baseURL, err := NormalizeLeoBaseURL(a.GetCredential("base_url"))
	if err != nil {
		return ""
	}
	return baseURL
}

func NormalizeLeoBaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	u, err := url.Parse(raw)
	if err != nil || u == nil {
		return "", invalidLeoCredentials("base_url must be a valid URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", invalidLeoCredentials("base_url scheme must be http or https")
	}
	if u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" {
		return "", invalidLeoCredentials("base_url must contain only scheme, host, port, and /v1 path")
	}
	u.Path = strings.TrimSuffix(u.Path, "/")
	if u.Path != "/v1" || u.RawPath != "" {
		return "", invalidLeoCredentials("base_url path must be /v1")
	}
	return u.String(), nil
}

func BuildLeoVideosGenerationsURL(baseURL string) (string, error) {
	baseURL, err := NormalizeLeoBaseURL(baseURL)
	if err != nil {
		return "", err
	}
	return baseURL + "/videos/generations", nil
}

func BuildLeoHealthURL(baseURL string) (string, error) {
	baseURL, err := NormalizeLeoBaseURL(baseURL)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(baseURL, "/v1") + "/health", nil
}

func ValidateLeoAccountCredentials(platform, accountType string, credentials map[string]any) error {
	if platform != PlatformLeo {
		return nil
	}
	if accountType != AccountTypeAPIKey {
		return invalidLeoCredentials("leo accounts must use API Key authentication")
	}
	baseURL, ok := credentials["base_url"].(string)
	if !ok {
		return invalidLeoCredentials("base_url is required for leo accounts")
	}
	if _, err := NormalizeLeoBaseURL(baseURL); err != nil {
		return err
	}
	apiKey, ok := credentials["api_key"].(string)
	if !ok || strings.TrimSpace(apiKey) == "" {
		return invalidLeoCredentials("api_key is required for leo accounts")
	}
	mapping, ok := credentials["model_mapping"].(map[string]any)
	if !ok || len(mapping) == 0 {
		return invalidLeoCredentials("model_mapping must contain at least one leo video model")
	}
	for from, rawTo := range mapping {
		to, ok := rawTo.(string)
		if !ok || strings.TrimSpace(from) == "" || strings.TrimSpace(to) == "" {
			return invalidLeoCredentials("model_mapping keys and values must be non-empty strings")
		}
	}
	return nil
}

func invalidLeoCredentials(message string) error {
	return infraerrors.BadRequest("LEO_ACCOUNT_CREDENTIALS_INVALID", message)
}
