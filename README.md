# Notify Contributions

![CI Status](https://github.com/gotoeveryone/notify-contributions/workflows/CI/badge.svg)

## Requirements

- Golang 1.23
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

## Run with AWS SAM (scheduled Lambda)

```console
$ sam build
$ sam deploy --guided
```

`template.yaml` では GitHub Actions と同じく毎日 `15:00 UTC` に実行するよう設定しています。

`template.yaml` は Secrets Manager を前提にしています。少なくとも以下のキーを持つ Secret を作成してください。

```json
{
  "GITHUB_TOKEN": "xxx",
  "GITLAB_TOKEN": "xxx",
  "SLACK_WEBHOOK_URL": "xxx"
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
