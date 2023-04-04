package order_api

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"time"
)

type OrderService struct {
	OrderRepository *repository.OrderRepository
}

func NewOrderService(orderRepository *repository.OrderRepository) *OrderService {
	orderService := &OrderService{
		OrderRepository: orderRepository,
	}
	return orderService
}

type IOrderService interface {
	GetAll() ([]models.Order, error)
	GetOrderById(id string) (models.Order, error)
	Insert(order models.Order) (models.Order, error)
	Update(user models.Order) (bool, error)
	Delete(id string) (bool, error)
	GetUser(userId string) (UserResponse, error)
}

func (b *OrderService) GetAll() ([]models.Order, error) {
	result, err := b.OrderRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *OrderService) GetOrderById(id string) (models.Order, error) {

	result, err := b.OrderRepository.GetOrderById(id)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b *OrderService) Insert(order models.Order) (models.Order, error) {
	// to create id and created date value
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()
	// we don't want to set null, so we put CreatedAt value.
	order.UpdatedAt = order.CreatedAt

	var total float64
	for _, product := range order.Product {
		total = product.Price * float64(product.Quantity)
		order.Total += total
	}

	result, err := b.OrderRepository.Insert(order)

	if err != nil || result == false {
		return models.Order{}, err
	}

	return order, nil
}

func (b *OrderService) Update(order models.Order) (bool, error) {
	// to create updated date value
	order.UpdatedAt = time.Now()

	var total float64
	for _, product := range order.Product {
		total = product.Price * float64(product.Quantity)
		order.Total += total
	}

	result, err := b.OrderRepository.Update(order)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func (b *OrderService) Delete(id string) (bool, error) {
	result, err := b.OrderRepository.Delete(id)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func (b *OrderService) GetUser(userId string) (UserResponse, error) {
	// => HTTP.CLIENT FIND USER
	// Create a new HTTP client with a timeout (to check user)
	client := http.Client{
		Timeout: time.Second * 20,
	}

	// Send a GET request to the User service to retrieve user information
	respUser, err := client.Get("http://localhost:8012/api/users" + "/" + userId)
	if err != nil || respUser.StatusCode != http.StatusOK {
		return UserResponse{}, errors.New("user cannot find")
	}
	defer func() {
		if err := respUser.Body.Close(); err != nil {
			log.Errorf("Something went wrong: %v", err)
		}
	}()

	// Read the response body
	respUserBody, err := io.ReadAll(respUser.Body)
	if err != nil {
		return UserResponse{}, err
	}

	// Unmarshal the response body into an Order struct
	var userResponse UserResponse
	err = json.Unmarshal(respUserBody, &userResponse)
	if err != nil {
		return UserResponse{}, err
	}

	return userResponse, nil
}
