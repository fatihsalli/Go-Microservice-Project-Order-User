package main

import (
	"OrderUserProject/cmd"
	"github.com/labstack/gommon/log"
	"os"
)

func main() {
	project := os.Getenv("project")

	if project == "orderAPI" {
		cmd.StartOrderAPI()
	} else if project == "userAPI" {
		cmd.StartUserAPI()
	} else {
		log.Fatal("Project cannot start!")
	}
}
