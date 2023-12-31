package orders

import (
	"context"
	"ordermngmt/internal/usecase"
	entity "ordermngmt/pkg/entity"
)

type Repository interface {
	GetOrder(ctx context.Context, id entity.OrderID) (entity.Order, error)
	AddOrder(ctx context.Context, order entity.Order) error
}

type Service struct {
	repo Repository
}

var _ usecase.OrderStorage = &Service{}

func New(repo Repository) *Service {
	return &Service{
		repo: repo}
}

func (s *Service) GetOrder(ctx context.Context, id entity.OrderID) (entity.Order, error) {
	order, err := s.repo.GetOrder(ctx, id)

	return order, err
}

func (s *Service) AddOrder(ctx context.Context, order entity.Order) error {

	return s.repo.AddOrder(ctx, order)
}
