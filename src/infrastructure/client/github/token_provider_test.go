package github

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestGitHubAppTokenProviderToken(t *testing.T) {
	privateKey := generatePrivateKeyPEM(t)
	provider, err := NewGitHubAppTokenProvider("12345", "67890", privateKey)
	if err != nil {
		t.Fatal(err)
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	requestCount := 0
	httpmock.RegisterResponder("POST", "https://api.github.com/app/installations/67890/access_tokens",
		func(req *http.Request) (*http.Response, error) {
			requestCount++
			if req.Header.Get("Authorization") == "" {
				t.Error("Authorization header is empty")
			}

			return httpmock.NewStringResponse(201, `{
				"token": "installation_token",
				"expires_at": "2099-01-01T00:00:00Z"
			}`), nil
		})

	token, err := provider.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if token != "installation_token" {
		t.Errorf("Token is not matched, actual: [%s], expected: [%s]", token, "installation_token")
	}

	token, err = provider.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if token != "installation_token" {
		t.Errorf("Token is not matched, actual: [%s], expected: [%s]", token, "installation_token")
	}
	if requestCount != 1 {
		t.Errorf("Request count is not matched, actual: [%d], expected: [%d]", requestCount, 1)
	}
}

func TestGitHubAppTokenProviderTokenExpiredCache(t *testing.T) {
	privateKey := generatePrivateKeyPEM(t)
	provider, err := NewGitHubAppTokenProvider("12345", "67890", privateKey)
	if err != nil {
		t.Fatal(err)
	}
	provider.token = "cached_token"
	provider.expiresAt = time.Now().Add(30 * time.Second)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://api.github.com/app/installations/67890/access_tokens",
		httpmock.NewStringResponder(201, `{
			"token": "refreshed_token",
			"expires_at": "2099-01-01T00:00:00Z"
		}`))

	token, err := provider.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if token != "refreshed_token" {
		t.Errorf("Token is not matched, actual: [%s], expected: [%s]", token, "refreshed_token")
	}
}

func generatePrivateKeyPEM(t *testing.T) string {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))
}
