package order_api

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OrderService struct {
	OrderRepository repository.IOrderRepository
	UserRepository  repository.IUserRepository
}

func NewOrderService(orderRepository repository.IOrderRepository, userRepository repository.IUserRepository) *OrderService {
	orderService := &OrderService{
		OrderRepository: orderRepository,
		UserRepository:  userRepository,
	}
	return orderService
}

type IOrderService interface {
	GetAll() ([]models.Order, error)
	GetOrderById(id string) (models.Order, error)
	Insert(order models.Order) (models.Order, error)
	Update(user models.Order) (bool, error)
	Delete(id string) (bool, error)
}

func (b OrderService) GetAll() ([]models.Order, error) {
	result, err := b.OrderRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b OrderService) GetOrderById(id string) (models.Order, error) {

	result, err := b.OrderRepository.GetOrderById(id)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b OrderService) Insert(order models.Order) (models.Order, error) {

	// to create id and created date value
	order.ID = uuid.New().String()
	order.CreatedDate = primitive.NewDateTimeFromTime(time.Now())

	result, err := b.OrderRepository.Insert(order)

	if err != nil || result == false {
		return order, err
	}

	addUserOrders, err := b.UserRepository.UpdateOrder(order.ID, order.UserId)

	if err != nil || addUserOrders == false {
		return order, err
	}

	return order, nil
}

func (b OrderService) Update(order models.Order) (bool, error) {
	// to create updated date value
	order.UpdatedDate = primitive.NewDateTimeFromTime(time.Now())

	result, err := b.OrderRepository.Update(order)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func (b OrderService) Delete(id string) (bool, error) {
	// to find userId => if we use jwt, we can find token payload
	order, err := b.OrderRepository.GetOrderById(id)

	result, err := b.OrderRepository.Delete(id)

	if err != nil || result == false {
		return false, err
	}

	updateUser, err := b.UserRepository.DeleteOrder(id, order.UserId)

	if err != nil || updateUser == false {
		return false, err
	}

	return true, nil
}
