package handler

import (
	"OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/apps/order-api/graphQL"
	"OrderUserProject/internal/configs"
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg"
	"OrderUserProject/pkg/kafka"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/graphql-go/graphql"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type OrderHandler struct {
	Service        *order_api.OrderService
	ElasticService *order_api.ElasticService
	Producer       *kafka.ProducerKafka
	Config         *configs.Config
	Validator      *validator.Validate
}

func NewOrderHandler(e *echo.Echo, service *order_api.OrderService, producer *kafka.ProducerKafka, config *configs.Config, v *validator.Validate, elasticService *order_api.ElasticService) *OrderHandler {
	router := e.Group("api/orders")
	b := &OrderHandler{Service: service, Producer: producer, Config: config, Validator: v, ElasticService: elasticService}

	e.Use(pkg.CustomErrorMiddleware)

	//Routes
	router.GET("", b.GetAllOrders)
	router.GET("/:id", b.GetOrderById)
	router.GET("/GraphQLWithStatus/:status", b.GraphQLWithStatus)
	router.POST("", b.CreateOrder, pkg.CheckOrderStatus)
	router.POST("/GenericEndpointFromMongo", b.GenericEndpointFromMongo)
	router.POST("/GenericEndpointFromElastic", b.GenericEndpointFromElastic)
	router.PUT("", b.UpdateOrder, pkg.CheckOrderStatus)
	router.DELETE("/:id", b.DeleteOrder)
	return b
}

