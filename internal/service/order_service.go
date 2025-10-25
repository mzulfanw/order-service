package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mzulfanw/order-service/internal/domain"
	"github.com/mzulfanw/order-service/internal/repository"
	"github.com/mzulfanw/order-service/pkg"
)

type OrderService struct {
	repo       repository.OrderRepository
	cache      pkg.Cache
	client     pkg.ProductClient
	messageBus pkg.RabbitMQPublisher
}

func NewOrderService(repo repository.OrderRepository, cache pkg.Cache, client pkg.ProductClient, messageBus pkg.RabbitMQPublisher) *OrderService {
	return &OrderService{
		repo:       repo,
		cache:      cache,
		client:     client,
		messageBus: messageBus,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, productID string, quantity int) (*domain.Order, error) {
	product, err := s.client.GetProduct(productID)

	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product.Qty < quantity {
		return nil, fmt.Errorf("insufficient product quantity")
	}

	order := &domain.Order{
		ProductID:  productID,
		TotalPrice: product.Price * float64(quantity),
		Status:     "CREATED",
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	s.messageBus.PublishEvent("order.created", order)
	return order, nil
}

func (s *OrderService) GetByProductID(ctx context.Context, productID string) ([]domain.Order, error) {
	if cached, found := s.cache.Get(fmt.Sprintf("orders:%d", productID)); found {
		var orders []domain.Order
		json.Unmarshal(cached, &orders)
		return orders, nil
	}
	orders, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	data, _ := json.Marshal(orders)
	s.cache.Set(fmt.Sprintf("orders:%d", productID), data)
	return orders, nil
}
