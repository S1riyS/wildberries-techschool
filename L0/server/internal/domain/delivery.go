package domain

import (
	"context"
	"errors"
)

type Delivery struct {
	Name    string `json:"name" faker:"name"`
	Phone   string `json:"phone" faker:"phone_number"`
	Zip     string `json:"zip" faker:"oneof: 10001, 90210, 60601, 75201, 94102"`
	City    string `json:"city" faker:"oneof: New York, Los Angeles, Chicago, Houston, Phoenix"`
	Address string `json:"address" faker:"sentence"`
	Region  string `json:"region" faker:"oneof: NY, CA, IL, TX, AZ"`
	Email   string `json:"email" faker:"email"`
}

type IDeliveryRepository interface {
	Save(ctx context.Context, delivery *Delivery) (int, error)
	Get(ctx context.Context, id int) (*Delivery, error)
}

func (d *Delivery) Validate() error {
	// TODO: validation heavily depends on business logic, which I don't have

	if d.Name == "" {
		return errors.New("name cannot be empty")
	}
	if d.Phone == "" {
		return errors.New("phone cannot be empty")
	}
	if d.Zip == "" {
		return errors.New("zip cannot be empty")
	}
	if d.City == "" {
		return errors.New("city cannot be empty")
	}
	if d.Address == "" {
		return errors.New("address cannot be empty")
	}
	// etc

	return nil
}
