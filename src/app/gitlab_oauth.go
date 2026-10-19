package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type gitLabTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const localGitLabRefreshTokenFile = ".gitlab-oauth-refresh-token"

// refreshGitLabToken rotates the refresh token and saves it before using the access token.
func refreshGitLabToken(cfg *Config) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if cfg.AppSecretARN == "" {
		refreshToken, err := readLocalGitLabRefreshToken(localGitLabRefreshTokenFile, cfg.GitLabOAuthRefreshToken)
		if err != nil {
			return "", fmt.Errorf("read local GitLab OAuth refresh token: %w", err)
		}
		cfg.GitLabOAuthRefreshToken = refreshToken
	}
	tokens, err := exchangeGitLabToken(ctx, cfg, &http.Client{Timeout: 10 * time.Second}, "https://gitlab.com/oauth/token")
	if err != nil {
		return "", err
	}
	if cfg.AppSecretARN == "" {
		err = writeLocalGitLabRefreshToken(localGitLabRefreshTokenFile, tokens.RefreshToken)
	} else {
		err = saveGitLabRefreshToken(ctx, cfg.AppSecretARN, tokens.RefreshToken)
	}
	if err != nil {
		return "", fmt.Errorf("save GitLab OAuth refresh token: %w", err)
	}
	return tokens.AccessToken, nil
}

func readLocalGitLabRefreshToken(path, initialToken string) (string, error) {
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return initialToken, nil
	}
	if err != nil {
		return "", err
	}
	token := strings.TrimSpace(string(b))
	if token == "" {
		return "", fmt.Errorf("local GitLab OAuth refresh token file is empty")
	}
	return token, nil
}

func writeLocalGitLabRefreshToken(path, token string) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".gitlab-oauth-refresh-token-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(token + "\n"); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

func exchangeGitLabToken(ctx context.Context, cfg *Config, client *http.Client, endpoint string) (*gitLabTokenResponse, error) {
	values := url.Values{
		"client_id":     {cfg.GitLabOAuthClientID},
		"client_secret": {cfg.GitLabOAuthClientSecret},
		"refresh_token": {cfg.GitLabOAuthRefreshToken},
		"grant_type":    {"refresh_token"},
		"redirect_uri":  {cfg.GitLabOAuthRedirectURI},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build GitLab OAuth request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh GitLab OAuth token: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		// Only expose the OAuth error code; the response body can contain secrets.
		var oauthError struct {
			Code string `json:"error"`
		}
		if json.NewDecoder(io.LimitReader(res.Body, 1<<12)).Decode(&oauthError) == nil {
			switch oauthError.Code {
			case "invalid_grant", "invalid_client", "invalid_request", "unauthorized_client", "unsupported_grant_type", "invalid_scope":
				return nil, fmt.Errorf("refresh GitLab OAuth token: HTTP %d (%s)", res.StatusCode, oauthError.Code)
			}
		}
		return nil, fmt.Errorf("refresh GitLab OAuth token: HTTP %d", res.StatusCode)
	}
	var tokens gitLabTokenResponse
	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("decode GitLab OAuth token response: %w", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		return nil, fmt.Errorf("GitLab OAuth response is missing a token")
	}

	return &tokens, nil
}

func saveGitLabRefreshToken(ctx context.Context, secretID, refreshToken string) error {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	client := secretsmanager.NewFromConfig(awsCfg)
	secret, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{SecretId: &secretID})
	if err != nil {
		return err
	}
	if secret.SecretString == nil {
		return fmt.Errorf("secret string is empty")
	}
	values := map[string]string{}
	if err := json.Unmarshal([]byte(*secret.SecretString), &values); err != nil {
		return err
	}
	values["GITLAB_OAUTH_REFRESH_TOKEN"] = refreshToken
	updated, err := json.Marshal(values)
	if err != nil {
		return err
	}
	value := string(updated)
	_, err = client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId: &secretID, SecretString: &value,
	})
	return err
}
