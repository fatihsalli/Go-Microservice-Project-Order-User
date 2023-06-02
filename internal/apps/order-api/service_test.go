package order_api

import (
	"OrderUserProject/internal/models"
	"OrderUserProject/internal/repository"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

// OrderRepositoryMock, IOrderRepository arayüzünü taklit eden bir mock yapısıdır.
type OrderRepositoryMock struct {
	mock.Mock
}

// GetAll, mock fonksiyonu ile OrderRepository'den tüm siparişleri döndürür.
func (m *OrderRepositoryMock) GetAll() ([]models.Order, error) {
	args := m.Called()
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetOrderById(id string) (models.Order, error) {
	args := m.Called()
	return args.Get(0).(models.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Insert(order models.Order) (bool, error) {
	args := m.Called()
	return args.Get(0).(bool), args.Error(1)
}

func (m *OrderRepositoryMock) Update(user models.Order) (bool, error) {
	args := m.Called()
	return args.Get(0).(bool), args.Error(1)
}

func (m *OrderRepositoryMock) Delete(id string) (bool, error) {
	args := m.Called()
	return args.Get(0).(bool), args.Error(1)
}

func (m *OrderRepositoryMock) GetOrdersWithFilter(filter bson.M, opt *options.FindOptions) ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func TestOrderService_GetAll(t *testing.T) {
	// Mock oluşturulur
	repoMock := &OrderRepositoryMock{}

	// Mock'in GetAll metodu için beklenen dönüş değerleri belirlenir
	expectedOrders := []models.Order{
		{ID: "1", Status: "Shipped"},
		{ID: "2", Status: "Not Shipped"},
	}
	repoMock.On("GetAll").Return(expectedOrders, nil)

	// IOrderRepository arayüzünü uygulayan bir değişken oluşturulur
	var orderRepo repository.IOrderRepository = repoMock

	// OrderService oluşturulur ve mock repository enjekte edilir
	service := NewOrderService(orderRepo)

	// GetAll çağrısı yapılır
	orders, _ := service.GetAll()

	assert.Equal(t, expectedOrders, orders)

	// Mock metodunun çağrıldığından emin olunur
	repoMock.AssertCalled(t, "GetAll")
}
