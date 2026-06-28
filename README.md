# Notify Contributions

![CI Status](https://github.com/gotoeveryone/notify-contributions/workflows/CI/badge.svg)

## Requirements

- Golang 1.26.0
- Slack Account (Optional)
- Twitter Account (Optional)

## Setup

```console
$ go mod download
```

## Run

```console
$ cp .env.example .env # please edit values
$ DEBUG=1 go run src/cmd/main.go
```

GitHub App の private key は `\n` を含む1行の値として `GITHUB_APP_PRIVATE_KEY` に設定してください。

## Run with AWS SAM (scheduled Lambda)

```console
$ sam build
$ sam deploy --guided
```

`template.yaml` では GitHub Actions と同じく毎日 `15:00 UTC` に実行するよう設定しています。

`template.yaml` は Secrets Manager を前提にしています。GitHub App 認証を使う場合は、以下のキーを持つ Secret を作成してください。

```json
{
  "GITHUB_APP_ID": "12345",
  "GITHUB_APP_INSTALLATION_ID": "67890",
  "GITHUB_APP_PRIVATE_KEY": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
  "GITLAB_TOKEN": "xxx",
  "SLACK_WEBHOOK_URL": "xxx"
}
```

PAT (`GITHUB_TOKEN`) は互換用の fallback として利用できます。ローカルなどで PAT を使う場合は、GitHub App の環境変数を未設定にしてください。

```json
{
  "GITHUB_TOKEN": "xxx"
}
```

CI からは `sam deploy` 実行時に次のようにパラメータを渡します。

```console
$ sam deploy \
  --parameter-overrides \
  AppSecretArn=$APP_SECRET_ARN \
  GithubUserName=$GITHUB_USER_NAME \
  GitlabUserId=$GITLAB_USER_ID
```

## Code check and format

```console
$ # Code check
$ go vet ./...
$ # Format
$ go fmt ./...
```

## Test

```console
$ go test ./...
```
