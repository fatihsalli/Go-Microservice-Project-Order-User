package handler

import (
	order_api "OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg"
	"fmt"
	"github.com/labstack/echo/v4"
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
	router.POST("", b.CreateOrder)

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

	// to response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      result.ID,
		Success: true,
	}

	log.Printf("{%v} with id is created.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusCreated, jsonSuccessResultId)
}
