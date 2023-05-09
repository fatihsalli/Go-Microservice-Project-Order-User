package order_api

import "time"

type OrderCreateRequest struct {
	UserId         string `json:"userId" bson:"userId" validate:"required,uuid4"`
	Status         string `json:"status" bson:"status" validate:"required,min=1,max=100"`
	Address        string `json:"address" bson:"address" validate:"required,uuid4"`
	InvoiceAddress string `json:"invoiceAddress" bson:"invoiceAddress" validate:"required,uuid4"`
	Product        []struct {
		Name     string  `json:"name" bson:"name" validate:"required,min=1,max=100"`
		Quantity int     `json:"quantity" bson:"quantity" validate:"required"`
		Price    float64 `json:"price" bson:"price" validate:"required"`
	} `json:"product" bson:"product" validate:"required"`
}

type OrderUpdateRequest struct {
	ID             string `json:"id" bson:"_id" validate:"required,uuid4"`
	UserId         string `json:"userId" bson:"userId" validate:"required,uuid4"`
	Status         string `json:"status" bson:"status" validate:"required,min=1,max=100"`
	Address        string `json:"address" bson:"address" validate:"required,uuid4"`
	InvoiceAddress string `json:"invoiceAddress" bson:"invoiceAddress" validate:"required,uuid4"`
	Product        []struct {
		Name     string  `json:"name" bson:"name" validate:"required,min=1,max=100"`
		Quantity int     `json:"quantity" bson:"quantity" validate:"required"`
		Price    float64 `json:"price" bson:"price" validate:"required"`
	} `json:"product" bson:"product" validate:"required"`
}

type OrderResponse struct {
	ID             string          `json:"id" bson:"_id"`
	UserId         string          `json:"userId" bson:"userId"`
	Status         string          `json:"status" bson:"status"`
	Address        AddressResponse `json:"address" bson:"address"`
	InvoiceAddress AddressResponse `json:"invoiceAddress" bson:"invoiceAddress"`
	Product        []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	} `json:"product" bson:"product"`
	Total     float64   `json:"total" bson:"total"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type AddressResponse struct {
	ID       string   `json:"id"`
	Address  string   `json:"address" bson:"address"`
	City     string   `json:"city" bson:"city"`
	District string   `json:"district" bson:"district"`
	Type     []string `json:"type" bson:"type"`
	Default  struct {
		IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
		IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
	} `json:"default" bson:"default"`
}

type UserResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Email     string            `json:"email"`
	Addresses []AddressResponse `json:"addresses"`
}

type OrderResponseForElastic struct {
	OrderID string `json:"orderID" bson:"orderID"`
	Status  string `json:"status" bson:"status"`
}

type OrderGenericResponse struct {
	ID               string `json:"id,omitempty" bson:"_id"`
	UserId           string `json:"userId,omitempty" bson:"userId"`
	Status           string `json:"status,omitempty" bson:"status"`
	AddressID        string `json:"addressID,omitempty" bson:"addressID"`
	InvoiceAddressID string `json:"invoiceAddressID,omitempty" bson:"invoiceAddressID"`
	Product          []struct {
		Name     string  `json:"name,omitempty" bson:"name"`
		Quantity int     `json:"quantity,omitempty" bson:"quantity"`
		Price    float64 `json:"price,omitempty" bson:"price"`
	} `json:"product,omitempty" bson:"product"`
	Total     float64 `json:"total,omitempty" bson:"total"`
	CreatedAt string  `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt string  `json:"updatedAt,omitempty" bson:"updatedAt"`
}

type OrderGetRequest struct {
	ExactFilters map[string][]interface{} `json:"exact_filters"`
	Fields       []string                 `json:"fields"`
	Match        []struct {
		MatchField string      `json:"match_field"`
		Parameter  string      `json:"parameter"`
		Value      interface{} `json:"value"`
	} `json:"match"`
	Sort map[string]int `json:"sort"`
}

type ConfigGenericEndpoint struct {
	ExactFilterArea      map[string]string
	MatchFilterParameter map[string]string
}

var ConfigsGeneric = map[string]ConfigGenericEndpoint{
	"mongoDB": {
		ExactFilterArea: map[string]string{
			"id":                      "_id",
			"_id":                     "_id",
			"userId":                  "userId",
			"userID":                  "userId",
			"status":                  "status",
			"address":                 "address",
			"address.id":              "address.id",
			"address.address":         "address.address",
			"address.city":            "address.city",
			"address.district":        "address.district",
			"address.type":            "address.type",
			"address.default":         "address.default",
			"invoiceAddress":          "invoiceAddress",
			"invoiceAddress.id":       "invoiceAddress.id",
			"invoiceAddress.address":  "invoiceAddress.address",
			"invoiceAddress.city":     "invoiceAddress.city",
			"invoiceAddress.district": "invoiceAddress.district",
			"invoiceAddress.type":     "invoiceAddress.type",
			"invoiceAddress.default":  "invoiceAddress.default",
			"product":                 "product",
			"product.name":            "product.name",
			"product.quantity":        "product.quantity",
			"product.price":           "product.price",
			"total":                   "total",
			"createdAt":               "createdAt",
			"createdAT":               "createdAt",
			"updatedAt":               "updatedAt",
			"updatedAT":               "updatedAt",
		}, MatchFilterParameter: map[string]string{
			"equal":            "$eq",
			"eq":               "$eq",
			"notEqual":         "$ne",
			"ne":               "$ne",
			"greaterThan":      "$gt",
			"gt":               "$gt",
			"greaterThanEqual": "$gte",
			"gte":              "$gte",
			"lessThan":         "$lt",
			"lt":               "$lt",
			"lessThanEqual":    "$lte",
			"lte":              "$lte",
			"in":               "$in",
			"nin":              "$nin",
			"exists":           "$exists",
			"regex":            "$regex",
		}},
}

func GetGenericConfig(database string) ConfigGenericEndpoint {
	return ConfigsGeneric[database]
}
