package internal

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Handler struct {
	Repository *repository.Repository
}

func NewBookHandler(e *echo.Echo, repository *repository.Repository) *Handler {
	router := e.Group("api/users")
	b := &Handler{Repository: repository}

	//Routes
	router.GET("", b.GetAllUsers)

	return b
}

func (h Handler) GetAllUsers(c echo.Context) error {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	userList, err := h.Repository.GetAll(ctx)

	if err != nil {

	}

	// to response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(userList),
		Data:           userList,
	}

	return c.JSON(http.StatusOK, jsonSuccessResultData)
}
