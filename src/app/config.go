package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"gotoeveryone/notify-contributions/src/domain/entity"
)

type Config struct {
	GitHubToken             string
	GitHubUser              string
	GitHubAppID             string
	GitHubAppInstallationID string
	GitHubAppPrivateKey     string
	GitLabUserID            string
	GitLabToken             string
	NotifyType              string
	SlackWebhook            string
	TwitterAuth             entity.TwitterAuth
}

func LoadConfig() (*Config, error) {
	source, err := newConfigSource(context.Background(), os.Getenv)
	if err != nil {
		return nil, err
	}
	return loadConfig(source)
}

type configSource struct {
	getenv  func(string) string
	secrets map[string]string
}

func newConfigSource(ctx context.Context, getenv func(string) string) (*configSource, error) {
	source := &configSource{getenv: getenv}

	appSecretArn := getenv("APP_SECRET_ARN")
	if appSecretArn == "" {
		return source, nil
	}

	secrets, err := loadSecretValues(ctx, appSecretArn)
	if err != nil {
		return nil, fmt.Errorf("load APP_SECRET_ARN: %w", err)
	}
	source.secrets = secrets

	return source, nil
}

func loadSecretValues(ctx context.Context, secretID string) (map[string]string, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)
	output, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretID,
	})
	if err != nil {
		return nil, err
	}
	if output.SecretString == nil {
		return nil, fmt.Errorf("secret string is empty")
	}

	values := map[string]string{}
	if err := json.Unmarshal([]byte(*output.SecretString), &values); err != nil {
		return nil, err
	}

	return values, nil
}

func (s *configSource) value(key string) string {
	if s.secrets != nil {
		if value, ok := s.secrets[key]; ok {
			return value
		}
	}
	return s.getenv(key)
}

func loadConfig(source *configSource) (*Config, error) {
	cfg := &Config{
		GitHubToken:             source.value("GITHUB_TOKEN"),
		GitHubUser:              source.value("GITHUB_USER_NAME"),
		GitHubAppID:             source.value("GITHUB_APP_ID"),
		GitHubAppInstallationID: source.value("GITHUB_APP_INSTALLATION_ID"),
		GitHubAppPrivateKey:     normalizePrivateKey(source.value("GITHUB_APP_PRIVATE_KEY")),
		GitLabUserID:            source.value("GITLAB_USER_ID"),
		GitLabToken:             source.value("GITLAB_TOKEN"),
		NotifyType:              source.value("NOTIFY_TYPE"),
		SlackWebhook:            source.value("SLACK_WEBHOOK_URL"),
		TwitterAuth: entity.TwitterAuth{
			ConsumerKey:       source.value("TWITTER_COMSUMER_KEY"),
			ConsumerSecret:    source.value("TWITTER_COMSUMER_SECRET"),
			AccessToken:       source.value("TWITTER_ACCESS_TOKEN"),
			AccessTokenSecret: source.value("TWITTER_ACCESS_TOKEN_SECRET"),
		},
	}

	if err := cfg.requireGitHubAuth(); err != nil {
		return nil, err
	}
	if err := require("GITHUB_USER_NAME", cfg.GitHubUser); err != nil {
		return nil, err
	}
	if err := require("GITLAB_USER_ID", cfg.GitLabUserID); err != nil {
		return nil, err
	}
	if err := require("GITLAB_TOKEN", cfg.GitLabToken); err != nil {
		return nil, err
	}
	if err := require("NOTIFY_TYPE", cfg.NotifyType); err != nil {
		return nil, err
	}

	switch cfg.NotifyType {
	case "slack":
		if err := require("SLACK_WEBHOOK_URL", cfg.SlackWebhook); err != nil {
			return nil, err
		}
	case "twitter":
		if err := require("TWITTER_COMSUMER_KEY", cfg.TwitterAuth.ConsumerKey); err != nil {
			return nil, err
		}
		if err := require("TWITTER_COMSUMER_SECRET", cfg.TwitterAuth.ConsumerSecret); err != nil {
			return nil, err
		}
		if err := require("TWITTER_ACCESS_TOKEN", cfg.TwitterAuth.AccessToken); err != nil {
			return nil, err
		}
		if err := require("TWITTER_ACCESS_TOKEN_SECRET", cfg.TwitterAuth.AccessTokenSecret); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid NOTIFY_TYPE: %s", cfg.NotifyType)
	}

	return cfg, nil
}

func (c *Config) UseGitHubApp() bool {
	return c.hasGitHubAppConfig() || c.GitHubToken == ""
}

func (c *Config) requireGitHubAuth() error {
	if !c.UseGitHubApp() {
		return require("GITHUB_TOKEN", c.GitHubToken)
	}

	if err := require("GITHUB_APP_ID", c.GitHubAppID); err != nil {
		return err
	}
	if err := require("GITHUB_APP_INSTALLATION_ID", c.GitHubAppInstallationID); err != nil {
		return err
	}
	if err := require("GITHUB_APP_PRIVATE_KEY", c.GitHubAppPrivateKey); err != nil {
		return err
	}

	return nil
}

func (c *Config) hasGitHubAppConfig() bool {
	return c.GitHubAppID != "" || c.GitHubAppInstallationID != "" || c.GitHubAppPrivateKey != ""
}

func require(key string, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", key)
	}
	return nil
}

func normalizePrivateKey(value string) string {
	return strings.ReplaceAll(value, `\n`, "\n")
}