// GetAllOrders godoc
// @Summary get all items in the order list
// @ID get-all-orders
// @Produce json
// @Success 200 {array} models.JSONSuccessResultData
// @Success 500 {object} pkg.InternalServerError
// @Router /orders [get]
func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	orderList, err := h.Service.GetAll()

	if err != nil {
		internalServerError := pkg.InternalServerError{
			Message:    fmt.Sprintf("StatusInternalServerError: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
		return internalServerError
	}

	// We can use automapper, but it will cause performance loss.
	var orderResponse order_api.OrderResponse
	var ordersResponse []order_api.OrderResponse
	for _, order := range orderList {
		orderResponse.ID = order.ID
		orderResponse.UserId = order.UserId
		orderResponse.Product = order.Product
		orderResponse.Address.ID = order.Address.ID
		orderResponse.Address.Address = order.Address.Address
		orderResponse.Address.City = order.Address.City
		orderResponse.Address.District = order.Address.District
		orderResponse.Address.Type = order.Address.Type
		orderResponse.Address.Default = order.Address.Default
		orderResponse.InvoiceAddress.ID = order.InvoiceAddress.ID
		orderResponse.InvoiceAddress.Address = order.InvoiceAddress.Address
		orderResponse.InvoiceAddress.City = order.InvoiceAddress.City
		orderResponse.InvoiceAddress.District = order.InvoiceAddress.District
		orderResponse.InvoiceAddress.Type = order.InvoiceAddress.Type
		orderResponse.InvoiceAddress.Default = order.InvoiceAddress.Default
		orderResponse.Product = order.Product
		orderResponse.Total = order.Total
		orderResponse.Status = order.Status
		orderResponse.CreatedAt = order.CreatedAt
		orderResponse.UpdatedAt = order.UpdatedAt
		ordersResponse = append(ordersResponse, orderResponse)
	}

	// Response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(ordersResponse),
		Data:           ordersResponse,
	}

	c.Logger().Info("All books are successfully listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// GetOrderById godoc
// @Summary get an order item by ID
// @ID get-order-by-id
// @Produce json
// @Param id path string true "order ID"
// @Success 200 {object} order_api.OrderResponse
// @Success 404 {object} pkg.NotFoundError
// @Success 500 {object} pkg.InternalServerError
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderById(c echo.Context) error {
	query := c.Param("id")

	order, err := h.Service.GetOrderById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			notFoundErr := pkg.NotFoundError{
				Message:    fmt.Sprintf("Not found exception: {%v} with id not found!", query),
				StatusCode: http.StatusNotFound,
			}
			return notFoundErr
		}
		internalServerError := pkg.InternalServerError{
			Message:    fmt.Sprintf("StatusInternalServerError: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
		return internalServerError
	}

	// We can use automapper, but it will cause performance loss.
	var orderResponse order_api.OrderResponse
	orderResponse.ID = order.ID
	orderResponse.UserId = order.UserId
	orderResponse.Product = order.Product
	orderResponse.Address.ID = order.Address.ID
	orderResponse.Address.Address = order.Address.Address
	orderResponse.Address.City = order.Address.City
	orderResponse.Address.District = order.Address.District
	orderResponse.Address.Type = order.Address.Type
	orderResponse.Address.Default = order.Address.Default
	orderResponse.InvoiceAddress.ID = order.InvoiceAddress.ID
	orderResponse.InvoiceAddress.Address = order.InvoiceAddress.Address
	orderResponse.InvoiceAddress.City = order.InvoiceAddress.City
	orderResponse.InvoiceAddress.District = order.InvoiceAddress.District
	orderResponse.InvoiceAddress.Type = order.InvoiceAddress.Type
	orderResponse.InvoiceAddress.Default = order.InvoiceAddress.Default
	orderResponse.Product = order.Product
	orderResponse.Total = order.Total
	orderResponse.Status = order.Status
	orderResponse.CreatedAt = order.CreatedAt
	orderResponse.UpdatedAt = order.UpdatedAt

	c.Logger().Info("{%v} with id is listed.", orderResponse.ID)
	return c.JSON(http.StatusOK, orderResponse)
}

// GraphQLWithStatus godoc
// @Summary get orders by status
// @ID get-order-by-status
// @Produce json
// @Param status path string true "status"
// @Success 200 {object} order_api.OrderResponse
// @Success 404 {object} pkg.BadRequestError
// @Router /orders/GraphQLWithStatus/{status} [get]
func (h *OrderHandler) GraphQLWithStatus(c echo.Context) error {
	query := c.Param("status")

	params := graphql.Params{
		Schema:         graphQL.Schema,
		RequestString:  graphQL.GenerateGraphQLQuery(query),
		VariableValues: nil,
	}
	result := graphql.Do(params)

	if len(result.Errors) > 0 {
		badRequestErr := pkg.BadRequestError{
			Message:    fmt.Sprintf("Bad Request. It cannot be reading! %v", result.Errors),
			StatusCode: http.StatusBadRequest,
		}
		return badRequestErr
	}

	// Count orders
	ordersData := result.Data.(map[string]interface{})
	orders := ordersData["orders"].([]interface{})

	// Response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(orders),
		Data:           orders,
	}

	c.Logger().Info("Orders are successfully listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// CreateOrder godoc
// @Summary add a new item to the order list
// @ID create-order
// @Produce json
// @Param data body order_api.OrderCreateRequest true "order data"
// @Success 201 {object} models.JSONSuccessResultId
// @Success 400 {object} pkg.BadRequestError
// @Success 404 {object} pkg.NotFoundError
// @Success 500 {object} pkg.InternalServerError
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	// Get order model from middleware because we bind it within middleware
	orderRequest := c.Get("order").(*order_api.OrderCreateRequest)

	// We parse the data as json into the struct
	/*	var orderRequest order_api.OrderCreateRequest
		if err := c.Bind(&orderRequest); err != nil {
			c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
			return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
				Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
			})
		}*/

	// Validate user input using the validator instance
	if err := h.Validator.Struct(orderRequest); err != nil {
		badRequestErr := pkg.BadRequestError{
			Message:    fmt.Sprintf("Bad Request. Please put valid order model ! %v", err),
			StatusCode: http.StatusBadRequest,
		}
		return badRequestErr
	}

	// Check user with http.Client
	user, err := h.Service.GetUser(orderRequest.UserId, h.Config.HttpClient.UserAPI)

	if err != nil {
		notFoundErr := pkg.NotFoundError{
			Message:    fmt.Sprintf("Not Found Exception: %v", err),
			StatusCode: http.StatusNotFound,
		}
		return notFoundErr
	}

	// Address check
	regularAddressCheck := false
	invoiceAddressCheck := false
	for _, regularAddress := range user.Addresses {
		if regularAddress.ID == orderRequest.Address {
			regularAddressCheck = true
		}

		if regularAddress.ID == orderRequest.InvoiceAddress {
			invoiceAddressCheck = true
		}
	}

	if regularAddressCheck == false || invoiceAddressCheck == false {
		notFoundErr := pkg.NotFoundError{
			Message:    "Not Found Exception: Address not found. Before order processing please put correct address id.",
			StatusCode: http.StatusNotFound,
		}
		return notFoundErr
	}

	// Mapping => we can use automapper, but it will cause performance loss.
	var order models.Order
	order.UserId = orderRequest.UserId
	order.Status = orderRequest.Status
	for _, regularAddress := range user.Addresses {
		if regularAddress.ID == orderRequest.Address {
			order.Address.ID = regularAddress.ID
			order.Address.Address = regularAddress.Address
			order.Address.City = regularAddress.City
			order.Address.District = regularAddress.District
			order.Address.Type = regularAddress.Type
			order.Address.Default = regularAddress.Default
		}
	}
	for _, invoiceAddress := range user.Addresses {
		if invoiceAddress.ID == orderRequest.InvoiceAddress {
			order.InvoiceAddress.ID = invoiceAddress.ID
			order.InvoiceAddress.Address = invoiceAddress.Address
			order.InvoiceAddress.City = invoiceAddress.City
			order.InvoiceAddress.District = invoiceAddress.District
			order.InvoiceAddress.Type = invoiceAddress.Type
			order.InvoiceAddress.Default = invoiceAddress.Default
		}
	}
	order.Product = []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	}(orderRequest.Product)

	// Service => Insert
	result, err := h.Service.Insert(order)

	// => SEND MESSAGE (OrderID)
	var orderKafka order_api.OrderResponseForElastic
	orderKafka.OrderID = result.ID
	orderKafka.Status = "Created"

	resultJson, errJson := json.Marshal(orderKafka)
	if errJson != nil {
		c.Logger().Errorf("Something went wrong convert to byte: %v", err)
	} else {
		err = h.Producer.SendToKafkaWithMessage(resultJson, h.Config.Kafka.TopicName["OrderID"])
		if err != nil {
			c.Logger().Errorf("Something went wrong cannot pushed: %v", err)
		} else {
			c.Logger().Infof("Order (%v) Pushed Successfully.", result.ID)
		}
	}

	// To response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      result.ID,
		Success: true,
	}

	c.Logger().Infof("{%v} with id is created.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusCreated, jsonSuccessResultId)
}

// GenericEndpointFromMongo godoc
// @Summary get orders list with filter
// @ID get-orders-with-filter-from-mongoDB
// @Produce json
// @Param data body order_api.OrderGetRequest true "order filter data"
// @Success 200 {object} models.JSONSuccessResultData
// @Success 400 {object} pkg.BadRequestError
// @Success 404 {object} pkg.NotFoundError
// @Router /orders/GenericEndpointFromMongo [post]
func (h *OrderHandler) GenericEndpointFromMongo(c echo.Context) error {
	var orderGetRequest order_api.OrderGetRequest

	if err := c.Bind(&orderGetRequest); err != nil {
		badRequestErr := pkg.BadRequestError{
			Message:    fmt.Sprintf("Bad Request. It cannot be binding! %v", err),
			StatusCode: http.StatusBadRequest,
		}
		return badRequestErr
	}

	// Create filter and find options for mongoDB (exact filter,sort,field and match)
	filter, findOptions := h.Service.FromModelConvertToFilter(orderGetRequest)

	// Get request with filter and find options for mongoDB
	orderList, err := h.Service.GetOrdersWithFilter(filter, findOptions)

	if err != nil {
		notFoundErr := pkg.NotFoundError{
			Message:    fmt.Sprintf("NotFoundError. %v", err),
			StatusCode: http.StatusNotFound,
		}
		return notFoundErr
	}

	// Response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(orderList),
		Data:           orderList,
	}

	c.Logger().Info("Orders are successfully listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// GenericEndpointFromElastic godoc
// @Summary get orders list with filter
// @ID get-orders-with-filter-from-elasticsearch
// @Produce json
// @Param data body order_api.OrderGetRequest true "order filter data"
// @Success 200 {object} models.JSONSuccessResultData
// @Success 400 {object} pkg.BadRequestError
// @Success 404 {object} pkg.NotFoundError
// @Router /orders/GenericEndpointFromElastic [post]
func (h *OrderHandler) GenericEndpointFromElastic(c echo.Context) error {
	var orderGetRequest order_api.OrderGetRequest

	if err := c.Bind(&orderGetRequest); err != nil {
		badRequestErr := pkg.BadRequestError{
			Message:    fmt.Sprintf("Bad Request. It cannot be binding! %v", err),
			StatusCode: http.StatusBadRequest,
		}
		return badRequestErr
	}

	// Create filter and find options (exact filter,sort,field and match)
	elasticQuery := h.ElasticService.FromModelConvertToElasticQuery(orderGetRequest)

	// Get orders from elasticsearch
	orderList, err := h.ElasticService.GetFromElasticsearch(elasticQuery)

	if err != nil {
		internalServerErr := pkg.InternalServerError{
			Message:    fmt.Sprintf("InternalServerError. %v", err),
			StatusCode: http.StatusInternalServerError,
		}
		return internalServerErr
	}

	// Response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(orderList),
		Data:           orderList,
	}

	c.Logger().Info("Orders are successfully listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// UpdateOrder godoc
// @Summary update an item to the order list
// @ID update-order
// @Produce json
// @Param data body order_api.OrderUpdateRequest true "order data"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 400 {object} pkg.BadRequestError
// @Success 500 {object} pkg.InternalServerError
// @Router /orders [put]
func (h *OrderHandler) UpdateOrder(c echo.Context) error {
	// Get order model from middleware because we bind it within middleware
	orderUpdateRequest := c.Get("order").(*order_api.OrderUpdateRequest)

	// We parse the data as json into the struct
	/*	var orderUpdateRequest order_api.OrderUpdateRequest
		if err := c.Bind(&orderUpdateRequest); err != nil {
			c.Logger().Errorf("Bad Request! %v", err)
			return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
				Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
			})
		}*/

	// Validate user input using the validator instance
	if err := h.Validator.Struct(orderUpdateRequest); err != nil {
		badRequestErr := pkg.BadRequestError{
			Message:    fmt.Sprintf("Bad Request. Please put valid order model ! %v", err),
			StatusCode: http.StatusBadRequest,
		}
		return badRequestErr
	}

	// Check order using with service
	if _, err := h.Service.GetOrderById(orderUpdateRequest.ID); err != nil {
		notFoundErr := pkg.NotFoundError{
			Message:    fmt.Sprintf("Not found exception: {%v} with id not found!", orderUpdateRequest.ID),
			StatusCode: http.StatusNotFound,
		}
		return notFoundErr
	}

	// Check user with http.Client
	user, err := h.Service.GetUser(orderUpdateRequest.UserId, h.Config.HttpClient.UserAPI)
	if err != nil {
		notFoundErr := pkg.NotFoundError{
			Message:    fmt.Sprintf("User with id (%v) cannot find!", orderUpdateRequest.UserId),
			StatusCode: http.StatusNotFound,
		}
		return notFoundErr
	}

	// Address check
	regularAddressCheck := false
	invoiceAddressCheck := false
	for _, regularAddress := range user.Addresses {
		if regularAddress.ID == orderUpdateRequest.Address {
			regularAddressCheck = true
		}
		if regularAddress.ID == orderUpdateRequest.InvoiceAddress {
			invoiceAddressCheck = true
		}
	}

	if regularAddressCheck == false || invoiceAddressCheck == false {
		notFoundErr := pkg.NotFoundError{
			Message:    fmt.Sprintf("Not Found Exception: Address not found. Before order processing please put correct address id."),
			StatusCode: http.StatusNotFound,
		}
		return notFoundErr
	}

	// Mapping => we can use automapper, but it will cause performance loss.
	var order models.Order
	order.ID = orderUpdateRequest.ID
	order.UserId = orderUpdateRequest.UserId
	order.Status = orderUpdateRequest.Status
	for _, regularAddress := range user.Addresses {
		if regularAddress.ID == orderUpdateRequest.Address {
			order.Address.ID = regularAddress.ID
			order.Address.Address = regularAddress.Address
			order.Address.City = regularAddress.City
			order.Address.District = regularAddress.District
			order.Address.Type = regularAddress.Type
			order.Address.Default = regularAddress.Default
		}
	}
	for _, invoiceAddress := range user.Addresses {
		if invoiceAddress.ID == orderUpdateRequest.InvoiceAddress {
			order.InvoiceAddress.ID = invoiceAddress.ID
			order.InvoiceAddress.Address = invoiceAddress.Address
			order.InvoiceAddress.City = invoiceAddress.City
			order.InvoiceAddress.District = invoiceAddress.District
			order.InvoiceAddress.Type = invoiceAddress.Type
			order.InvoiceAddress.Default = invoiceAddress.Default
		}
	}
	order.Product = []struct {
		Name     string  `json:"name" bson:"name"`
		Quantity int     `json:"quantity" bson:"quantity"`
		Price    float64 `json:"price" bson:"price"`
	}(orderUpdateRequest.Product)

	// Service => Update
	result, err := h.Service.Update(order)

	if err != nil || result == false {
		internalServerError := pkg.InternalServerError{
			Message:    fmt.Sprintf("StatusInternalServerError: {%v} ", err),
			StatusCode: http.StatusInternalServerError,
		}
		return internalServerError
	}

	// => SEND MESSAGE (OrderID)
	var orderKafka order_api.OrderResponseForElastic
	orderKafka.OrderID = order.ID
	orderKafka.Status = "Updated"

	resultJson, errJson := json.Marshal(orderKafka)
	if errJson != nil {
		c.Logger().Errorf("Something went wrong convert to byte: %v", err)
	} else {
		err = h.Producer.SendToKafkaWithMessage(resultJson, h.Config.Kafka.TopicName["OrderID"])
		if err != nil {
			c.Logger().Errorf("Something went wrong cannot pushed: %v", err)
		} else {
			c.Logger().Infof("Order (%v) Pushed Successfully.", order.ID)
		}
	}

	// To response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      order.ID,
		Success: result,
	}

	c.Logger().Infof("{%v} with id is updated.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}

// DeleteOrder godoc
// @Summary delete an order item by ID
// @ID delete-order-by-id
// @Produce json
// @Param id path string true "order ID"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 404 {object} pkg.NotFoundError
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	query := c.Param("id")

	result, err := h.Service.Delete(query)

	if err != nil || result == false {
		notFoundErr := pkg.NotFoundError{
			Message:    fmt.Sprintf("Not found exception: {%v} with id not found!", query),
			StatusCode: http.StatusNotFound,
		}
		return notFoundErr
	}

	// => SEND MESSAGE (OrderID)
	var orderKafka order_api.OrderResponseForElastic
	orderKafka.OrderID = query
	orderKafka.Status = "Deleted"

	resultJson, errJson := json.Marshal(orderKafka)
	if errJson != nil {
		c.Logger().Errorf("Something went wrong convert to byte: %v", err)
	} else {
		err = h.Producer.SendToKafkaWithMessage(resultJson, h.Config.Kafka.TopicName["OrderID"])
		if err != nil {
			c.Logger().Errorf("Something went wrong cannot pushed: %v", err)
		} else {
			c.Logger().Infof("Order (%v) Pushed Successfully.", query)
		}
	}

	// To response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      query,
		Success: result,
	}

	c.Logger().Infof("{%v} with id is deleted.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}
