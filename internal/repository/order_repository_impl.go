package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/mzulfanw/order-service/internal/domain"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func exponentialBackoff(attempt int) time.Duration {
	base := 100 * time.Millisecond
	maxBackoff := 2 * time.Second
	d := time.Duration(float64(base) * math.Pow(2, float64(attempt)))
	if d > maxBackoff {
		return maxBackoff
	}
	return d
}

func (r orderRepository) Create(ctx context.Context, o *domain.Order) error {
	var lastErr error
	maxAttempts := 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			lastErr = err
			time.Sleep(exponentialBackoff(attempt))
			continue
		}

		query := `INSERT INTO orders (product_id, total_price, status, created_at)
		          VALUES ($1, $2, $3, $4) RETURNING id`
		err = tx.QueryRowContext(ctx, query, o.ProductID, o.TotalPrice, o.Status, time.Now()).Scan(&o.ID)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("rollback error: %v\n", rbErr)
			}
			lastErr = err
			time.Sleep(exponentialBackoff(attempt))
			continue
		}
		if err := tx.Commit(); err != nil {
			lastErr = err
			time.Sleep(exponentialBackoff(attempt))
			continue
		}
		return nil
	}
	return errors.New("order creation failed after retries: " + lastErr.Error())
}

func (r orderRepository) GetByProductID(ctx context.Context, productID string) ([]domain.Order, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, product_id, total_price, status, created_at 
	FROM orders WHERE product_id=$1`, productID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("rows close error: %v\n", err)
		}
	}()

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
