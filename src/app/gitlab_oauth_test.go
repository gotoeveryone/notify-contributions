package app

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestExchangeGitLabToken(t *testing.T) {
	cfg := &Config{
		GitLabOAuthClientID:     "client-id",
		GitLabOAuthClientSecret: "client-secret",
		GitLabOAuthRedirectURI:  "https://example.com/callback",
		GitLabOAuthRefreshToken: "old-refresh-token",
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "https://gitlab.com/oauth/token", func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost || r.URL.Path != "/oauth/token" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		for key, want := range map[string]string{
			"client_id": "client-id", "client_secret": "client-secret",
			"redirect_uri":  "https://example.com/callback",
			"refresh_token": "old-refresh-token", "grant_type": "refresh_token",
		} {
			if got := r.PostForm.Get(key); got != want {
				t.Errorf("%s = %q, want %q", key, got, want)
			}
		}
		return httpmock.NewStringResponse(http.StatusOK, `{"access_token":"new-access-token","refresh_token":"new-refresh-token"}`), nil
	})

	tokens, err := exchangeGitLabToken(context.Background(), cfg, &http.Client{}, "https://gitlab.com/oauth/token")
	if err != nil {
		t.Fatal(err)
	}
	if tokens.AccessToken != "new-access-token" || tokens.RefreshToken != "new-refresh-token" {
		t.Fatal("OAuth response tokens were not decoded")
	}
}

func TestExchangeGitLabTokenReportsSafeOAuthError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(http.MethodPost, "https://gitlab.com/oauth/token",
		httpmock.NewStringResponder(http.StatusBadRequest, `{"error":"invalid_grant","error_description":"secret marker"}`))
	_, err := exchangeGitLabToken(context.Background(), &Config{}, &http.Client{}, "https://gitlab.com/oauth/token")
	if err == nil || !strings.Contains(err.Error(), "invalid_grant") || strings.Contains(err.Error(), "secret marker") {
		t.Fatalf("unexpected OAuth error: %v", err)
	}
}

func TestLocalGitLabRefreshTokenIsReused(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".gitlab-oauth-refresh-token")
	initial, err := readLocalGitLabRefreshToken(path, "initial-token")
	if err != nil || initial != "initial-token" {
		t.Fatalf("initial token = %q, err = %v", initial, err)
	}
	if err := writeLocalGitLabRefreshToken(path, "rotated-token"); err != nil {
		t.Fatal(err)
	}
	got, err := readLocalGitLabRefreshToken(path, "initial-token")
	if err != nil || got != "rotated-token" {
		t.Fatalf("stored token = %q, err = %v", got, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("token file permissions = %o, want 600", info.Mode().Perm())
	}
}
