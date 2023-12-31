package cache

import (
	"context"
	"ordermngmt/internal/usecase"
	model "ordermngmt/pkg/entity"
)

type Cache interface {
	Get(ctx context.Context, id string) (model.Order, error)
	Add(ctx context.Context, order model.Order) error
}

type Service struct {
	repo Cache
}

var _ usecase.Cache = &Service{}

func New(repo Cache) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetOrder(ctx context.Context, id model.OrderID) (model.Order, error) {
	order, err := s.repo.Get(ctx, id.ID)
	return order, err
}

func (s *Service) AddOrder(ctx context.Context, order model.Order) error {

	return s.repo.Add(ctx, order)
}
