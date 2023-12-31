package orderrepo

import (
	"context"
	"log"
	entity "ordermngmt/pkg/entity"
)

func (repo *Repository) GetOrder(ctx context.Context, id entity.OrderID) (entity.Order, error) {
	order := entity.Order{}

	err := repo.pgsql.QueryRowContext(ctx,
		`SELECT 
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
		FROM orders WHERE order_uid = $1`, id.ID).Scan(
		&order.Order_uid,
		&order.Track_number,
		&order.Entry,
		&order.Locale,
		&order.Internal_signature,
		&order.Customer_id,
		&order.Delivery_service,
		&order.Shardkey,
		&order.Sm_id,
		&order.Date_created,
		&order.Oof_shard,
	)
	if err != nil {
		log.Println("Error querying orders table:", err)
		return entity.Order{}, err
	}

	err = repo.pgsql.QueryRowContext(ctx,
		`SELECT 
			name,
			phone,
			zip,
			city,
			address,
			region,
			email
		FROM delivery WHERE order_uid = $1`, id.ID).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		log.Println("Error querying delivery table:", err)
		return entity.Order{}, err
	}

	err = repo.pgsql.QueryRowContext(ctx,
		`SELECT 
			request_id,
			currency,
			provider,
			amount,
			payment_dt,
			bank,
			delivery_cost,
			goods_total,
			custom_fee
		FROM payment WHERE order_uid = $1`, id.ID).Scan(
		&order.Payment.Request_id,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.Payment_dt,
		&order.Payment.Bank,
		&order.Payment.Delivery_cost,
		&order.Payment.Goods_total,
		&order.Payment.Custom_fee,
	)
	if err != nil {
		log.Println("Error querying payment table:", err)
		return entity.Order{}, err
	}

	rows, err := repo.pgsql.QueryContext(ctx,
		`SELECT 
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
		FROM items WHERE order_uid = $1`, id.ID)
	if err != nil {
		log.Println("Error querying items table:", err)
		return entity.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		item := entity.Item{}
		err = rows.Scan(
			&item.Chrt_id,
			&item.Track_number,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.Total_price,
			&item.Nm_id,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return entity.Order{}, err
		}
		order.Items = append(order.Items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Println("Error iterating over rows:", err)
		return entity.Order{}, err
	}

	return order, nil
}
