package application

import "context"

type Service interface {
	ListApplications(ctx context.Context, userId int64) ([]ApplicationResponse, error)
}

type service struct {
	appRepo Repository
}

func NewService(appRepo Repository) Service {
	s := &service{
		appRepo: appRepo,
	}
	return s
}

func (s *service) ListApplications(ctx context.Context, userId int64) ([]ApplicationResponse, error) {

}
