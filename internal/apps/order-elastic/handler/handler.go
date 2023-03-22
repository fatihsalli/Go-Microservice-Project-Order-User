package handler

import (
	"OrderUserProject/internal/apps/order-elastic"
	"OrderUserProject/pkg/kafka"
	"github.com/labstack/echo/v4"
	"net/http"
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

	c.Logger().Infof("{%v} with id is created.", string(result))
	return c.JSON(http.StatusOK, string(result))
}
