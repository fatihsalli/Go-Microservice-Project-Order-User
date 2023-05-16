package tests

import (
	"OrderUserProject/internal/models"
	"github.com/golang/mock/gomock"
)

// MockOrderRepository is a mock implementation of IOrderRepository
type MockOrderRepository struct {
	mockCtrl *gomock.Controller
}

// NewMockOrderRepository creates a new instance of MockOrderRepository
func NewMockOrderRepository(mockCtrl *gomock.Controller) *MockOrderRepository {
	return &MockOrderRepository{
		mockCtrl: mockCtrl,
	}
}

// Insert is a mock implementation of Insert method
func (m *MockOrderRepository) Insert(order models.Order) (bool, error) {
	// Mock the behavior you expect during testing
	// For example, return true to simulate a successful insert
	return true, nil
}
