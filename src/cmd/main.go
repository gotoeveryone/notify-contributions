package main

import (
	"context"
	"gotoeveryone/notify-github-contributions/src/registry"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/sirupsen/logrus"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	userName := os.Getenv("USER_NAME")
	baseDate := time.Now().UTC()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "failed", err
	}
	ssmClient := ssm.NewFromConfig(cfg)

	cc := registry.NewGitHubClient()
	tc, err := registry.NewTwitterClient(*ssmClient)
	if err != nil {
		if err != nil {
			return "failed", err
		}
	}

	u := registry.NewContributionUsecase(cc, tc)
	if err := u.Exec(userName, baseDate); err != nil {
		return "failed", err
	}
	return "success", err
}

func main() {
	if os.Getenv("DEBUG") == "1" {
		logrus.SetFormatter(&logrus.JSONFormatter{})

		res, err := HandleRequest(context.TODO(), MyEvent{Name: "debug"})
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		logrus.Info(res)
		return
	}

	lambda.Start(HandleRequest)
}
