package usecase

import (
	"errors"
	"testing"
	"time"

	"gotoeveryone/notify-contributions/src/domain/client"
	"gotoeveryone/notify-contributions/src/domain/entity"
	"gotoeveryone/notify-contributions/src/mock"
)

func TestContributionExec(t *testing.T) {
	u := contributionUsecase{
		sources: []client.Contribution{&mock.ContributionClient{Entity: &entity.Contribution{
			BaseDate:       time.Now(),
			YesterdayCount: 1,
			BaseDateCount:  2,
		}}},
		notification: &mock.NotificationClient{},
	}

	err := u.Exec(time.Now())
	if err != nil {
		t.Errorf("Failed: Error is not nil, actual: [%s]", err.Error())
	}
}

func TestContributionExecContributionNotFound(t *testing.T) {
	u := contributionUsecase{
		sources:      []client.Contribution{&mock.ContributionClient{Entity: nil}},
		notification: &mock.NotificationClient{},
	}

	err := u.Exec(time.Now())
	if !errors.Is(err, ErrContributionNotFound) {
		t.Errorf("Failed: Error is not matched, actual: [%s], expected: [%s]", err.Error(), ErrContributionNotFound.Error())
	}
}
