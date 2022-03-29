package mysql

import (
	"fmt"
	"time"

	"github.com/kaduartur/go-planet/pkg/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	ID               int      `gorm:"primary_key" json:"id"`
	StoreID          string   `gorm:"type:string;index;not null" json:"store_id"`
	ConsumerID       string   `gorm:"type:string;index;not null" json:"consumer_id"`
	Products         Products `gorm:"foreignKey:OrderID" json:"products"`
	ShippingID       string   `gorm:"type:string;not null" json:"shipping_id"`
	ShippingAmount   int      `gorm:"type:int;not null" json:"shipping_amount"`
	ShippingScale    int      `gorm:"type:int;not null" json:"shipping_scale"`
	ShippingCurrency string   `gorm:"type:string;not null" json:"shipping_currency"`
	SubTotalAmount   int      `gorm:"type:int;not null" json:"sub_total_amount"`
	SubTotalScale    int      `gorm:"type:int;not null" json:"sub_total_scale"`
	SubTotalCurrency string   `gorm:"type:string;not null" json:"sub_total_currency"`
	TotalAmount      int      `gorm:"type:int;not null" json:"total_amount"`
	TotalScale       int      `gorm:"type:int;not null" json:"total_scale"`
	TotalCurrency    string   `gorm:"type:string;not null" json:"total_currency"`
}

type Product struct {
	ID            int    `gorm:"primary_key" json:"id"`
	OrderID       int    `gorm:"type:int;not null" json:"order_id"`
	SkuID         string `gorm:"type:string;not null" json:"sku_id"`
	Name          string `gorm:"type:string;not null" json:"name"`
	TotalAmount   int    `gorm:"type:int;not null" json:"total_amount"`
	TotalScale    int    `gorm:"type:int;not null" json:"total_scale"`
	TotalCurrency string `gorm:"type:string;not null" json:"total_currency"`
}

type Products []Product

func NewConnection(cfg env.Database) (*gorm.DB, error) {
	mysqlCredentials := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		cfg.Username,
		cfg.Password,
		cfg.HostName,
		cfg.Port,
		cfg.Name,
	)
	db, err := gorm.Open(mysql.Open(mysqlCredentials))
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxLifetime(time.Second * 10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	doMigration(db)

	return db, nil
}

func doMigration(db *gorm.DB) {
	if !db.Migrator().HasTable(&Order{}) {
		err := db.Migrator().CreateTable(&Order{})
		if err != nil {
			panic(err)
		}
	}
	if !db.Migrator().HasTable(&Product{}) {
		err := db.Migrator().CreateTable(&Product{})
		if err != nil {
			panic(err)
		}
	}
}
