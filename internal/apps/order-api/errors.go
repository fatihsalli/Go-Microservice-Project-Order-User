package order_api

import (
	"OrderUserProject/pkg"
	"fmt"
	"net/http"
)

var OrderNotFound = pkg.NotFoundError{
	Message:    fmt.Sprintf("Not found exception: {%v} with id not found!", "query"),
	StatusCode: http.StatusNotFound,
}

// TODO: Direkt olarak CustomError ile işi çözebiliriz diğer.

func OrderNotFoundMethod(message string) error {
	var OrderNotFound = pkg.NotFoundError{
		Message:    fmt.Sprintf("Not found exception: {%v} with id not found!", message),
		StatusCode: http.StatusNotFound,
	}
	return OrderNotFound
}
