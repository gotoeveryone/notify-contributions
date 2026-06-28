package app

import "testing"

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
