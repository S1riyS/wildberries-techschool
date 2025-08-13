package domain

import (
	"context"
	"errors"
)

type Payment struct {
	Transaction  string `json:"transaction" faker:"uuid_hyphenated"`
	RequestID    string `json:"request_id" faker:"uuid_digit"`
	Currency     string `json:"currency" faker:"currency"`
	Provider     string `json:"provider" faker:"oneof: wbpay, paypal, stripe, applepay, googlepay"`
	Amount       int    `json:"amount" faker:"boundary_start=1000, boundary_end=100000"`
	PaymentDT    int64  `json:"payment_dt" faker:"unix_time"`
	Bank         string `json:"bank" faker:"oneof: alpha, sber, tinkoff, vtb, gazprom, raiffeisen"`
	DeliveryCost int    `json:"delivery_cost" faker:"boundary_start=500, boundary_end=5000"`
	GoodsTotal   int    `json:"goods_total" faker:"boundary_start=300, boundary_end=20000"`
	CustomFee    int    `json:"custom_fee" faker:"boundary_start=0, boundary_end=500"`
}

type IPaymentRepository interface {
	Save(ctx context.Context, payment *Payment) error
	Get(ctx context.Context, transaction string) (*Payment, error)
}

func (p *Payment) Validate() error {
	// TODO: validation heavily depends on business logic, which I don't have

	if p.Transaction == "" {
		return errors.New("transaction cannot be empty")
	}
	if p.RequestID == "" {
		return errors.New("request_id cannot be empty")
	}
	if p.Currency == "" {
		return errors.New("currency cannot be empty")
	}
	if p.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	// etc

	return nil
}
