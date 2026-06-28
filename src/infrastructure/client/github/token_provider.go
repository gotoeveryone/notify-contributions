package github

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const githubAPIBaseURL = "https://api.github.com"

// TokenProvider provides a GitHub API access token.
type TokenProvider interface {
	Token(ctx context.Context) (string, error)
}

type staticTokenProvider struct {
	token string
}

func NewStaticTokenProvider(token string) TokenProvider {
	return &staticTokenProvider{token: token}
}

func (p *staticTokenProvider) Token(_ context.Context) (string, error) {
	return p.token, nil
}

type GitHubAppTokenProvider struct {
	appID          int64
	installationID string
	privateKey     *rsa.PrivateKey
	httpClient     *http.Client

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

func NewGitHubAppTokenProvider(appID string, installationID string, privateKeyPEM string) (*GitHubAppTokenProvider, error) {
	parsedAppID, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid github app id: %w", err)
	}

	privateKey, err := parsePrivateKey(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	return &GitHubAppTokenProvider{
		appID:          parsedAppID,
		installationID: installationID,
		privateKey:     privateKey,
		httpClient:     http.DefaultClient,
	}, nil
}

func (p *GitHubAppTokenProvider) Token(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.token != "" && time.Now().Before(p.expiresAt.Add(-1*time.Minute)) {
		return p.token, nil
	}

	jwt, err := p.createJWT(time.Now())
	if err != nil {
		return "", err
	}

	token, expiresAt, err := p.createInstallationToken(ctx, jwt)
	if err != nil {
		return "", err
	}

	p.token = token
	p.expiresAt = expiresAt

	return token, nil
}

func (p *GitHubAppTokenProvider) createJWT(now time.Time) (string, error) {
	header := map[string]string{
		"alg": "RS256",
		"typ": "JWT",
	}
	claims := map[string]int64{
		"iat": now.Add(-1 * time.Minute).Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
		"iss": p.appID,
	}

	encodedHeader, err := encodeSegment(header)
	if err != nil {
		return "", err
	}
	encodedClaims, err := encodeSegment(claims)
	if err != nil {
		return "", err
	}

	unsignedToken := encodedHeader + "." + encodedClaims
	hashed := sha256.Sum256([]byte(unsignedToken))
	signature, err := rsa.SignPKCS1v15(rand.Reader, p.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("github app jwt signing failed: %w", err)
	}

	return unsignedToken + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

func (p *GitHubAppTokenProvider) createInstallationToken(ctx context.Context, jwt string) (string, time.Time, error) {
	endpoint := fmt.Sprintf("%s/app/installations/%s/access_tokens", githubAPIBaseURL, p.installationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader([]byte("{}")))
	if err != nil {
		return "", time.Time{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	res, err := p.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("github app installation token request failed: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", time.Time{}, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", time.Time{}, fmt.Errorf("github app installation token request failed: status=%d body=%s", res.StatusCode, string(body))
	}

	var response struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", time.Time{}, fmt.Errorf("github app installation token response decode failed: %w", err)
	}
	if response.Token == "" {
		return "", time.Time{}, errors.New("github app installation token response did not include token")
	}

	return response.Token, response.ExpiresAt, nil
}

func encodeSegment(value interface{}) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("github app private key must be PEM encoded")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("github app private key parse failed: %w", err)
	}

	key, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("github app private key must be RSA private key")
	}

	return key, nil
}
