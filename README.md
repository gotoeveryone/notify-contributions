# Notify Github Contributions

![Build Status](https://github.com/gotoeveryone/notify-github-contributions/workflows/Build/badge.svg)

## Requirements

- Golang
- AWS account (use to Lambda and Secrets Manager)
- Twitter account

## Setup

```console
$ go mod download
```

## Run (Local)

```console
$ USER_NAME={user_name} DEBUG=1 go run src/cmd/main.go
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

## Deploy

```console
$ cp .chalice/config.json.example .chalice/config.json # Please edit the value.
$ pipenv run deploy
```
