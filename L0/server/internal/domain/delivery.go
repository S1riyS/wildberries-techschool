package domain

type Delivery struct {
	Name    string `json:"name" faker:"name"`
	Phone   string `json:"phone" faker:"phone_number"`
	Zip     string `json:"zip" faker:"oneof: 10001, 90210, 60601, 75201, 94102"`
	City    string `json:"city" faker:"oneof: New York, Los Angeles, Chicago, Houston, Phoenix"`
	Address string `json:"address" faker:"sentence"`
	Region  string `json:"region" faker:"oneof: NY, CA, IL, TX, AZ"`
	Email   string `json:"email" faker:"email"`
}

func (d *Delivery) Validate() error {
	// TODO: validation heavily depends on business logic, which I don't have
	return nil
}
