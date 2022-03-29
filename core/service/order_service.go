package service

import (
	"context"

	"github.com/kaduartur/go-planet/core/port"
)

type orderService struct {
	repo port.OrderRepository
}

func NewOrderService(repo port.OrderRepository) port.OrderService {
	return &orderService{
		repo: repo,
	}
}

func (os *orderService) Create(ctx context.Context, createOrder port.CreateOrder) error {
	order := createOrder.ToDomain()
	if err := os.repo.Create(ctx, order); err != nil {
		return err
	}

	return nil
}
