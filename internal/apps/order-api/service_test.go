package order_api

import (
	"OrderUserProject/internal/models"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

var ordersList = []models.Order{
	{
		ID:     "2b45ac31-6906-4e1e-82db-d9bcdbdb2143",
		UserId: "fcd20a19-6171-4737-a2ed-23e293cae7b5",
		Status: "Shipped",
		Address: models.Address{
			ID:       "130beada-8339-4ee6-a754-725f43b8da98",
			Address:  "Levent",
			City:     "İstanbul",
			District: "Beşiktaş",
			Type: []string{
				"Regular",
			},
			Default: struct {
				IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
				IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
			}{
				IsDefaultInvoiceAddress: false,
				IsDefaultRegularAddress: true,
			},
		},
		InvoiceAddress: models.Address{
			ID:       "4e4ed986-6d06-4bcf-b577-3d68a72949a7",
			Address:  "Suadiye",
			City:     "İstanbul",
			District: "Kadıköy",
			Type: []string{
				"Invoice",
			},
			Default: struct {
				IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
				IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
			}{
				IsDefaultInvoiceAddress: true,
				IsDefaultRegularAddress: false,
			},
		},
		Product: []struct {
			Name     string  `json:"name" bson:"name"`
			Quantity int     `json:"quantity" bson:"quantity"`
			Price    float64 `json:"price" bson:"price"`
		}{
			{
				Name:     "Asus Notebook",
				Quantity: 1,
				Price:    20000.0,
			},
			{
				Name:     "Airpods",
				Quantity: 1,
				Price:    4000.0,
			},
		},
		Total:     24000.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:     "41840818-6f62-4378-a82a-e6badf225bcc",
		UserId: "72df3086-564a-470d-a68c-82476e988a54",
		Status: "Delivered",
		Address: models.Address{
			ID:       "56f97563-6941-4793-b248-9051ed6aa256",
			Address:  "Narlıdere",
			City:     "İzmir",
			District: "Narlıdere",
			Type: []string{
				"Regular", "Invoice",
			},
			Default: struct {
				IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
				IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
			}{
				IsDefaultInvoiceAddress: true,
				IsDefaultRegularAddress: true,
			},
		},
		InvoiceAddress: models.Address{
			ID:       "56f97563-6941-4793-b248-9051ed6aa256",
			Address:  "Narlıdere",
			City:     "İzmir",
			District: "Narlıdere",
			Type: []string{
				"Regular", "Invoice",
			},
			Default: struct {
				IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
				IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
			}{
				IsDefaultInvoiceAddress: true,
				IsDefaultRegularAddress: true,
			},
		},
		Product: []struct {
			Name     string  `json:"name" bson:"name"`
			Quantity int     `json:"quantity" bson:"quantity"`
			Price    float64 `json:"price" bson:"price"`
		}{
			{
				Name:     "Iphone 12",
				Quantity: 1,
				Price:    24000.0,
			},
		},
		Total:     24000.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

// MockOrderRepository is a mock implementation of IOrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) GetAll() ([]models.Order, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	// []models.Order => 1.Return model || error => 2.Return model
	return args.Get(0).([]models.Order), nil
}

func (m *MockOrderRepository) GetOrderById(id string) (models.Order, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return models.Order{}, args.Error(1)
	}
	return args.Get(0).(models.Order), nil
}

func (m *MockOrderRepository) Insert(order models.Order) (bool, error) {
	args := m.Called(order)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return true, nil
}

func (m *MockOrderRepository) Update(order models.Order) (bool, error) {
	args := m.Called(order)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return true, nil
}

func (m *MockOrderRepository) Delete(id string) (bool, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return true, nil
}

func (m *MockOrderRepository) GetOrdersWithFilter(filter bson.M, opt *options.FindOptions) ([]interface{}, error) {
	args := m.Called(filter, opt)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), nil
}

func TestOrderService_GetAll_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockOrderRepository)

	mockRepo.On("GetAll").Return(ordersList, nil)

	// Create an instance of OrderService with the mock repository
	orderService := NewOrderService(mockRepo)

	// Call the GetAll method
	orders, err := orderService.GetAll()

	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, ordersList, orders)

	// Verify that the mock method was called
	mockRepo.AssertCalled(t, "GetAll")
}

func TestOrderService_GetOrderById_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockOrderRepository)

	// Define the expected result
	mockRepo.On("GetOrderById", "2b45ac31-6906-4e1e-82db-d9bcdbdb2143").Return(ordersList[0], nil)

	// Create an instance of OrderService with the mock repository
	orderService := NewOrderService(mockRepo)

	// Call the GetOrderById method
	order, err := orderService.GetOrderById("2b45ac31-6906-4e1e-82db-d9bcdbdb2143")

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, ordersList[0], order)

	// Verify that the mock method was called
	mockRepo.AssertCalled(t, "GetOrderById", "2b45ac31-6906-4e1e-82db-d9bcdbdb2143")
}

func TestOrderService_GetOrderById_NotFoundFail(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockOrderRepository)

	expectedError := errors.New("not found error")

	// Define the expected result
	mockRepo.On("GetOrderById", "2b45ac31-6906-4e1e-82db-d9bcdbdb2141").
		Return(models.Order{}, expectedError)

	// Create an instance of OrderService with the mock repository
	orderService := NewOrderService(mockRepo)

	// Call the GetOrderById method
	order, err := orderService.GetOrderById("2b45ac31-6906-4e1e-82db-d9bcdbdb2141")

	// Check error
	if !errors.Is(err, expectedError) {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}

	// Check nil order model
	if order.ID != "" {
		t.Error("Expected empty order model, but got a non-empty model!")
	}
}

