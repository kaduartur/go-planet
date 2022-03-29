package domain

import (
	"math/rand"
)

func newOrderID() OrderID {
	return OrderID{rand.Int()}
}

type OrderID struct {
	id int
}

func (o OrderID) Id() int {
	return o.id
}

func NewOrder(storeID, consumerID string, products Products, shipping Shipping, subTotal, total Money) Order {
	id := newOrderID()
	return Order{
		ID:         id,
		StoreID:    storeID,
		ConsumerID: consumerID,
		Products:   products,
		Shipping:   shipping,
		SubTotal:   subTotal,
		Total:      total,
	}
}

type Order struct {
	ID         OrderID
	StoreID    string
	ConsumerID string
	Products   Products
	Shipping   Shipping
	SubTotal   Money
	Total      Money
}

func NewMoney(amount, scale int, currency string) Money {
	return Money{
		Amount:   amount,
		Scale:    scale,
		Currency: currency,
	}
}

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

func NewProduct(ID string, skuID string, name string, quantity int, price Money) Product {
	return Product{
		ID:       ID,
		SkuID:    skuID,
		Name:     name,
		Quantity: quantity,
		Price:    price,
	}
}

type Products []Product

func NewShipping(ID string, price Money) Shipping {
	return Shipping{
		ID:    ID,
		Price: price,
	}
}

type Shipping struct {
	ID    string
	Price Money
}
