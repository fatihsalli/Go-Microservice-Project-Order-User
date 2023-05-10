package pkg

import (
	"OrderUserProject/internal/apps/order-api"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neko-neko/echo-logrus/v2/log"
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			reqSize := req.Header.Get(echo.HeaderContentLength)
			if reqSize == "" {
				reqSize = "0"
			}

			log.Infof("%s %s [%v] %s %-7s %s %3d %s %s %13v %s %s",
				id,
				c.RealIP(),
				stop.Format(time.RFC3339),
				req.Host,
				req.Method,
				req.RequestURI,
				res.Status,
				reqSize,
				strconv.FormatInt(res.Size, 10),
				stop.Sub(start).String(),
				req.Referer(),
				req.UserAgent(),
			)
			return err
		}
	}
}

// CheckOrderStatus => Middleware: Status Check using Reflection for Update and Post method
func CheckOrderStatus(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Method
		if method != http.MethodPost && method != http.MethodPut {
			return next(c)
		}

		// To check OrderCreateRequest or OrderUpdateRequest
		var order interface{}
		var orderType reflect.Type

		if method == http.MethodPost {
			order = &order_api.OrderCreateRequest{}
			orderType = reflect.TypeOf(&order_api.OrderCreateRequest{})
		} else if method == http.MethodPut {
			order = &order_api.OrderUpdateRequest{}
			orderType = reflect.TypeOf(&order_api.OrderUpdateRequest{})
		}

		if err := c.Bind(order); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, BadRequestError{
				Message: "Invalid request payload!",
			})
		}

		if reflect.TypeOf(order).AssignableTo(orderType) {
			statusList := []string{
				"Shipped",
				"Not Shipped",
				"Delivered",
				"Not Delivered",
				"Canceled",
				"Closed",
			}

			orderValue := reflect.ValueOf(order).Elem()
			orderStatusValue := orderValue.FieldByName("Status").String()
			for _, status := range statusList {
				if orderStatusValue == status {
					return next(c) // Aynı context ile devam et
				}
			}
			return echo.NewHTTPError(http.StatusBadRequest, BadRequestError{
				Message: "Please write a valid status value!",
			})
		}

		return next(c)
	}
}
