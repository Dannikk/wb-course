package model

type OrderID struct {
	ID string
}

type Order struct {
	Order_uid          string   `json:"order_uid" validate:"nonzero" db:"order_uid" binding:"required"`
	Track_number       string   `json:"track_number"`
	Entry              string   `json:"entry"`
	Delivery           Delivery `json:"delivery"`
	Payment            Payment  `json:"payment"`
	Items              []Item   `json:"items"`
	Locale             string   `json:"locale"`
	Internal_signature string   `json:"internal_signature"`
	Customer_id        string   `json:"customer_id"`
	Delivery_service   string   `json:"delivery_service"`
	Shardkey           string   `json:"shardkey"`
	Sm_id              int64    `json:"sm_id"`
	Date_created       string   `json:"date_created"`
	Oof_shard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction   string `json:"transaction"`
	Request_id    string `json:"request_id"`
	Currency      string `json:"currency"`
	Provider      string `json:"provider"`
	Amount        uint64 `json:"amount"`
	Payment_dt    uint64 `json:"payment_id"`
	Bank          string `json:"bank"`
	Delivery_cost uint64 `json:"delivery_cost"`
	Goods_total   uint64 `json:"goods_total"`
	Custom_fee    uint64 `json:"custom_fee"`
}

type Item struct {
	Chrt_id      uint64 `json:"chrt_id"`
	Track_number string `json:"track_number"`
	Price        uint64 `json:"price"`
	Rid          string `json:"rid"`
	Name         string `json:"name"`
	Sale         uint64 `json:"sale"`
	Size         string `json:"size"`
	Total_price  uint64 `json:"total_price"`
	Nm_id        uint64 `json:"nm_id"`
	Brand        string `json:"brand"`
	Status       uint16 `json:"status"`
}
