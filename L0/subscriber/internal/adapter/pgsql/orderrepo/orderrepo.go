package orderrepo

import (
	"database/sql"
	orderstorage "ordermngmt/internal/service/orderstorage"
)

type Repository struct {
	pgsql *sql.DB
}

var _ orderstorage.Repository = &Repository{}

func New(pgsql_connect *sql.DB) *Repository {
	return &Repository{
		pgsql: pgsql_connect,
	}
}
