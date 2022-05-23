package usecase

import (
	"errors"
	"testing"
	"time"

	"gotoeveryone/notify-github-contributions/src/domain/entity"
	"gotoeveryone/notify-github-contributions/src/mock"
)

func TestContributionExec(t *testing.T) {
	u := contributionUsecase{
		source: &mock.ContributionClient{Entity: &entity.Contribution{
			UserName:       "test",
			BaseDate:       time.Now(),
			YesterdayCount: 1,
			BaseDateCount:  2,
		}},
		notification: &mock.NotificationClient{},
	}

	err := u.Exec("test message", time.Now())
	if err != nil {
		t.Errorf("Failed: Error is not nil, actual: [%s]", err.Error())
	}
}

func TestContributionExecContributionNotFound(t *testing.T) {
	u := contributionUsecase{
		source:       &mock.ContributionClient{Entity: nil},
		notification: &mock.NotificationClient{},
	}

	err := u.Exec("test message", time.Now())
	if !errors.Is(err, ErrContributionNotFound) {
		t.Errorf("Failed: Error is not matched, actual: [%s], expected: [%s]", err.Error(), ErrContributionNotFound.Error())
	}
}
