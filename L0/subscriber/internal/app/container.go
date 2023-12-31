package app

import (
	"database/sql"
	"log"
	"ordermngmt/internal/config"
	"ordermngmt/internal/nats"
	"ordermngmt/internal/usecase"

	cacherepo "ordermngmt/internal/adapter/cache"
	orderrepo "ordermngmt/internal/adapter/pgsql/orderrepo"
	"ordermngmt/internal/service/cache"
	"ordermngmt/internal/service/orderstorage"
)

type Container struct {
	pgsql        *sql.DB
	cache        cache.Cache
	msgStream    chan []byte
	subscription *nats.Subscription
}

func NewContainer(cfg config.Config, pgsqlConnect *sql.DB) *Container {
	msgStream := make(chan []byte)
	sub, err := nats.NewSubscription(cfg.Stan, msgStream)
	if err != nil {
		log.Fatalf("can't subscribe: %v", err)
	}
	return &Container{
		pgsql:        pgsqlConnect,
		cache:        cacherepo.New(),
		msgStream:    msgStream,
		subscription: sub,
	}
}

func (c *Container) GetUseCase() *usecase.UseCase {
	return usecase.New(
		orders.New(orderrepo.New(c.pgsql)),
		cache.New(c.cache),
	)
}
