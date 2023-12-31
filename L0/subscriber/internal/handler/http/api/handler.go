package api

import (
	"context"
	entity "ordermngmt/pkg/entity"
)

//go:generate mockgen -source=handler.go -destination=mocks/ucmock.go

type UseCase interface {
	GetOrder(ctx context.Context, id entity.OrderID) (entity.Order, error)
}

type Handler struct {
	uc UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{uc: uc}
}
