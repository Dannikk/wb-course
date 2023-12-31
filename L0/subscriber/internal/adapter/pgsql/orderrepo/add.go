package orderrepo

import (
	"context"
	"fmt"
	"log"
	entity "ordermngmt/pkg/entity"
	"strings"
)

func (repo *Repository) AddOrder(ctx context.Context, order entity.Order) error {
	log.Printf("adding order: [%v]\n", order.Order_uid)
	query := `
		INSERT INTO orders (
			order_uid,
			track_number,
			entry,
			locale,
			internal_signature,
			customer_id,
			delivery_service,
			shardkey,
			sm_id,
			date_created,
			oof_shard
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`
	_, err := repo.pgsql.ExecContext(ctx, query,
		order.Order_uid,
		order.Track_number,
		order.Entry,
		order.Locale,
		order.Internal_signature,
		order.Customer_id,
		order.Delivery_service,
		order.Shardkey,
		order.Sm_id,
		order.Date_created,
		order.Oof_shard,
	)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO delivery (
			order_uid,
			name,
			phone,
			zip,
			city,
			address,
			region,
			email
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err = repo.pgsql.ExecContext(ctx, query,
		order.Order_uid,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO payment (
			order_uid,
			request_id,
			currency,
			provider,
			amount,
			payment_dt,
			bank,
			delivery_cost,
			goods_total,
			custom_fee
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	_, err = repo.pgsql.ExecContext(ctx, query,
		order.Order_uid,
		order.Payment.Request_id,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.Payment_dt,
		order.Payment.Bank,
		order.Payment.Delivery_cost,
		order.Payment.Goods_total,
		order.Payment.Custom_fee,
	)
	if err != nil {
		return err
	}

	query = `
	INSERT INTO items (
		order_uid,
		chrt_id,
		track_number, 
		price,        
		rid,          
		name,         
		sale,         
		size,         
		total_price,  
		nm_id,        
		brand,        
		status     
	)
	VALUES `
	args := make([]interface{}, 0, len(order.Items)*12)
	for i, item := range order.Items {
		query += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			i*12+1, i*12+2, i*12+3, i*12+4, i*12+5, i*12+6, i*12+7, i*12+8, i*12+9, i*12+10, i*12+11, i*12+12)
		args = append(args,
			order.Order_uid,
			item.Chrt_id,
			item.Track_number,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.Total_price,
			item.Nm_id,
			item.Brand,
			item.Status)
	}
	query = strings.TrimSuffix(query, ",") + ";"

	_, err = repo.pgsql.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert items: %w", err)
	}

	log.Printf("Inserted order [%v] that consists of %v items\n", order.Order_uid, len(order.Items))

	return nil
}
