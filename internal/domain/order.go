package domain

import "time"

type Order struct {
	ID         string    `json:"id"`
	ProductID  string    `json:"productId"`
	TotalPrice float64   `json:"totalPrice"`
	Quantity   int       `json:"quantity"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}