func TestOrderService_Insert_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockOrderRepository)

	// Define the input and expected result
	order := models.Order{
		ID:     "",
		UserId: "4ae3b4a1-1cab-460e-a1bf-0d3a73f2787f",
		Status: "Not Shipped",
		Address: models.Address{
			ID:       "ddf1e162-d438-4167-b131-bbc7b767fa9d",
			Address:  "Levent",
			City:     "İstanbul",
			District: "Beşiktaş",
			Type: []string{
				"Regular",
			},
			Default: struct {
				IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
				IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
			}{
				IsDefaultInvoiceAddress: false,
				IsDefaultRegularAddress: true,
			},
		},
		InvoiceAddress: models.Address{
			ID:       "a7d8b4ae-cb77-4bb3-a8e8-accff9075bf3",
			Address:  "Bulgurlu",
			City:     "İstanbul",
			District: "Üsküdar",
			Type: []string{
				"Invoice",
			},
			Default: struct {
				IsDefaultInvoiceAddress bool `json:"isDefaultInvoiceAddress" bson:"isDefaultInvoiceAddress"`
				IsDefaultRegularAddress bool `json:"isDefaultRegularAddress" bson:"isDefaultRegularAddress"`
			}{
				IsDefaultInvoiceAddress: true,
				IsDefaultRegularAddress: false,
			},
		},
		Product: []struct {
			Name     string  `json:"name" bson:"name"`
			Quantity int     `json:"quantity" bson:"quantity"`
			Price    float64 `json:"price" bson:"price"`
		}{
			{
				Name:     "LG Smart Tv",
				Quantity: 1,
				Price:    20000.0,
			},
			{
				Name:     "Bosch Filter Coffee Machine",
				Quantity: 1,
				Price:    2500.0,
			},
		},
		Total:     0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	// We don't know exact order model because in service we have changed order model
	mockRepo.On("Insert", mock.AnythingOfType("models.Order")).Return(true, nil)

	// Create an instance of OrderService with the mock repository
	orderService := NewOrderService(mockRepo)

	// Call the Insert method
	result, err := orderService.Insert(order)

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, order.UserId, result.UserId)
	assert.Equal(t, order.Status, result.Status)
	assert.Equal(t, float64(22500), result.Total)
	assert.Equal(t, order.Product, result.Product)
	assert.Equal(t, order.Address, result.Address)
	assert.Equal(t, order.InvoiceAddress, result.InvoiceAddress)

	// We don't know exact order model because in service we have changed order model
	mockRepo.AssertCalled(t, "Insert", mock.AnythingOfType("models.Order"))
}

func TestOrderService_Update_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockOrderRepository)

	// Define the input and expected result
	order := ordersList[0]

	// We don't know exact order model because in service we have changed order model
	mockRepo.On("Update", mock.AnythingOfType("models.Order")).Return(true, nil)

	// Create an instance of OrderService with the mock repository
	orderService := NewOrderService(mockRepo)

	// Call the Insert method
	result, err := orderService.Update(order)

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, true, result)

	// We don't know exact order model because in service we have changed order model
	mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("models.Order"))
}

func TestOrderService_Delete_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockOrderRepository)

	id := "2b45ac31-6906-4e1e-82db-d9bcdbdb2143"

	// We don't know exact order model because in service we have changed order model
	mockRepo.On("Delete", id).Return(true, nil)

	// Create an instance of OrderService with the mock repository
	orderService := NewOrderService(mockRepo)

	// Call the Insert method
	result, err := orderService.Delete(id)

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, true, result)

	// We don't know exact order model because in service we have changed order model
	mockRepo.AssertCalled(t, "Delete", id)
}

func TestOrderService_GetOrdersWithFilter_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockOrderRepository)

	orderRequest := OrderGetRequest{
		ExactFilters: map[string][]interface{}{
			"address.city": {"İzmir"},
		},
		Fields: []string{
			"userId", "status", "total",
		},
		Match: []struct {
			MatchField string      `json:"match_field"`
			Parameter  string      `json:"parameter"`
			Value      interface{} `json:"value"`
		}{{
			MatchField: "address.address",
			Parameter:  "eq",
			Value:      "Narlıdere",
		}},
		Sort: map[string]int{"total": -1},
	}

	// Create an instance of OrderService with the mock repository
	orderService := NewOrderService(mockRepo)

	selectedOrder := ordersList[0]
	filteredOrder := struct {
		id     string
		userId string
		status string
		total  float64
	}{
		id:     selectedOrder.ID,
		userId: selectedOrder.UserId,
		status: selectedOrder.Status,
		total:  selectedOrder.Total,
	}

	orderAsInterface := interface{}(filteredOrder)

	var orders []interface{}
	orders = append(orders, orderAsInterface)

	filter, opt := orderService.FromModelConvertToFilter(orderRequest)

	// We don't know exact order model because in service we have changed order model
	mockRepo.On("GetOrdersWithFilter", filter, opt).Return(orders, nil)

	orderServiceLast := NewOrderService(mockRepo)

	// Call the Insert method
	result, err := orderServiceLast.GetOrdersWithFilter(filter, opt)

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, orders, result)

	// We don't know exact order model because in service we have changed order model
	mockRepo.AssertCalled(t, "GetOrdersWithFilter", filter, opt)
}
