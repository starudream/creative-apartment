package iroute

import (
	"net/http"

	"github.com/starudream/creative-apartment/config"
)

type ListCustomersResp struct {
	Customers []Customer `json:"customers"`
}

type Customer struct {
	Phone string `json:"phone"`
}

func ListCustomers(c *Context) {
	resp := ListCustomersResp{}

	for _, customer := range config.GetCustomers() {
		resp.Customers = append(resp.Customers, Customer{Phone: customer.Phone})
	}

	c.JSON(http.StatusOK, resp)
}
