package handler

import (
	"OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg"
	"OrderUserProject/pkg/kafka"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type OrderHandler struct {
	Service  *order_api.OrderService
	Producer *kafka.ProducerKafka
}

func NewOrderHandler(e *echo.Echo, service *order_api.OrderService, producer *kafka.ProducerKafka) *OrderHandler {
	router := e.Group("api/orders")
	b := &OrderHandler{Service: service, Producer: producer}

	//Routes
	router.GET("", b.GetAllOrders)
	router.GET("/:id", b.GetOrderById)
	router.POST("", b.CreateOrder)
	router.PUT("", b.UpdateOrder)
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

	// to test GracefulShutdown
	// time.Sleep(5 * time.Second)

	orderList, err := h.Service.GetAll()

	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// we can use automapper, but it will cause performance loss.
	var orderResponse order_api.OrderResponse
	var ordersResponse []order_api.OrderResponse
	for _, order := range orderList {
		orderResponse.ID = order.ID
		orderResponse.UserId = order.UserId
		orderResponse.Product = order.Product
		orderResponse.Address.Address = order.Address.Address
		orderResponse.Address.City = order.Address.City
		orderResponse.Address.District = order.Address.District
		orderResponse.Address.Type = order.Address.Type
		orderResponse.Address.Default = order.Address.Default
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

	// to response success result data
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
	var orderResponse order_api.OrderResponse
	orderResponse.ID = order.ID
	orderResponse.UserId = order.UserId
	orderResponse.Product = order.Product
	orderResponse.Address.Address = order.Address.Address
	orderResponse.Address.City = order.Address.City
	orderResponse.Address.District = order.Address.District
	orderResponse.Address.Type = order.Address.Type
	orderResponse.Address.Default = order.Address.Default
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

	// We parse the data as json into the struct
	var orderRequest order_api.OrderCreateRequest
	if err := c.Bind(&orderRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Check user with http.Client
	err := h.Service.CheckUser(orderRequest.UserId)
	if err != nil {
		c.Logger().Errorf("Not Found Exception: %v", err.Error())
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("User with id (%v) cannot find!", orderRequest.UserId),
		})
	}

	// Mapping => we can use automapper, but it will cause performance loss.
	var order models.Order
	order.UserId = orderRequest.UserId
	order.Status = orderRequest.Status
	order.Address.Address = orderRequest.Address.Address
	order.Address.City = orderRequest.Address.City
	order.Address.District = orderRequest.Address.District
	order.Address.Type = orderRequest.Address.Type
	order.Address.Default = orderRequest.Address.Default
	order.InvoiceAddress.Address = orderRequest.InvoiceAddress.Address
	order.InvoiceAddress.City = orderRequest.InvoiceAddress.City
	order.InvoiceAddress.District = orderRequest.InvoiceAddress.District
	order.InvoiceAddress.Type = orderRequest.InvoiceAddress.Type
	order.InvoiceAddress.Default = orderRequest.InvoiceAddress.Default
	order.Product = orderRequest.Product

	// Service => Insert
	result, err := h.Service.Insert(order)

	// => SEND MESSAGE (OrderID)
	err = h.Producer.SendToKafkaWithMessage([]byte(result.ID))
	if err != nil {
		c.Logger().Errorf("Something went wrong cannot pushed: %v", err)
	} else {
		c.Logger().Infof("Order (%v) Pushed Successfully.", result.ID)
	}

	// TODO : Silinecek kafka-test için yazıldı
	go func() {
		time.Sleep(2000 * time.Millisecond)

		result := kafka.ListenFromKafka("orderID-created-v01")

		c.Logger().Infof("Message is: %v", string(result))
	}()

	// To response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      result.ID,
		Success: true,
	}

	c.Logger().Infof("{%v} with id is created.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusCreated, jsonSuccessResultId)
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

	// We parse the data as json into the struct
	var orderUpdateRequest order_api.OrderUpdateRequest
	if err := c.Bind(&orderUpdateRequest); err != nil {
		c.Logger().Errorf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Check order using with service
	if _, err := h.Service.GetOrderById(orderUpdateRequest.ID); err != nil {
		c.Logger().Errorf("Not found exception: {%v} with id not found!", orderUpdateRequest.ID)
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", orderUpdateRequest.ID),
		})
	}

	// Check user with http.Client
	err := h.Service.CheckUser(orderUpdateRequest.UserId)
	if err != nil {
		c.Logger().Errorf("Not Found Exception: %v", err.Error())
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("User with id (%v) cannot find!", orderUpdateRequest.UserId),
		})
	}

	// Mapping => we can use automapper, but it will cause performance loss.
	var order models.Order
	order.ID = orderUpdateRequest.ID
	order.UserId = orderUpdateRequest.UserId
	order.Status = orderUpdateRequest.Status
	order.Address.Address = orderUpdateRequest.Address.Address
	order.Address.City = orderUpdateRequest.Address.City
	order.Address.District = orderUpdateRequest.Address.District
	order.Address.Type = orderUpdateRequest.Address.Type
	order.Address.Default = orderUpdateRequest.Address.Default
	order.InvoiceAddress.Address = orderUpdateRequest.InvoiceAddress.Address
	order.InvoiceAddress.City = orderUpdateRequest.InvoiceAddress.City
	order.InvoiceAddress.District = orderUpdateRequest.InvoiceAddress.District
	order.InvoiceAddress.Type = orderUpdateRequest.InvoiceAddress.Type
	order.InvoiceAddress.Default = orderUpdateRequest.InvoiceAddress.Default
	order.Product = orderUpdateRequest.Product
	var total float64
	for _, product := range order.Product {
		total = product.Price * float64(product.Quantity)
		order.Total += total
	}

	// Service => Update
	result, err := h.Service.Update(order)

	if err != nil || result == false {
		c.Logger().Errorf("StatusInternalServerError: {%v} ", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
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
		c.Logger().Errorf("Not found exception: {%v} with id not found!", query)
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
		})
	}

	// To response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      query,
		Success: result,
	}

	c.Logger().Infof("{%v} with id is deleted.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}
