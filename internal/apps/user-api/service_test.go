package user_api

import (
	"OrderUserProject/internal/models"
	"errors"
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
	// []models.User => 1.Return model || error => 2.Return model
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

	// Create an instance of UserService with the mock repository
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

func TestUserService_GetUserById_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockUserRepository)

	id := "4f6f687e-522a-4203-810b-827bc6c09180"

	// Define the expected result
	mockRepo.On("GetUserById", id).Return(userList[0], nil)

	// Create an instance of UserService with the mock repository
	userService := NewUserService(mockRepo)

	// Call the GetUserById method
	user, err := userService.GetUserById(id)

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, userList[0], user)

	// Verify that the mock method was called
	mockRepo.AssertCalled(t, "GetUserById", id)
}

func TestUserService_GetUserById_NotFoundFail(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockUserRepository)

	id := "4f6f687e-522a-4203-810b-827bc6c09185"
	expectedError := errors.New("not found error")

	// Define the expected result
	mockRepo.On("GetUserById", id).Return(models.User{}, expectedError)

	// Create an instance of UserService with the mock repository
	userService := NewUserService(mockRepo)

	// Call the GetUserById method
	user, err := userService.GetUserById(id)

	// Check error
	if !errors.Is(err, expectedError) {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}

	// Check nil user model
	if user.ID != "" {
		t.Error("Expected empty user model, but got a non-empty model!")
	}
}

func TestUserService_Insert_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockUserRepository)

	// Define the input and expected result
	user := models.User{
		ID:       "",
		Name:     "Fatih Şallı",
		Email:    "sallifatih@hotmail.com",
		Password: []byte("Password12*"),
		Addresses: []models.Address{
			{
				ID:       "2a98a38e-3e0c-4485-b124-7c216c91333a",
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
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	// We don't know exact user model because in service we have changed user model
	mockRepo.On("Insert", mock.AnythingOfType("models.User")).Return(true, nil)

	// Create an instance of UserService with the mock repository
	userService := NewUserService(mockRepo)

	// Call the Insert method
	result, err := userService.Insert(user)

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Addresses, result.Addresses)
	assert.Equal(t, user.Email, result.Email)

	// We don't know exact user model because in service we have changed user model
	mockRepo.AssertCalled(t, "Insert", mock.AnythingOfType("models.User"))
}

func TestUserService_Update_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockUserRepository)

	// Define the input and expected result
	user := userList[0]
	user.Name = "Emine Gamsız"

	// We don't know exact user model because in service we have changed user model
	mockRepo.On("Update", mock.AnythingOfType("models.User")).Return(true, nil)

	// Create an instance of UserService with the mock repository
	userService := NewUserService(mockRepo)

	// Call the Insert method
	result, err := userService.Update(user)

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, true, result)
	assert.Equal(t, "Emine Gamsız", user.Name)

	// We don't know exact user model because in service we have changed user model
	mockRepo.AssertCalled(t, "Update", mock.AnythingOfType("models.User"))
}

func TestUserService_Delete_Success(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockUserRepository)

	id := "4f6f687e-522a-4203-810b-827bc6c09180"

	// Define the expected result
	mockRepo.On("Delete", id).Return(true, nil)

	// Create an instance of UserService with the mock repository
	userService := NewUserService(mockRepo)

	// Call the Insert method
	result, err := userService.Delete(id)

	// Assert the result
	if err != nil {
		t.Error(err)
	}

	// Assert the result
	assert.Equal(t, true, result)

	// We don't know exact user model because in service we have changed user model
	mockRepo.AssertCalled(t, "Delete", id)
}

func TestUserService_Delete_NotFoundFail(t *testing.T) {
	// Create a mock instance
	mockRepo := new(MockUserRepository)

	expectedError := errors.New("not found error")
	id := "2b45ac31-6906-4e1e-82db-d9bcdbdb2143"

	// Define the expected result
	mockRepo.On("Delete", id).Return(false, expectedError)

	// Create an instance of UserService with the mock repository
	userService := NewUserService(mockRepo)

	// Call the Insert method
	result, err := userService.Delete(id)

	// Check error
	if !errors.Is(err, expectedError) {
		t.Errorf("Expected error: %v, but got: %v", expectedError, err)
	}

	// Assert the result
	assert.Equal(t, false, result)

	// We don't know exact user model because in service we have changed user model
	mockRepo.AssertCalled(t, "Delete", id)
}
