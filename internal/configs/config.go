package configs

type Config struct {
	Server struct {
		Port map[string]string
		Host string
	}
	Database struct {
		Connection          string
		DatabaseName        string
		UserCollectionName  string
		OrderCollectionName string
	}
	Elasticsearch struct {
		Addresses map[string]string
		IndexName map[string]string
	}
	Kafka struct {
		Address   string
		TopicName map[string]string
	}
	HttpClient struct {
		UserAPI  string
		OrderAPI string
	}
}

var Configs = map[string]Config{
	"test": {
		Server: struct {
			Port map[string]string
			Host string
		}{
			Port: map[string]string{
				"orderAPI": ":8011",
				"userAPI":  ":8012",
			},
			Host: "localhost",
		},
		Database: struct {
			Connection          string
			DatabaseName        string
			UserCollectionName  string
			OrderCollectionName string
		}{
			Connection:          "mongodb://localhost:27017",
			DatabaseName:        "ProjectDB",
			UserCollectionName:  "Users",
			OrderCollectionName: "Orders",
		},
		Elasticsearch: struct {
			Addresses map[string]string
			IndexName map[string]string
		}{
			Addresses: map[string]string{
				"Address 1": "http://localhost:9200",
			},
			IndexName: map[string]string{
				"OrderSave": "order_duplicate_v01",
			},
		},
		Kafka: struct {
			Address   string
			TopicName map[string]string
		}{
			Address: "localhost:9092",
			TopicName: map[string]string{
				"OrderID":    "orderID-created-v01",
				"OrderModel": "orderDuplicate-created-v01",
			},
		},
		HttpClient: struct {
			UserAPI  string
			OrderAPI string
		}{
			UserAPI:  "http://localhost:8012/api/users",
			OrderAPI: "http://localhost:8011/api/orders",
		},
	},
	"development": {
		Server: struct {
			Port map[string]string
			Host string
		}{
			Port: map[string]string{
				"orderAPI": ":8011",
				"userAPI":  ":8012",
			},
			Host: "",
		},
		Database: struct {
			Connection          string
			DatabaseName        string
			UserCollectionName  string
			OrderCollectionName string
		}{
			Connection:          "mongodb://172.28.0.51:27017",
			DatabaseName:        "ProjectDB",
			UserCollectionName:  "Users",
			OrderCollectionName: "Orders",
		},
		Elasticsearch: struct {
			Addresses map[string]string
			IndexName map[string]string
		}{
			Addresses: map[string]string{
				"Address 1": "http://172.28.0.55:9200",
			},
			IndexName: map[string]string{
				"OrderSave": "order_duplicate_v01",
			},
		},
		Kafka: struct {
			Address   string
			TopicName map[string]string
		}{
			Address: "172.28.0.53:9092",
			TopicName: map[string]string{
				"OrderID":    "orderID-created-v01",
				"OrderModel": "orderDuplicate-created-v01",
			},
		},
		HttpClient: struct {
			UserAPI  string
			OrderAPI string
		}{
			UserAPI:  "http://user-api:8012/api/users",
			OrderAPI: "http://order-api:8011/api/orders",
		},
	},
	"qa":   {},
	"prod": {},
}

func GetConfig(env string) Config {
	if conf, ok := Configs[env]; ok {
		return conf
	}

	return Configs["test"]
}

type GenericEndpointConfig struct {
	ExactFilterArea      map[string]string
	MatchFilterParameter map[string]string
}

var GenericEndpointConfigs = map[string]GenericEndpointConfig{
	"mongoDB": {
		ExactFilterArea: map[string]string{
			"id":                      "_id",
			"_id":                     "_id",
			"userId":                  "userId",
			"userID":                  "userId",
			"status":                  "status",
			"product.name":            "product.name",
			"product.quantity":        "product.quantity",
			"product.price":           "product.price",
			"total":                   "total",
			"createdAt":               "createdAt",
			"createdAT":               "createdAt",
			"updatedAt":               "updatedAt",
			"updatedAT":               "updatedAt",
			"address.id":              "address.id",
			"address.address":         "address.address",
			"address.city":            "address.city",
			"address.district":        "address.district",
			"address.type":            "address.type",
			"invoiceAddress.id":       "invoiceAddress.id",
			"invoiceAddress.address":  "invoiceAddress.address",
			"invoiceAddress.city":     "invoiceAddress.city",
			"invoiceAddress.district": "invoiceAddress.district",
			"invoiceAddress.type":     "invoiceAddress.type",
			"address.default.isDefaultInvoiceAddress":        "address.default.isDefaultInvoiceAddress",
			"address.default.isDefaultRegularAddress":        "address.default.isDefaultRegularAddress",
			"invoiceAddress.default.isDefaultInvoiceAddress": "invoiceAddress.default.isDefaultInvoiceAddress",
			"invoiceAddress.default.isDefaultRegularAddress": "invoiceAddress.default.isDefaultRegularAddress",
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
	"elasticsearch": {
		ExactFilterArea: map[string]string{
			"id":                      "id.keyword",
			"_id":                     "id.keyword",
			"userId":                  "userId.keyword",
			"userID":                  "userId.keyword",
			"status":                  "status.keyword",
			"product.name":            "product.name.keyword",
			"product.quantity":        "product.quantity",
			"product.price":           "product.price",
			"total":                   "total",
			"createdAt":               "createdAt",
			"createdAT":               "createdAt",
			"updatedAt":               "updatedAt",
			"updatedAT":               "updatedAt",
			"address.id":              "address.id.keyword",
			"address.address":         "address.address.keyword",
			"address.city":            "address.city.keyword",
			"address.district":        "address.district.keyword",
			"address.type":            "address.type.keyword",
			"invoiceAddress.id":       "invoiceAddress.id.keyword",
			"invoiceAddress.address":  "invoiceAddress.address.keyword",
			"invoiceAddress.city":     "invoiceAddress.city.keyword",
			"invoiceAddress.district": "invoiceAddress.district.keyword",
			"invoiceAddress.type":     "invoiceAddress.type.keyword",
			"address.default.isDefaultInvoiceAddress":        "address.default.isDefaultInvoiceAddress",
			"address.default.isDefaultRegularAddress":        "address.default.isDefaultRegularAddress",
			"invoiceAddress.default.isDefaultInvoiceAddress": "invoiceAddress.default.isDefaultInvoiceAddress",
			"invoiceAddress.default.isDefaultRegularAddress": "invoiceAddress.default.isDefaultRegularAddress",
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

func GetGenericEndpointConfig(database string) GenericEndpointConfig {
	return GenericEndpointConfigs[database]
}
