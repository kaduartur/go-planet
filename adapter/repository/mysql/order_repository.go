package mysql

import (
	"context"
	"math/rand"

	"github.com/kaduartur/go-planet/pkg/log"
	"gorm.io/gorm"

	"github.com/kaduartur/go-planet/core/domain"
	"github.com/kaduartur/go-planet/core/port"
)

type orderRepository struct {
	log log.Logger
	db  *gorm.DB
}

func NewOrderRepo(log log.Logger, db *gorm.DB) port.OrderRepository {
	return &orderRepository{
		log: log,
		db:  db,
	}
}

func (r *orderRepository) Create(ctx context.Context, order domain.Order) error {
	products := make(Products, len(order.Products))
	for i, p := range order.Products {
		products[i] = Product{
			ID:            rand.Int(),
			OrderID:       order.ID.Id(),
			SkuID:         p.SkuID,
			Name:          p.Name,
			TotalAmount:   p.Price.Amount,
			TotalScale:    p.Price.Scale,
			TotalCurrency: p.Price.Currency,
		}
	}

	o := Order{
		ID:               order.ID.Id(),
		StoreID:          order.StoreID,
		ConsumerID:       order.ConsumerID,
		Products:         products,
		ShippingID:       order.Shipping.ID,
		ShippingAmount:   order.Shipping.Price.Amount,
		ShippingScale:    order.Shipping.Price.Scale,
		ShippingCurrency: order.Shipping.Price.Currency,
		SubTotalAmount:   order.SubTotal.Amount,
		SubTotalScale:    order.SubTotal.Scale,
		SubTotalCurrency: order.SubTotal.Currency,
		TotalAmount:      order.Total.Amount,
		TotalScale:       order.Total.Scale,
		TotalCurrency:    order.Total.Currency,
	}

	result := r.db.WithContext(ctx).
		Table("orders").
		Create(&o)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
