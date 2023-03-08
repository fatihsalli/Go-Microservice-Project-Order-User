package main

import (
	"OrderUserProject/internal"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	config := configs.GetConfig("test")

	mongoUserCollection := configs.ConnectDB(config.Database.Connection).
		Database(config.Database.DatabaseName).Collection(config.Database.UserCollectionName)

	UserRepo := repository.NewRepository(mongoUserCollection, func() interface{} {
		return &models.User{}
	})

	internal.NewBookHandler(e, UserRepo)

	e.Logger.Fatal(e.Start(config.Server.Port))
}
