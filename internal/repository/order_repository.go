package repository

import (
	"context"

	"github.com/mzulfanw/order-service/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, o *domain.Order) error
	GetByProductID(ctx context.Context, productID string) ([]domain.Order, error)
}
