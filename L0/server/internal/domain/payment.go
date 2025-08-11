package domain

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

func (p *Payment) Validate() error {
	// TODO: validation heavily depends on business logic, which I don't have
	return nil
}
