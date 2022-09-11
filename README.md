# Notify Github Contributions

![CI Status](https://github.com/gotoeveryone/notify-github-contributions/workflows/CI/badge.svg)

## Requirements

- Golang 1.19
- Twitter account

## Setup

```console
$ go mod download
```

## Run

```console
$ cp .env.example .env # please edit values
$ DEBUG=1 go run src/cmd/main.go
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
