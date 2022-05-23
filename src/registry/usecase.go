package registry

import (
	"gotoeveryone/notify-github-contributions/src/domain/client"
	"gotoeveryone/notify-github-contributions/src/usecase"
)

// NewContributionUsecase is create new usecase instance for contribution
func NewContributionUsecase(source client.Contribution, notification client.Notification) usecase.Contribution {
	return usecase.NewContributionUsecase(source, notification)
}
