package port

import (
	"context"

	"github.com/kaduartur/go-planet/core/domain"
)

type Money struct {
	Amount   int
	Scale    int
	Currency string
}

type Product struct {
	ID       string
	SkuID    string
	Name     string
	Quantity int
	Price    Money
}

type Products []Product

func (pp Products) ToDomain() domain.Products {
	products := make(domain.Products, 0)
	for _, p := range pp {
		price := domain.NewMoney(p.Price.Amount, p.Price.Scale, p.Price.Currency)
		product := domain.NewProduct(p.ID, p.SkuID, p.Name, p.Quantity, price)
		products = append(products, product)
	}

	return products
}

type CreateOrder struct {
	StoreID    string
	ConsumerID string
	Products   Products
	Shipping   struct {
		ID    string
		Price Money
	}
	SubTotal Money
	Total    Money
}

func (c CreateOrder) ToDomain() domain.Order {
	products := c.Products.ToDomain()
	shippingPrice := domain.NewMoney(c.Shipping.Price.Amount, c.Shipping.Price.Scale, c.Shipping.Price.Currency)
	shipping := domain.NewShipping(c.Shipping.ID, shippingPrice)
	subTotal := domain.NewMoney(c.SubTotal.Amount, c.SubTotal.Scale, c.SubTotal.Currency)
	total := domain.NewMoney(c.Total.Amount, c.Total.Scale, c.Total.Currency)

	return domain.NewOrder(c.StoreID, c.ConsumerID, products, shipping, subTotal, total)
}

type OrderService interface {
	Create(ctx context.Context, create CreateOrder) error
}

type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) error
}
