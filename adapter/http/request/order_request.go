package request

import (
	"github.com/kaduartur/go-planet/core/port"
)

type Money struct {
	Amount   int    `json:"amount"`
	Scale    int    `json:"scale"`
	Currency string `json:"currency"`
}

func (m Money) ToMoneyCommand() port.Money {
	return port.Money{
		Amount:   m.Amount,
		Scale:    m.Scale,
		Currency: m.Currency,
	}
}

type Product struct {
	ID       string `json:"id"`
	SkuID    string `json:"sku_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    Money  `json:"price"`
}

func (p Product) ToProductCommand() port.Product {
	return port.Product{
		ID:       p.ID,
		SkuID:    p.SkuID,
		Name:     p.Name,
		Quantity: p.Quantity,
		Price:    p.Price.ToMoneyCommand(),
	}
}

type Products []Product

func (pp Products) ToProductsCommand() port.Products {
	products := make(port.Products, 0)
	for _, p := range pp {
		product := p.ToProductCommand()
		products = append(products, product)
	}

	return products
}

type CreateOrder struct {
	StoreID    string   `json:"store_id"`
	ConsumerID string   `json:"consumer_id"`
	Products   Products `json:"products"`
	Shipping   struct {
		ID    string `json:"id"`
		Price Money  `json:"price"`
	} `json:"shipping"`
	SubTotal Money `json:"sub_total"`
	Total    Money `json:"total"`
}

func (c CreateOrder) ToCommand() port.CreateOrder {
	return port.CreateOrder{
		StoreID:    c.StoreID,
		ConsumerID: c.ConsumerID,
		Products:   c.Products.ToProductsCommand(),
		Shipping: struct {
			ID    string
			Price port.Money
		}{
			ID:    c.Shipping.ID,
			Price: c.Shipping.Price.ToMoneyCommand(),
		},
		SubTotal: c.SubTotal.ToMoneyCommand(),
		Total:    c.Total.ToMoneyCommand(),
	}
}
