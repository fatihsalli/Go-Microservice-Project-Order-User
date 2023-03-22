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
}

var Configs = map[string]Config{
	"test": {
		Server: struct {
			Port map[string]string
			Host string
		}{
			Port: map[string]string{
				"orderAPI":     ":8011",
				"userAPI":      ":8012",
				"orderElastic": ":8013",
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
