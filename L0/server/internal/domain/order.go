package domain

import (
	"context"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid" faker:"uuid_hyphenated"`
	TrackNumber       string    `json:"track_number" faker:"len=15"`
	Entry             string    `json:"entry" faker:"len=4"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items" faker:"len=1,5"` // Generates 1-5 items
	Locale            string    `json:"locale" faker:"oneof: en, ru, de, fr, es, it"`
	InternalSignature string    `json:"internal_signature" faker:"word"`
	CustomerID        string    `json:"customer_id" faker:"uuid_digit"`
	DeliveryService   string    `json:"delivery_service" faker:"oneof: meest, fedex, ups, dhl, usps"`
	Shardkey          string    `json:"shardkey" faker:"oneof: 1, 2, 3, 4, 5, 6, 7, 8, 9"`
	SMID              int       `json:"sm_id" faker:"boundary_start=1, boundary_end=100"`
	DateCreated       time.Time `json:"date_created"`
	OOFShard          string    `json:"oof_shard" faker:"oneof: 1, 2, 3"`
}

type IOrderRepository interface {
	Save(ctx context.Context, order *Order) error
	Get(ctx context.Context, orderID string) (*Order, error)
}

type IOrderCache interface {
	Save(ctx context.Context, order *Order) error
	Get(ctx context.Context, orderID string) (*Order, error)
}
