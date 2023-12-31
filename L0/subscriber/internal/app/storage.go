package app

import (
	"database/sql"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
	"log"
	"ordermngmt/internal/config"
)

func newPostgresqlConnection(cfg config.PostgresCfg) (*sql.DB, error) {
	pgConnString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.Hostname,
		cfg.DB,
		cfg.User,
		cfg.Password,
	)

	// Open the connection
	db, err := sql.Open("postgres", pgConnString)

	if err != nil {
		log.Printf("error opening connection: %v\n", err)
		return nil, err
	}

	var attemptnum int
	// check the connection
	err = backoff.Retry(
		func() error {
			err := db.Ping()
			attemptnum += 1
			log.Printf("Failed to ping db. Attempt=%v: %v\n", attemptnum, err)
			return err
		},
		backoff.NewExponentialBackOff(),
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initDB(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		order_uid          VARCHAR PRIMARY KEY,
		track_number       VARCHAR,
		entry              VARCHAR,
		locale             VARCHAR,
		internal_signature VARCHAR,
		customer_id        VARCHAR,
		delivery_service   VARCHAR,
		shardkey           VARCHAR,
		sm_id              INT NOT NULL,
		date_created       DATE,
		oof_shard          VARCHAR
	);`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create orders table: %v", err)
	}

	query = `
	CREATE TABLE IF NOT EXISTS delivery (
		order_uid 	VARCHAR PRIMARY KEY,
		name    	VARCHAR,
		phone   	VARCHAR,
		zip     	VARCHAR,
		city    	VARCHAR,
		address 	VARCHAR,
		region  	VARCHAR,
		email   	VARCHAR,
		FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
	);`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create delivery table: %v", err)
	}

	query = `
	CREATE TABLE IF NOT EXISTS payment (
		order_uid      	VARCHAR PRIMARY KEY,
		transaction     VARCHAR,
		request_id      VARCHAR,
		currency        VARCHAR,
		provider        VARCHAR,
		amount          BIGINT,
		payment_dt      BIGINT,
		bank            VARCHAR,
		delivery_cost   BIGINT,
		goods_total     BIGINT,
		custom_fee      BIGINT,
		FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
	);`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create payment table: %v", err)
	}

	query = `
	CREATE TABLE IF NOT EXISTS items (
		order_uid    VARCHAR,
		chrt_id      BIGINT,
		track_number VARCHAR,
		price        BIGINT,
		rid          VARCHAR,
		name         VARCHAR,
		sale         BIGINT,
		size         VARCHAR,
		total_price  BIGINT,
		nm_id        BIGINT,
		brand        VARCHAR,
		status       SMALLINT,
		FOREIGN KEY (order_uid) REFERENCES orders(order_uid)
	);`
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create items table: %v", err)
	}

	return nil
}

func CloseConnection(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db pointer is nil")
	}

	log.Println("Closing connection...")
	if err := db.Close(); err != nil {
		return err
	}
	log.Println("Connection closed!")
	return nil
}
