package app

import "testing"

func TestLoadConfigFromEnv(t *testing.T) {
	cfg, err := loadConfig(&configSource{
		getenv: mapGetenv(map[string]string{
			"GITHUB_APP_ID":              "12345",
			"GITHUB_APP_INSTALLATION_ID": "67890",
			"GITHUB_APP_PRIVATE_KEY":     "private_key",
			"GITHUB_USER_NAME":           "github-user",
			"GITLAB_USER_ID":             "12345678",
			"GITLAB_TOKEN":               "gitlab-token",
			"NOTIFY_TYPE":                "slack",
			"SLACK_WEBHOOK_URL":          "slack-webhook",
		}),
	})
	if err != nil {
		t.Fatal(err)
	}

	if cfg.GitHubAppID != "12345" {
		t.Errorf("GitHubAppID is not matched, actual: [%s], expected: [%s]", cfg.GitHubAppID, "12345")
	}
	if cfg.SlackWebhook != "slack-webhook" {
		t.Errorf("SlackWebhook is not matched, actual: [%s], expected: [%s]", cfg.SlackWebhook, "slack-webhook")
	}
}

func TestLoadConfigPrefersSecretValuesAndFallsBackToEnv(t *testing.T) {
	cfg, err := loadConfig(&configSource{
		getenv: mapGetenv(map[string]string{
			"GITHUB_APP_ID":              "env-github-app-id",
			"GITHUB_APP_INSTALLATION_ID": "env-installation-id",
			"GITHUB_APP_PRIVATE_KEY":     "env-private-key",
			"GITHUB_USER_NAME":           "github-user",
			"GITLAB_USER_ID":             "12345678",
			"GITLAB_TOKEN":               "env-gitlab-token",
			"NOTIFY_TYPE":                "slack",
			"SLACK_WEBHOOK_URL":          "env-slack-webhook",
		}),
		secrets: map[string]string{
			"GITHUB_APP_ID":              "secret-github-app-id",
			"GITHUB_APP_INSTALLATION_ID": "secret-installation-id",
			"GITHUB_APP_PRIVATE_KEY":     "secret-private-key",
			"GITLAB_TOKEN":               "secret-gitlab-token",
			"SLACK_WEBHOOK_URL":          "secret-slack-webhook",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if cfg.GitHubAppID != "secret-github-app-id" {
		t.Errorf("GitHubAppID is not matched, actual: [%s], expected: [%s]", cfg.GitHubAppID, "secret-github-app-id")
	}
	if cfg.GitHubUser != "github-user" {
		t.Errorf("GitHubUser is not matched, actual: [%s], expected: [%s]", cfg.GitHubUser, "github-user")
	}
	if cfg.GitLabToken != "secret-gitlab-token" {
		t.Errorf("GitLabToken is not matched, actual: [%s], expected: [%s]", cfg.GitLabToken, "secret-gitlab-token")
	}
	if cfg.GitLabUserID != "12345678" {
		t.Errorf("GitLabUserID is not matched, actual: [%s], expected: [%s]", cfg.GitLabUserID, "12345678")
	}
	if cfg.SlackWebhook != "secret-slack-webhook" {
		t.Errorf("SlackWebhook is not matched, actual: [%s], expected: [%s]", cfg.SlackWebhook, "secret-slack-webhook")
	}
}

func TestConfigRequireGitHubAuthWithToken(t *testing.T) {
	cfg := Config{GitHubToken: "token"}

	if err := cfg.requireGitHubAuth(); err != nil {
		t.Fatal(err)
	}
}

func TestConfigRequireGitHubAuthWithGitHubApp(t *testing.T) {
	cfg := Config{
		GitHubAppID:             "12345",
		GitHubAppInstallationID: "67890",
		GitHubAppPrivateKey:     "private_key",
	}

	if err := cfg.requireGitHubAuth(); err != nil {
		t.Fatal(err)
	}
}

func TestConfigRequireGitHubAuthDefaultsToGitHubApp(t *testing.T) {
	cfg := Config{}

	err := cfg.requireGitHubAuth()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "GITHUB_APP_ID is required" {
		t.Errorf("Error is not matched, actual: [%s], expected: [%s]", err.Error(), "GITHUB_APP_ID is required")
	}
}

func TestConfigRequireGitHubAuthRequiresGitHubAppPrivateKey(t *testing.T) {
	cfg := Config{
		GitHubAppID:             "12345",
		GitHubAppInstallationID: "67890",
	}

	if err := cfg.requireGitHubAuth(); err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func mapGetenv(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}
