package user_api

import (
	"OrderUserProject/internal/models"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

var userList = []models.User{
	{
		ID:       "4f6f687e-522a-4203-810b-827bc6c09180",
		Name:     "Fatih Yerebakan",
		Email:    "fatihyerebakan@gmail.com",
		Password: []byte("Password12*"),
		Addresses: []models.Address{
			{
				ID:       "130beada-8339-4ee6-a754-725f43b8da98",
				Address:  "Levent",
				City:     "İstanbul",
				District: "Beşiktaş",
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
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:       "4f6f687e-522a-4203-810b-827bc6c09180",
		Name:     "Fatih Yerebakan",
		Email:    "fatihyerebakan@gmail.com",
		Password: []byte("Password12*"),
		Addresses: []models.Address{
			{
				ID:       "cc05be98-25af-4b22-b95c-bda2401bf6bc",
				Address:  "Kadıköy",
				City:     "İstanbul",
				District: "Kadıköy",
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
			{
				ID:       "d094bf9c-1f51-4936-afa9-ce5c5d807f09",
				Address:  "Üsküdar",
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
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

// MockUserRepository is a mock implementation of IUserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	// []models.Order => 1.Return model || error => 2.Return model
	return args.Get(0).([]models.User), nil
}

func (m *MockUserRepository) GetUserById(id string) (models.User, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return models.User{}, args.Error(1)
	}
	return args.Get(0).(models.User), nil
}

func (m *MockUserRepository) Insert(user models.User) (bool, error) {
	args := m.Called(user)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return true, nil
}

func (m *MockUserRepository) Update(user models.User) (bool, error) {
	args := m.Called(user)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return true, nil
}

func (m *MockUserRepository) Delete(id string) (bool, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return false, args.Error(1)
	}
	return true, nil
}

func TestUserService_GetAll_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetAll").Return(userList, nil)

	// Create an instance of OrderService with the mock repository
	userService := NewUserService(mockRepo)

	// Call the GetAll method
	users, err := userService.GetAll()

	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, userList, users)

	// Verify that the mock method was called
	mockRepo.AssertCalled(t, "GetAll")
}
