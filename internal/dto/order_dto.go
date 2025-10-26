package dto

type CreateOrderDTO struct {
	ProductID string `json:"productId" validate:"required,uuid4"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}
