package handler

import (
	order_api "OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/kafka"
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type OrderHandler struct {
	Service order_api.IOrderService
}

func NewOrderHandler(e *echo.Echo, service order_api.IOrderService) *OrderHandler {
	router := e.Group("api/orders")
	b := &OrderHandler{Service: service}

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
func (h OrderHandler) GetAllOrders(c echo.Context) error {
	orderList, err := h.Service.GetAll()

	if err != nil {
		log.Printf("StatusInternalServerError: %v", err)
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
		orderResponse.Total = order.Total

		ordersResponse = append(ordersResponse, orderResponse)
	}

	// to response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(ordersResponse),
		Data:           ordersResponse,
	}

	log.Print("All books are listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// GetOrderById godoc
// @Summary get a order item by ID
// @ID get-order-by-id
// @Produce json
// @Param id path string true "order ID"
// @Success 200 {object} models.JSONSuccessResultData
// @Success 404 {object} pkg.NotFoundError
// @Success 500 {object} pkg.InternalServerError
// @Router /orders/{id} [get]
func (h OrderHandler) GetOrderById(c echo.Context) error {
	query := c.Param("id")

	order, err := h.Service.GetOrderById(query)

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

	// we can use automapper, but it will cause performance loss.
	var orderResponse order_api.OrderResponse
	orderResponse.ID = order.ID
	orderResponse.UserId = order.UserId
	orderResponse.Product = order.Product
	orderResponse.Total = order.Total

	// to response success result data => single one
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: 1,
		Data:           orderResponse,
	}

	log.Printf("{%v} with id is listed.", orderResponse.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// CreateOrder godoc
// @Summary add a new item to the order list
// @ID create-order
// @Produce json
// @Param data body order_api.OrderCreateRequest true "order data"
// @Success 201 {object} models.JSONSuccessResultId
// @Success 400 {object} pkg.BadRequestError
// @Success 500 {object} pkg.InternalServerError
// @Router /orders [post]
func (h OrderHandler) CreateOrder(c echo.Context) error {

	var orderRequest order_api.OrderCreateRequest

	// We parse the data as json into the struct
	if err := c.Bind(&orderRequest); err != nil {
		log.Printf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	var order models.Order

	// we can use automapper, but it will cause performance loss.
	order.UserId = orderRequest.UserId
	order.Product = orderRequest.Product
	order.Total = orderRequest.Total

	result, err := h.Service.Insert(order)

	if err != nil {
		log.Printf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
	}

	// publish event
	// convert body into bytes and send it to kafka
	orderInBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("There was a problem when convert to byte format: %v", err.Error())
	}

	// create topic name
	topic := "order-test-v01"

	// sending data
	err = kafka.SendToKafka(topic, orderInBytes)
	if err != nil {
		log.Printf("There was a problem when sending message: %v", err.Error())
	}
	log.Printf("Order (%v) Pushed Successfully.", result.ID)

	// listening data
	kafka.ListenFromKafka(topic)

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      result.ID,
		Success: true,
	}

	log.Printf("{%v} with id is created.", jsonSuccessResultId.ID)
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
func (h OrderHandler) UpdateOrder(c echo.Context) error {

	var orderUpdateRequest order_api.OrderUpdateRequest

	// we parse the data as json into the struct
	if err := c.Bind(&orderUpdateRequest); err != nil {
		log.Printf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	if _, err := h.Service.GetOrderById(orderUpdateRequest.ID); err != nil {
		log.Printf("Not found exception: {%v} with id not found!", orderUpdateRequest.ID)
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", orderUpdateRequest.ID),
		})
	}

	var order models.Order

	// we can use automapper, but it will cause performance loss.
	order.ID = orderUpdateRequest.ID
	order.Product = orderUpdateRequest.Product
	order.Total = orderUpdateRequest.Total
	order.UserId = orderUpdateRequest.UserId

	result, err := h.Service.Update(order)

	if err != nil || result == false {
		log.Printf("StatusInternalServerError: {%v} ", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      order.ID,
		Success: result,
	}

	log.Printf("{%v} with id is updated.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusOK, jsonSuccessResultId)
}

// DeleteOrder godoc
// @Summary delete a order item by ID
// @ID delete-order-by-id
// @Produce json
// @Param id path string true "order ID"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 404 {object} pkg.NotFoundError
// @Router /orders/{id} [delete]
func (h OrderHandler) DeleteOrder(c echo.Context) error {
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
