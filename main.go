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
	} else if project == "orderElastic" {
		cmd.StartOrderElastic()
	} else {
		log.Fatal("Project cannot start!")
	}
}
