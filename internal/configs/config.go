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
			Connection:          "mongodb://172.18.0.2:27017",
			DatabaseName:        "ProjectDB",
			UserCollectionName:  "Users",
			OrderCollectionName: "Orders",
		},
		Elasticsearch: struct {
			Addresses map[string]string
			IndexName map[string]string
		}{
			Addresses: map[string]string{
				"Address 1": "http://elasticsearch:9200",
			},
			IndexName: map[string]string{
				"OrderSave": "order_duplicate_v01",
			},
		},
		Kafka: struct {
			Address   string
			TopicName map[string]string
		}{
			Address: "kafka:9092",
			TopicName: map[string]string{
				"OrderID":    "orderID-created-v01",
				"OrderModel": "orderDuplicate-created-v01",
			},
		},
		HttpClient: struct {
			UserAPI  string
			OrderAPI string
		}{
			UserAPI:  "http://user-api:80/api/users",
			OrderAPI: "http://order-api:80/api/orders",
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
