package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  []byte    `json:"password" bson:"password"`
	Addresses []Address `json:"addresses" bson:"addresses"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Order struct {
	ID             string  `json:"id" bson:"_id"`
	UserId         string  `json:"userId" bson:"userId"`
	Status         string  `json:"status" bson:"status"`
	Address        Address `json:"address" bson:"address"`
	InvoiceAddress Address `json:"invoiceAddress" bson:"invoiceAddress"`
	Product        []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	} `json:"product" bson:"product"`
	Total     float64   `json:"total" bson:"total"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Address struct {
	ID       string   `json:"id" bson:"_id"`
	Address  string   `json:"address" bson:"address"`
	City     string   `json:"city" bson:"city"`
	District string   `json:"district" bson:"district"`
	Type     []string `json:"type" bson:"type"` // 2 alan gelebilir invoice veya regular ikisi de olabilir
	Default  struct {
		IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
		IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
	} `json:"default" bson:"default"`
}

// TODO: Tüm defaultları kaldırmamalı en son 1 tane değer kaldığında default olarak kalmalı kaldıramamalı
