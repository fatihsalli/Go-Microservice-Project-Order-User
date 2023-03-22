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
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"time"
)

type OrderElasticHandler struct {
	Service order_elastic.IOrderElasticService
}

var ClientBaseUrl = map[string]string{
	"order":         "http://localhost:8011/api/orders",
	"user":          "http://localhost:8012/api/users",
	"order-elastic": "http://localhost:8013/api/orders-elastic",
}

func NewOrderElasticHandler(e *echo.Echo, service order_elastic.IOrderElasticService) *OrderElasticHandler {
	router := e.Group("api/orders-elastic")
	b := &OrderElasticHandler{Service: service}

	//Routes
	router.GET("", b.CreateOrderElastic)

	return b
}

func (h OrderElasticHandler) CreateOrderElastic(c echo.Context) error {
	// => RECEIVE MESSAGE
	// create topic name
	topic := "orderID-created-v01"

	result := kafka.ListenFromKafka(topic)

	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * 20,
	}

	// Send a GET request to the User service to retrieve user information
	respOrder, err := client.Get(ClientBaseUrl["order"] + "/" + string(result))
	if err != nil || respOrder.StatusCode != http.StatusOK {
		c.Logger().Errorf("User with id {%v} cannot find!", string(result))
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("User with id {%v} cannot find!", string(result)),
		})
	}
	defer func() {
		if err := respOrder.Body.Close(); err != nil {
			c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		}
	}()

	// Read the response body
	respOrderBody, err := io.ReadAll(respOrder.Body)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot read from response body. Please check the logs.",
		})
	}

	// Unmarshal the response body into an Order struct
	var data order_api.OrderResponse
	err = json.Unmarshal(respOrderBody, &data)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot convert to JSON format. Please check the logs.",
		})
	}

	var order models.Order

	// we can use automapper, but it will cause performance loss.
	order.UserId = data.UserId
	order.Status = data.Status
	// mapping from AddressResponse to Address
	order.Address.Address = data.Address.Address
	order.Address.City = data.Address.City
	order.Address.District = data.Address.District
	order.Address.Type = data.Address.Type
	order.Address.Default = data.Address.Default
	order.InvoiceAddress.Address = data.InvoiceAddress.Address
	order.InvoiceAddress.City = data.InvoiceAddress.City
	order.InvoiceAddress.District = data.InvoiceAddress.District
	order.InvoiceAddress.Type = data.InvoiceAddress.Type
	order.InvoiceAddress.Default = data.InvoiceAddress.Default
	order.Product = data.Product

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Errorf("Error marshalling order:", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot convert to [] byte format. Please check the logs.",
		})
	}

	// => SEND MESSAGE
	// create topic name
	topicOrder := "orderDuplicate-created-v01"
	// sending data
	kafka.SendToKafka(topicOrder, orderJSON)
	c.Logger().Infof("Order (%v) Pushed Successfully.", order.ID)

	go func() {
		err = order_elastic.CreateOrderDuplicate()
		if err != nil {
			log.Errorf("Something went wrong:", err)
		}
	}()

	return c.JSON(http.StatusOK, order.ID)
}
