package handler

import (
	"OrderUserProject/internal/apps/order-api"
	"OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/pkg"
	"OrderUserProject/pkg/kafka"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"time"
)

// TODO: Elasticsearch delete and update with DeleteRequest-IndexRequest

type OrderElasticHandler struct {
	Service order_elastic.IOrderElasticService
}

var ClientBaseUrl = map[string]string{
	"order":         "http://localhost:8011/api/orders",
	"user":          "http://localhost:8012/api/users",
	"order-elastic": "http://localhost:8013/api/orders-elastic",
}

func NewOrderElasticHandler(e *echo.Echo, service order_elastic.OrderElasticService) *OrderElasticHandler {

	router := e.Group("api/orders-elastic")
	b := &OrderElasticHandler{Service: service}

	//Routes
	router.GET("", b.CreateOrderElastic)

	return b
}

func (h OrderElasticHandler) CreateOrderElastic(c echo.Context) error {
	// => RECEIVE MESSAGE (OrderID)
	// create topic name
	topic := "orderID-created-v01"
	result := kafka.ListenFromKafka(topic)

	// => HTTP.CLIENT FIND ORDER
	// Create a new HTTP client with a timeout
	client := http.Client{
		Timeout: time.Second * 20,
	}

	// Send a GET request to the Order service to retrieve order information
	respOrder, err := client.Get(ClientBaseUrl["order"] + "/" + string(result))
	if err != nil || respOrder.StatusCode != http.StatusOK {
		c.Logger().Errorf("Order with id {%v} cannot find!", string(result))
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("Order with id {%v} cannot find!", string(result)),
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
	var orderResponse order_api.OrderResponse
	err = json.Unmarshal(respOrderBody, &orderResponse)
	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot convert to JSON format. Please check the logs.",
		})
	}

	orderJSON, err := json.Marshal(orderResponse)
	if err != nil {
		c.Logger().Errorf("Error marshalling order:", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Cannot convert to [] byte format. Please check the logs.",
		})
	}

	// => SEND MESSAGE (Order Model)
	// create topic name
	topicOrder := "orderDuplicate-created-v01"
	// sending data
	kafka.SendToKafka(topicOrder, orderJSON)
	c.Logger().Infof("Order (%v) Pushed Successfully.", orderResponse.ID)

	// => RECEIVE MESSAGE AND SAVE ON ELASTICSEARCH (Asynchronous)
	go func() {
		order, err := h.Service.ConsumeOrderDuplicate(topicOrder)
		if err != nil {
			c.Logger().Errorf("Something went wrong:", err)
		}

		err = h.Service.SaveOrderToElasticsearch(order)
		if err != nil {
			c.Logger().Errorf("Something went wrong:", err)
		}

		c.Logger().Infof("This id (%v) successfully saved on elastic.", order.ID)
	}()

	return c.JSON(http.StatusOK, orderResponse.ID)
}
