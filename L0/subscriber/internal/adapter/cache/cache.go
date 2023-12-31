package cache_repo

import (
	"context"
	"ordermngmt/internal/custom_error"
	"ordermngmt/internal/service/cache"
	model "ordermngmt/pkg/entity"
	"sync"
)

type Repository struct {
	repo sync.Map
}

func New() *Repository {
	return &Repository{
		repo: sync.Map{},
	}
}
func (r *Repository) Add(ctx context.Context, order model.Order) error {
	r.repo.Store(order.Order_uid, order)
	return nil
}

// Get implements cache.Cache.
func (r *Repository) Get(ctx context.Context, id string) (model.Order, error) {
	order, ok := r.repo.Load(id)
	if !ok {
		return model.Order{}, custom_error.ErrNotFoundCache
	}

	if order, isOrder := order.(model.Order); isOrder {
		return order, nil
	}

	return model.Order{}, custom_error.ErrNonOrderType

}

var _ cache.Cache = &Repository{}
