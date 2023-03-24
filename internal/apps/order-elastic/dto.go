package order_elastic

type OrderResponseElastic struct {
	ID             string
	UserId         string          `json:"userId" bson:"userId"`
	Status         string          `json:"status" bson:"status"`
	Address        AddressResponse `json:"address" bson:"address"`
	InvoiceAddress AddressResponse `json:"invoiceAddress" bson:"invoiceAddress"`
	Product        []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	} `json:"product" bson:"product"`
	Total float64 `json:"total" bson:"total"`
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
