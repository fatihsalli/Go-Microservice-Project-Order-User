package order_api

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	GetUser(userId string, userURL string) (UserResponse, error)
	FromModelConvertToFilter(req OrderGetRequest) (bson.M, *options.FindOptions)
	GetOrdersWithFilter(filter bson.M, opt *options.FindOptions) ([]interface{}, error)
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
		return models.Order{}, err
	}

	return result, nil
}

func (b *OrderService) Insert(order models.Order) (models.Order, error) {
	// Create id and created date value
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()
	// We don't want to set null, so we put CreatedAt value.
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
	// Create updated date value
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

func (b *OrderService) GetUser(userId string, userURL string) (UserResponse, error) {
	// => HTTP.CLIENT FIND USER
	// Create a new HTTP client with a timeout (to check user)
	client := http.Client{
		Timeout: time.Second * 20,
	}

	// Send a GET request to the User service to retrieve user information
	respUser, err := client.Get(userURL + "/" + userId)
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

func (b *OrderService) FromModelConvertToFilter(req OrderGetRequest) (bson.M, *options.FindOptions) {

	// Create a filter based on the exact filters and matches provided in the request
	filter := bson.M{}

	// Add exact filter criteria to filter if provided
	if len(req.ExactFilters) > 0 {
		for key, values := range req.ExactFilters {
			filter[key] = bson.M{"$in": values}
		}
	}

	// Add match criteria to filter if provided
	if len(req.Match) > 0 {
		match := bson.M{}
		for key, value := range req.Match {
			match[key] = value
		}
		filter = bson.M{
			"$and": []bson.M{
				filter,
				match,
			},
		}
	}

	// Create options for the find operation, including the requested fields and sort order
	findOptions := options.Find()

	// Add projection criteria to find options if provided
	if len(req.Fields) > 0 {
		projection := bson.M{}
		findOptions.SetProjection(projection)
		for _, field := range req.Fields {
			projection[field] = 1
		}
	}

	// Add sort criteria to find options if provided
	if len(req.Sort) > 0 {
		sort := bson.M{}
		for key, value := range req.Sort {
			sort[key] = value
		}
		findOptions.SetSort(sort)
	}

	return filter, findOptions
}

func (b *OrderService) GetOrdersWithFilter(filter bson.M, opt *options.FindOptions) ([]interface{}, error) {
	result, err := b.OrderRepository.GetOrdersWithFilter(filter, opt)

	if err != nil {
		return nil, err
	}

	return result, nil
}
