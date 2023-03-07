package configs

type Config struct {
	Server struct {
		Port string
		Host string
	}
	Database struct {
		Connection     string
		DatabaseName   string
		CollectionName string
	}
}

var Configs = map[string]Config{
	"test": {
		Server: struct {
			Port string
			Host string
		}{
			Port: ":8080",
			Host: "localhost",
		},
		Database: struct {
			Connection     string
			DatabaseName   string
			CollectionName string
		}{
			Connection:     "mongodb://localhost:27017",
			DatabaseName:   "booksDB",
			CollectionName: "books",
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
