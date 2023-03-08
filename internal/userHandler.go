package internal

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"OrderUserProject/pkg"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

type UserHandler struct {
	Repository *repository.Repository
}

func NewBookHandler(e *echo.Echo, repository *repository.Repository) *UserHandler {
	router := e.Group("api/users")
	b := &UserHandler{Repository: repository}

	//Routes
	router.GET("", b.GetAllUsers)
	router.POST("", b.CreateUser)
	return b
}

func (h UserHandler) GetAllUsers(c echo.Context) error {
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

func (h UserHandler) CreateUser(c echo.Context) error {

	var user models.User

	// We parse the data as json into the struct
	if err := c.Bind(&user); err != nil {
		log.Printf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// to create id and created date value
	user.ID = uuid.New().String()
	user.CreatedDate = primitive.NewDateTimeFromTime(time.Now())

	err := h.Repository.Create(ctx, user)

	if err != nil {
		log.Printf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      "test",
		Success: true,
	}

	log.Printf("{%v} with id is created.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusCreated, jsonSuccessResultId)
}
