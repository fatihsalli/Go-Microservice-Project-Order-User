package controller

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/service"
	"OrderUserProject/pkg"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type UserHandler struct {
	Service service.IUserService
}

func NewUserHandler(e *echo.Echo, service service.IUserService) *UserHandler {
	router := e.Group("api/users")
	b := &UserHandler{Service: service}

	//Routes
	router.GET("", b.GetAllUsers)
	router.POST("", b.CreateUser)

	return b
}

// GetAllBooks => To get request for listing all of books

// GetAllUsers godoc
// @Summary get all items in the user list
// @ID get-all-users
// @Produce json
// @Success 200 {array} models.JSONSuccessResultData
// @Success 500 {object} pkg.InternalServerError
// @Router /users [get]
func (h UserHandler) GetAllUsers(c echo.Context) error {
	userList, err := h.Service.GetAll()

	if err != nil {
		log.Printf("StatusInternalServerError: %v", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// we can use automapper, but it will cause performance loss.
	var userResponse models.UserResponse
	var usersResponse []models.UserResponse

	for _, user := range userList {
		userResponse.ID = user.ID
		userResponse.Name = user.Name
		userResponse.Email = user.Email
		userResponse.Password = user.Password
		userResponse.Address = user.Address
		userResponse.Orders = user.Orders

		usersResponse = append(usersResponse, userResponse)
	}

	// to response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(usersResponse),
		Data:           usersResponse,
	}

	log.Print("All books are listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// CreateBook => To post request for creating new a book

// CreateUser godoc
// @Summary add a new item to the user list
// @ID create-user
// @Produce json
// @Param data body models.UserCreateRequest true "book data"
// @Success 201 {object} models.JSONSuccessResultId
// @Success 400 {object} pkg.BadRequestError
// @Success 500 {object} pkg.InternalServerError
// @Router /users [post]
func (h UserHandler) CreateUser(c echo.Context) error {

	var userRequest models.UserCreateRequest

	// We parse the data as json into the struct
	if err := c.Bind(&userRequest); err != nil {
		log.Printf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	var user models.User

	// we can use automapper, but it will cause performance loss.
	user.Name = userRequest.Name
	user.Email = userRequest.Email
	user.Password = userRequest.Password
	user.Address = userRequest.Address

	result, err := h.Service.Insert(user)

	if err != nil {
		log.Printf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      result.ID,
		Success: true,
	}

	log.Printf("{%v} with id is created.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusCreated, jsonSuccessResultId)
}
