package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mzulfanw/order-service/internal/domain"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r orderRepository) Create(ctx context.Context, o *domain.Order) error {
	fmt.Println("here")
	query := `INSERT INTO orders (product_id, total_price, status, created_at)
	          VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRowContext(ctx, query, o.ProductID, o.TotalPrice, o.Status, time.Now()).
		Scan(&o.ID)
}

func (r orderRepository) GetByProductID(ctx context.Context, productID string) ([]domain.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, product_id, total_price, status, created_at 
	FROM orders WHERE product_id=$1`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.TotalPrice, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
