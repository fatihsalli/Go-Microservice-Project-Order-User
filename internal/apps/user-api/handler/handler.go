package handler

import (
	"OrderUserProject/internal/apps/user-api"
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserHandler struct {
	Service   *user_api.UserService
	Validator *validator.Validate
}

func NewUserHandler(e *echo.Echo, service *user_api.UserService, v *validator.Validate) *UserHandler {
	router := e.Group("api/users")
	b := &UserHandler{Service: service, Validator: v}

	//Routes
	router.GET("", b.GetAllUsers)
	router.GET("/:id", b.GetUserById)
	router.POST("", b.CreateUser)
	router.PUT("", b.UpdateUser)
	router.PUT("/add-address/:id", b.AddAddress)
	router.PUT("/change-address/:id", b.ChangeAddress)
	router.PUT("/delete-address/:id", b.DeleteAddress)
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
func (h *UserHandler) GetAllUsers(c echo.Context) error {
	userList, err := h.Service.GetAll()

	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// we can use automapper, but it will cause performance loss.
	var userResponse user_api.UserResponse
	var usersResponse []user_api.UserResponse
	var addressResponse user_api.AddressResponse
	for _, user := range userList {
		userResponse.ID = user.ID
		userResponse.Name = user.Name
		userResponse.Email = user.Email
		for _, address := range user.Addresses {
			addressResponse.ID = address.ID
			addressResponse.Address = address.Address
			addressResponse.City = address.City
			addressResponse.District = address.District
			addressResponse.Type = address.Type
			addressResponse.Default = address.Default
			userResponse.Addresses = append(userResponse.Addresses, addressResponse)
		}
		usersResponse = append(usersResponse, userResponse)
	}

	// to response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(usersResponse),
		Data:           usersResponse,
	}

	c.Logger().Info("All users are listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// GetUserById godoc
// @Summary get a user item by ID
// @ID get-user-by-id
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} user_api.UserResponse
// @Success 404 {object} pkg.NotFoundError
// @Success 500 {object} pkg.InternalServerError
// @Router /users/{id} [get]
func (h *UserHandler) GetUserById(c echo.Context) error {
	query := c.Param("id")

	user, err := h.Service.GetUserById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Logger().Errorf("Not found exception: {%v} with id not found!", query)
			return c.JSON(http.StatusNotFound, pkg.NotFoundError{
				Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
			})
		}
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// we can use automapper, but it will cause performance loss.
	var userResponse user_api.UserResponse
	var addressResponse user_api.AddressResponse
	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	for _, address := range user.Addresses {
		addressResponse.ID = address.ID
		addressResponse.Address = address.Address
		addressResponse.City = address.City
		addressResponse.District = address.District
		addressResponse.Type = address.Type
		addressResponse.Default = address.Default
		userResponse.Addresses = append(userResponse.Addresses, addressResponse)
	}

	c.Logger().Infof("{%v} with id is listed.", userResponse.ID)
	return c.JSON(http.StatusOK, userResponse)
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
func (h *UserHandler) CreateUser(c echo.Context) error {

	var userRequest user_api.UserCreateRequest

	// We parse the data as json into the struct
	if err := c.Bind(&userRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Check address
	if len(userRequest.Addresses) < 1 {
		c.Logger().Error("Address value is empty.")
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: "Address value cannot empty. At least you have to put one address!",
		})
	}

	// Validate user input using the validator instance
	if err := h.Validator.Struct(userRequest); err != nil {
		c.Logger().Errorf("Bad Request. Please put valid user model ! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. Please put valid user model! %v", err.Error()),
		})
	}

	// we can use automapper, but it will cause performance loss.
	var user models.User
	var address models.Address
	user.Name = userRequest.Name
	user.Email = userRequest.Email
	for _, addressRequest := range userRequest.Addresses {
		address.ID = uuid.New().String()
		address.Address = addressRequest.Address
		address.City = addressRequest.City
		address.District = addressRequest.District
		address.Type = addressRequest.Type
		address.Default = addressRequest.Default
		user.Addresses = append(user.Addresses, address)
	}

	// using 'bcrypt' to hash password
	password := []byte(userRequest.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Errorf("Bad Request. It cannot be hashing! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be hashing! %v", err.Error()),
		})
	}
	user.Password = hashedPassword

	// Invoice and regular addresses check
	userCheckAddress := h.Service.InvoiceRegularAddressCheck(user)

	result, err := h.Service.Insert(userCheckAddress)

	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      result.ID,
		Success: true,
	}

	c.Logger().Infof("{%v} with id is created.", jsonSuccessResultId.ID)
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
func (h *UserHandler) UpdateUser(c echo.Context) error {

	var userUpdateRequest user_api.UserUpdateRequest

	// we parse the data as json into the struct
	if err := c.Bind(&userUpdateRequest); err != nil {
		c.Logger().Errorf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Validate user input using the validator instance
	if err := h.Validator.Struct(userUpdateRequest); err != nil {
		c.Logger().Errorf("Bad Request. Please put valid user model ! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. Please put valid user model! %v", err.Error()),
		})
	}

	// to find user
	userExist, err := h.Service.GetUserById(userUpdateRequest.ID)
	if err != nil {
		c.Logger().Errorf("Not found exception: {%v} with id not found!", userUpdateRequest.ID)
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", userUpdateRequest.ID),
		})
	}

	// we can use automapper, but it will cause performance loss.
	var user models.User
	user.ID = userUpdateRequest.ID
	user.Name = userUpdateRequest.Name
	user.Email = userUpdateRequest.Email
	user.Addresses = userExist.Addresses

	// using 'bcrypt' to check password (tested)
	err = bcrypt.CompareHashAndPassword(userExist.Password, []byte(userUpdateRequest.Password))
	if err != nil {
		c.Logger().Error("Password is wrong. Please put correct password!")
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprint("Password is wrong. Please put correct password!"),
		})
	}

	result, err := h.Service.Update(user)

	if err != nil || result == false {
		c.Logger().Errorf("StatusInternalServerError: {%v} ", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "User cannot update! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      user.ID,
		Success: result,
	}

	c.Logger().Infof("{%v} with id is updated.", jsonSuccessResultId.ID)
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
func (h *UserHandler) DeleteUser(c echo.Context) error {
	query := c.Param("id")

	result, err := h.Service.Delete(query)

	if err != nil || result == false {
		c.Logger().Errorf("Not found exception: {%v} with id not found!", query)
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      query,
		Success: result,
	}

	c.Logger().Infof("{%v} with id is deleted.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}

// AddAddress godoc
// @Summary add a user's address by userID
// @ID add-address-with-userID
// @Produce json
// @Param id path string true "user ID"
// @Param data body user_api.AddressCreateRequest true "address data"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 404 {object} pkg.NotFoundError
// @Success 500 {object} pkg.InternalServerError
// @Router /users/add-address/{id} [put]
func (h *UserHandler) AddAddress(c echo.Context) error {
	query := c.Param("id")

	user, err := h.Service.GetUserById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Logger().Errorf("Not found exception: {%v} with id not found!", query)
			return c.JSON(http.StatusNotFound, pkg.NotFoundError{
				Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
			})
		}
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	var userAddress user_api.AddressCreateRequest

	// we parse the data as json into the struct
	if err := c.Bind(&userAddress); err != nil {
		c.Logger().Errorf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Validate user input using the validator instance
	if err := h.Validator.Struct(userAddress); err != nil {
		c.Logger().Errorf("Bad Request. Please put valid user address model ! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. Please put valid user address model! %v", err.Error()),
		})
	}

	var userAddressModel models.Address
	userAddressModel.ID = uuid.New().String()
	userAddressModel.Address = userAddress.Address
	userAddressModel.City = userAddress.City
	userAddressModel.District = userAddress.District
	userAddressModel.Type = userAddress.Type
	userAddressModel.Default = userAddress.Default

	user.Addresses = append(user.Addresses, userAddressModel)

	userAddressCheck := h.Service.InvoiceRegularAddressCheck(user)

	result, err := h.Service.Update(userAddressCheck)

	if err != nil || result == false {
		c.Logger().Errorf("StatusInternalServerError: {%v} ", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "User cannot update! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      userAddressModel.ID,
		Success: result,
	}

	c.Logger().Infof("{%v} with id is updated.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}

// ChangeAddress godoc
// @Summary change a user's address by userID
// @ID change-address-with-userID
// @Produce json
// @Param id path string true "user ID"
// @Param data body user_api.AddressUpdateRequest true "address data"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 404 {object} pkg.NotFoundError
// @Success 500 {object} pkg.InternalServerError
// @Router /users/change-address/{id} [put]
func (h *UserHandler) ChangeAddress(c echo.Context) error {
	query := c.Param("id")

	user, err := h.Service.GetUserById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Logger().Errorf("Not found exception: {%v} with id not found!", query)
			return c.JSON(http.StatusNotFound, pkg.NotFoundError{
				Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
			})
		}
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	var userAddress user_api.AddressUpdateRequest

	// we parse the data as json into the struct
	if err := c.Bind(&userAddress); err != nil {
		c.Logger().Errorf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Validate user input using the validator instance
	if err := h.Validator.Struct(userAddress); err != nil {
		c.Logger().Errorf("Bad Request. Please put valid user address model ! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. Please put valid user address model! %v", err.Error()),
		})
	}

	var userAddressModel models.Address
	userAddressModel.ID = userAddress.ID
	userAddressModel.Address = userAddress.Address
	userAddressModel.City = userAddress.City
	userAddressModel.District = userAddress.District
	userAddressModel.Type = userAddress.Type
	userAddressModel.Default = userAddress.Default

	for i, address := range user.Addresses {
		if address.ID == userAddressModel.ID {
			user.Addresses[i] = userAddressModel
		}
	}

	userAddressCheck := h.Service.InvoiceRegularAddressCheck(user)

	result, err := h.Service.Update(userAddressCheck)

	if err != nil || result == false {
		c.Logger().Errorf("StatusInternalServerError: {%v} ", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "User cannot update! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      userAddressModel.ID,
		Success: result,
	}

	c.Logger().Infof("{%v} with id is updated.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}

// DeleteAddress godoc
// @Summary delete a user's address by userID
// @ID delete-address-with-userID
// @Produce json
// @Param id path string true "user ID"
// @Param data body user_api.AddressDeleteRequest true "address data"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 404 {object} pkg.NotFoundError
// @Success 500 {object} pkg.InternalServerError
// @Router /users/delete-address/{id} [put]
func (h *UserHandler) DeleteAddress(c echo.Context) error {
	query := c.Param("id")

	user, err := h.Service.GetUserById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.Logger().Errorf("Not found exception: {%v} with id not found!", query)
			return c.JSON(http.StatusNotFound, pkg.NotFoundError{
				Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
			})
		}
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	if len(user.Addresses) < 2 {
		c.Logger().Errorf("BadRequestError: %v", errors.New("you cannot delete user's address"))
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: "You cannot delete user's address. Because there is just one address.",
		})
	}

	var userAddress user_api.AddressDeleteRequest

	// we parse the data as json into the struct
	if err := c.Bind(&userAddress); err != nil {
		c.Logger().Errorf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Validate user input using the validator instance
	if err := h.Validator.Struct(userAddress); err != nil {
		c.Logger().Errorf("Bad Request. Please put valid user address model ! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. Please put valid user address model! %v", err.Error()),
		})
	}

	for i, address := range user.Addresses {
		if address.ID == userAddress.AddressID {
			user.Addresses = append(user.Addresses[:i], user.Addresses[i+1:]...)
		}
	}

	userAddressCheck := h.Service.InvoiceRegularAddressCheck(user)

	result, err := h.Service.Update(userAddressCheck)

	if err != nil || result == false {
		c.Logger().Errorf("StatusInternalServerError: {%v} ", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "User cannot update! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      userAddress.AddressID,
		Success: result,
	}

	c.Logger().Infof("{%v} with id is updated.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}
