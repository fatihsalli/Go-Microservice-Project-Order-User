package handler

import (
	"OrderUserProject/internal/apps/user-api"
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type UserHandler struct {
	Service user_api.IUserService
}

func NewUserHandler(e *echo.Echo, service user_api.IUserService) *UserHandler {
	router := e.Group("api/users")
	b := &UserHandler{Service: service}

	//Routes
	router.GET("", b.GetAllUsers)
	router.GET("/:id", b.GetUserById)
	router.POST("", b.CreateUser)
	router.PUT("", b.UpdateUser)
	router.DELETE("/:id", b.DeleteUser)

	return b
}

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
	var userResponse user_api.UserResponse
	var usersResponse []user_api.UserResponse

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

// GetUserById godoc
// @Summary get a user item by ID
// @ID get-user-by-id
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} models.JSONSuccessResultData
// @Success 404 {object} pkg.NotFoundError
// @Success 500 {object} pkg.InternalServerError
// @Router /users/{id} [get]
func (h UserHandler) GetUserById(c echo.Context) error {
	query := c.Param("id")

	user, err := h.Service.GetUserById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Not found exception: {%v} with id not found!", query)
			return c.JSON(http.StatusNotFound, pkg.NotFoundError{
				Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
			})
		}
		log.Printf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// mapping
	var userResponse user_api.UserResponse
	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	userResponse.Password = user.Password
	userResponse.Address = user.Address
	userResponse.Orders = user.Orders

	// to response success result data => single one
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: 1,
		Data:           userResponse,
	}

	log.Printf("{%v} with id is listed.", userResponse.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// CreateUser godoc
// @Summary add a new item to the user list
// @ID create-user
// @Produce json
// @Param data body user_api.UserCreateRequest true "user data"
// @Success 201 {object} models.JSONSuccessResultId
// @Success 400 {object} pkg.BadRequestError
// @Success 500 {object} pkg.InternalServerError
// @Router /users [post]
func (h UserHandler) CreateUser(c echo.Context) error {

	var userRequest user_api.UserCreateRequest

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

// UpdateUser godoc
// @Summary update an item to the user list
// @ID update-user
// @Produce json
// @Param data body user_api.UserUpdateRequest true "user data"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 400 {object} pkg.BadRequestError
// @Success 500 {object} pkg.InternalServerError
// @Router /users [put]
func (h UserHandler) UpdateUser(c echo.Context) error {

	var userUpdateRequest user_api.UserUpdateRequest

	// we parse the data as json into the struct
	if err := c.Bind(&userUpdateRequest); err != nil {
		log.Printf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	if _, err := h.Service.GetUserById(userUpdateRequest.ID); err != nil {
		log.Printf("Not found exception: {%v} with id not found!", userUpdateRequest.ID)
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", userUpdateRequest.ID),
		})
	}

	var user models.User

	// we can use automapper, but it will cause performance loss.
	user.ID = userUpdateRequest.ID
	user.Name = userUpdateRequest.Name
	user.Email = userUpdateRequest.Email
	user.Password = userUpdateRequest.Password
	user.Address = userUpdateRequest.Address

	result, err := h.Service.Update(user)

	if err != nil || result == false {
		log.Printf("StatusInternalServerError: {%v} ", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      user.ID,
		Success: result,
	}

	log.Printf("{%v} with id is updated.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}

// DeleteUser godoc
// @Summary delete a user item by ID
// @ID delete-user-by-id
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 404 {object} pkg.NotFoundError
// @Router /users/{id} [delete]
func (h UserHandler) DeleteUser(c echo.Context) error {
	query := c.Param("id")

	result, err := h.Service.Delete(query)

	if err != nil || result == false {
		log.Printf("Not found exception: {%v} with id not found!", query)
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      query,
		Success: result,
	}

	log.Printf("{%v} with id is deleted.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}
