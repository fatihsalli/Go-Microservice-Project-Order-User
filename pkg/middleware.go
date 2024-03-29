package pkg

import (
	"OrderUserProject/internal/apps/order-api"
	"fmt"
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

// CheckOrderStatus => Middleware: Status Check using Reflection for Update and Post method (Learning Reflection!)
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
			// To work with reflection we use pointer
			order = &order_api.OrderCreateRequest{}
			orderType = reflect.TypeOf(&order_api.OrderCreateRequest{})
		} else if method == http.MethodPut {
			// To work with reflection we use pointer
			order = &order_api.OrderUpdateRequest{}
			orderType = reflect.TypeOf(&order_api.OrderUpdateRequest{})
		}

		if err := c.Bind(order); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, CustomError{
				Message:    "Invalid request payload!",
				StatusCode: http.StatusBadRequest,
			})
		}

		// It is not necessary, it made to learn reflection
		// Checking if it is assignable to the order model
		if reflect.TypeOf(order).AssignableTo(orderType) {

			statusList := []string{
				"Shipped",
				"Not Shipped",
				"Delivered",
				"Not Delivered",
				"Canceled",
				"Closed",
			}

			// If "order" represents the "pointer" then "Elem()" is used to reach the target value of "pointer"
			orderValue := reflect.ValueOf(order).Elem()
			orderStatusValue := orderValue.FieldByName("Status").String()
			for _, status := range statusList {
				if orderStatusValue == status {
					// To reach value of order we can c.Set and c.Get.Otherwise we cannot bind context twice
					c.Set("order", order)
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusBadRequest, CustomError{
				Message:    "Please write a valid status value!",
				StatusCode: http.StatusBadRequest,
			})
		}

		return echo.NewHTTPError(http.StatusBadRequest, CustomError{
			Message:    "Something wrong! Type of model inconsistent.",
			StatusCode: http.StatusBadRequest,
		})
	}
}

// CustomErrorMiddleware => Middleware: To control log level and return error model
func CustomErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			// To find panic error
			if r := recover(); r != nil {
				err := fmt.Errorf("!Panic occurred: %v", r)
				c.Response().Status = 500
				logMessage := fmt.Sprintf("Error: %v | Request: %v | Response: %v",
					err, c.Request(), c.Response())
				c.Logger().Error(logMessage)
				_ = c.JSON(http.StatusInternalServerError, CustomError{
					Message:    "Oops. Something wrong!",
					StatusCode: http.StatusInternalServerError,
				})
			}
		}()

		// Call next middleware
		err := next(c)

		// CustomError check
		if customError, ok := err.(CustomError); ok {
			if customError.StatusCode == 400 {
				c.Logger().Info(customError.Message)
				return c.JSON(http.StatusBadRequest, CustomError{
					Message: customError.Message,
				})
			}

			if customError.StatusCode == 404 {
				c.Logger().Info(customError.Message)
				return c.JSON(http.StatusNotFound, CustomError{
					Message: customError.Message,
				})
			}

			if customError.StatusCode >= 400 && customError.StatusCode <= 499 {
				c.Logger().Info(customError.Message)
				return c.JSON(customError.StatusCode, CustomError{
					Message: customError.Message,
				})
			}

			if customError.StatusCode >= 500 {
				c.Response().Status = customError.StatusCode
				logMessage := fmt.Sprintf("Error: %v | Request: %v | Response: %v",
					customError.Message, c.Request(), c.Response())
				c.Logger().Error(logMessage)
				return c.JSON(http.StatusInternalServerError, CustomError{
					Message: "Oops. Something wrong!",
				})
			}
		}

		// Handle error if occurred in subsequent middleware or handler
		if err != nil {
			c.Response().Status = 500
			logMessage := fmt.Sprintf("Error: %v | Request: %v | Response: %v",
				err, c.Request(), c.Response())
			c.Logger().Error(logMessage)
			return c.JSON(http.StatusInternalServerError, CustomError{
				Message: "Oops. Something wrong!",
			})
		}

		return err
	}
}
