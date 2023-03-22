package handler

import (
	order_api "OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/internal/models"
	"OrderUserProject/pkg"
	"OrderUserProject/pkg/kafka"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"time"
)

type OrderElasticHandler struct {
	Service order_elastic.IOrderElasticService
}

var ClientBaseUrl = map[string]string{
	"user":  "http://localhost:8012/api/users",
	"order": "http://localhost:8011/api/orders",
}

func NewOrderElasticHandler(e *echo.Echo, service order_elastic.IOrderElasticService) *OrderElasticHandler {
	router := e.Group("api/orders-elastic")
	b := &OrderElasticHandler{Service: service}

	//Routes
	router.GET("", b.CreateOrderElastic)

	return b
}

func (h OrderElasticHandler) CreateOrderElastic(c echo.Context) error {
	topic := "orderID-created-v01"

	orderId := kafka.ListenFromKafka(topic)

	// Create a new HTTP client with a timeout (to check user)
	client := http.Client{
		Timeout: time.Second * 20,
	}

	// Send a GET request to the User service to retrieve user information
	resp, err := client.Get(ClientBaseUrl["order"] + "/" + string(orderId))
	if err != nil || resp.StatusCode != http.StatusOK {
		c.Logger().Errorf("User with id {%v} cannot find!", string(orderId))
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("User with id {%v} cannot find!", string(orderId)),
		})
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		}
	}()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot read from response body. Please check the logs.",
		})
	}

	// Unmarshal the response body into an Order struct
	var orderResponse order_api.OrderResponse
	err = json.Unmarshal(respBody, &orderResponse)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot convert to JSON format. Please check the logs.",
		})
	}

	var order models.Order

	// we can use automapper, but it will cause performance loss.
	order.UserId = orderResponse.UserId
	order.Status = orderResponse.Status
	// mapping from AddressResponse to Address
	order.Address.Address = orderResponse.Address.Address
	order.Address.City = orderResponse.Address.City
	order.Address.District = orderResponse.Address.District
	order.Address.Type = orderResponse.Address.Type
	order.Address.Default = orderResponse.Address.Default
	order.InvoiceAddress.Address = orderResponse.InvoiceAddress.Address
	order.InvoiceAddress.City = orderResponse.InvoiceAddress.City
	order.InvoiceAddress.District = orderResponse.InvoiceAddress.District
	order.InvoiceAddress.Type = orderResponse.InvoiceAddress.Type
	order.InvoiceAddress.Default = orderResponse.InvoiceAddress.Default
	order.Product = orderResponse.Product

	return nil

	//
}
