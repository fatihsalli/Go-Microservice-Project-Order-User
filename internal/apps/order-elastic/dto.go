package order_elastic

import "time"

type OrderResponseElastic struct {
	ID             string          `json:"orderId"`
	UserId         string          `json:"userId"`
	Status         string          `json:"status"`
	Address        AddressResponse `json:"address"`
	InvoiceAddress AddressResponse `json:"invoiceAddress"`
	Product        []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	} `json:"product" bson:"product"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AddressResponse struct {
	Address  string   `json:"address" bson:"address"`
	City     string   `json:"city" bson:"city"`
	District string   `json:"district" bson:"district"`
	Type     []string `json:"type" bson:"type"`
	Default  struct {
		IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
		IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
	} `json:"default" bson:"default"`
}
