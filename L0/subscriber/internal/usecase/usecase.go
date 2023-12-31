package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"ordermngmt/internal/custom_error"
	handler "ordermngmt/internal/handler/http/api"
	entity "ordermngmt/pkg/entity"
)

type OrderStorage interface {
	GetOrder(ctx context.Context, id entity.OrderID) (entity.Order, error)
	AddOrder(ctx context.Context, order entity.Order) error
}

type Cache interface {
	GetOrder(ctx context.Context, id entity.OrderID) (entity.Order, error)
	AddOrder(ctx context.Context, order entity.Order) error
}

type UseCase struct {
	storage OrderStorage
	cache   Cache
}

var _ handler.UseCase = &UseCase{}

func New(storage OrderStorage, cache Cache) *UseCase {
	return &UseCase{
		storage,
		cache,
	}
}

func (uc *UseCase) GetOrder(ctx context.Context, id entity.OrderID) (entity.Order, error) {
	order, err := uc.cache.GetOrder(ctx, id)

	if err == nil {
		return order, nil
	} else if err != custom_error.ErrNotFoundCache {
		return order, err
	}

	order, err = uc.storage.GetOrder(ctx, id)

	if err != nil {
		return order, err
	}

	err = uc.cache.AddOrder(ctx, order)

	return order, err
}

func (uc *UseCase) AddOrder(ctx context.Context, data []byte) error {
	order := entity.Order{}
	if err := json.Unmarshal(data, &order); err != nil {
		return fmt.Errorf("while unmarshal order: %w", err)
	}

	return uc.storage.AddOrder(ctx, order)
}
