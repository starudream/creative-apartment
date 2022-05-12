package config

import (
	"github.com/starudream/creative-apartment/internal/json"
)

type Customer struct {
	Phone string `json:"phone" yaml:"phone" validate:"required,min=1"`
	Token string `json:"token" yaml:"token" validate:"required,min=1"`
}

func (c *Customer) GetToken() string {
	if c == nil {
		return ""
	}
	return c.Token
}

var customers []*Customer

func SetCustomers(v any) {
	customers = json.ReMustUnmarshalTo[[]*Customer](v)
}

func GetCustomers() []*Customer {
	return customers
}

func GetCustomer(i int) *Customer {
	if i < 0 || i >= len(customers) {
		return nil
	}
	return customers[i]
}
