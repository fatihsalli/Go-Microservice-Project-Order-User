package handler

import (
	order_api "OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/elastic"
	"OrderUserProject/internal/kafka"
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"time"
)

type AggregatorHandler struct {
}

var ClientBaseUrl = map[string]string{
	"order": "http://localhost:8011/api/orders",
	"user":  "http://localhost:8012/api/users",
}

func NewGatewayHandler(e *echo.Echo) *AggregatorHandler {
	router := e.Group("api/")

	b := &AggregatorHandler{}

	//Routes
	router.POST("CreateOrder", b.CreateOrder)

	return b
}

// CreateOrder godoc
// @Summary add a new item to the order list
// @ID create-order
// @Produce json
// @Param data body order_api.OrderCreateRequest true "order data"
// @Success 201 {object} models.JSONSuccessResultId
// @Success 400 {object} pkg.BadRequestError
// @Success 500 {object} pkg.InternalServerError
// @Router /CreateOrder [post]
func (h AggregatorHandler) CreateOrder(c echo.Context) error {

	var orderRequest order_api.OrderCreateRequest

	// We parse the data as json into the struct
	if err := c.Bind(&orderRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * 20,
	}

	// Send a GET request to the User service to retrieve user information
	respUser, err := client.Get(ClientBaseUrl["user"] + "/" + orderRequest.UserId)
	if err != nil || respUser.StatusCode != http.StatusOK {
		c.Logger().Errorf("User with id {%v} cannot find!", orderRequest.UserId)
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("User with id {%v} cannot find!", orderRequest.UserId),
		})
	}
	defer func() {
		if err := respUser.Body.Close(); err != nil {
			c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		}
	}()

	// Convert the payload to JSON bytes
	orderReqBytes, err := json.Marshal(orderRequest)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Order cannot convert to byte format.",
		})
	}

	// Create a new request with the JSON payload
	req, err := http.NewRequest("POST", ClientBaseUrl["order"], bytes.NewBuffer(orderReqBytes))
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot create a new request with the JSON Payload.",
		})
	}

	// Set the request header
	req.Header.Set("Content-Type", "application/json")

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot send the request. Please check the order service.",
		})
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot create a new order. Somethings happen. Please check the logs.",
		})
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot read from response body. Please check the logs.",
		})
	}

	// Unmarshal the response body into an Order struct
	var data models.JSONSuccessResultId
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot convert to JSON format. Please check the logs.",
		})
	}

	// Listen the event with Kafka
	// create topic name
	topic := "order-create-v01"
	go kafka.ListenFromKafka(topic)

	go elastic.ReadFromElastic()

	c.Logger().Infof("{%v} with id is successfully created.", data.ID)
	return c.JSON(http.StatusCreated, data)
}
