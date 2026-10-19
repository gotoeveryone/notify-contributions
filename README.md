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

ローカル実行では、初回の認可コード交換で得た `refresh_token` と OAuth アプリの ID、Secret、認可時と同じ Redirect URI を `.env` に設定してください。初回実行で更新トークンを交換し、新しい値を `.gitlab-oauth-refresh-token` に保存します。次回からはこのファイルの値を使用します。ファイルは Git の管理対象外です。`DEBUG=1` を指定すると通知は送らず、取得結果を表示します。

GitHub App の private key は `\n` を含む1行の値として `GITHUB_APP_PRIVATE_KEY` に設定してください。

`APP_SECRET_ARN` を設定した場合は、Secrets Manager から secret JSON を取得し、GitLab OAuth の更新トークンを交換・保存します。未設定の場合は `.env` とローカルのトークンファイルを使用します。ローカルのトークンファイルを消した後は、失効済みの `.env` の更新トークンでは再実行できないため、再認可が必要です。

## Run with AWS SAM (scheduled Lambda)

```console
$ sam build
$ sam deploy --guided
```

`template.yaml` では GitHub Actions と同じく毎日 `15:00 UTC` に実行するよう設定しています。

`template.yaml` は Secrets Manager を前提にしています。Lambda の環境変数には `APP_SECRET_ARN` のみを渡し、secret 値は実行時に取得します。

GitHub App 認証を使う場合は、以下のキーを持つ Secret を作成してください。

```json
{
  "GITHUB_APP_ID": "12345",
  "GITHUB_APP_INSTALLATION_ID": "67890",
  "GITHUB_APP_PRIVATE_KEY": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
  "GITLAB_OAUTH_CLIENT_ID": "GitLab OAuth application ID",
  "GITLAB_OAUTH_CLIENT_SECRET": "GitLab OAuth application secret",
  "GITLAB_OAUTH_REDIRECT_URI": "https://example.com/callback",
  "GITLAB_OAUTH_REFRESH_TOKEN": "initial refresh token",
  "SLACK_WEBHOOK_URL": "xxx"
}
```

GitLab OAuth アプリを `read_user` スコープで作成し、認可コードフローで取得した初回の更新トークンを Secret に登録してください。`GITLAB_OAUTH_REDIRECT_URI` は認可時と同じ値にします。実行ごとに更新トークンを交換し、Secrets Manager の同じ Secret に新しい値を保存します。Lambda には `GetSecretValue` と `PutSecretValue` が必要です。更新トークンは単回利用のため、この Lambda を同時実行しないでください。トークン交換後に Secret の保存が失敗した場合は、再認可して新しい更新トークンを登録する必要があります。[GitLab OAuth API](https://docs.gitlab.com/api/oauth2/)

PAT (`GITHUB_TOKEN`) は互換用の fallback として利用できます。ローカルなどで PAT を使う場合は、GitHub App の環境変数、または Secrets Manager の GitHub App 関連キーを未設定にしてください。

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
