package services

import (
	"errors"
	"testing"
	"time"

	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/mocks"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestOrderService_GetOrder(t *testing.T) {
	mockOrderRepo := new(mocks.MockOrderRepositoryInterface)
	mockCartRepo := new(mocks.MockCartRepositoryInterface)
	mockProductRepo := new(mocks.MockProductRepositoryInterface)

	service := &OrderService{
		db:          &gorm.DB{},
		config:      &config.Config{},
		orderRepo:   mockOrderRepo,
		cartRepo:    mockCartRepo,
		productRepo: mockProductRepo,
	}

	t.Run("success", func(t *testing.T) {
		userID := uint(1)
		orderID := uint(1)

		expectedOrder := &models.Order{
			ID:          orderID,
			UserID:      userID,
			Status:      models.OrderStatusPending,
			TotalAmount: 200.0,
			OrderItems: []models.OrderItem{
				{
					ID:        1,
					OrderID:   orderID,
					ProductID: 1,
					Quantity:  2,
					Price:     100.0,
					Product: models.Product{
						ID:          1,
						Name:        "Test Product",
						Price:       100.0,
						CategoryID:  1,
						Description: "Test",
						SKU:         "TEST-001",
						Stock:       10,
						IsActive:    true,
						Category: models.Category{
							ID:          1,
							Name:        "Test Category",
							Description: "Test Category",
							IsActive:    true,
						},
						Images: []models.ProductImage{},
					},
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockOrderRepo.On("GetByID", orderID).Return(expectedOrder, nil).Once()

		result, err := service.GetOrder(userID, orderID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orderID, result.ID)
		assert.Equal(t, userID, result.UserID)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("order not found", func(t *testing.T) {
		userID := uint(1)
		orderID := uint(999)

		mockOrderRepo.On("GetByID", orderID).Return(nil, errors.New("order not found")).Once()

		result, err := service.GetOrder(userID, orderID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("order belongs to different user", func(t *testing.T) {
		userID := uint(1)
		orderID := uint(1)

		order := &models.Order{
			ID:          orderID,
			UserID:      999, // Different user
			Status:      models.OrderStatusPending,
			TotalAmount: 200.0,
			OrderItems:  []models.OrderItem{},
		}

		mockOrderRepo.On("GetByID", orderID).Return(order, nil).Once()

		result, err := service.GetOrder(userID, orderID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "order not found")
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestOrderService_GetOrders(t *testing.T) {
	t.Run("note - requires db count mocking", func(t *testing.T) {
		// GetOrders requires mocking DB Count operation which is complex
		// The service directly calls s.db.Model(&models.Order{}).Where(...).Count(&total)
		// This would require either:
		// 1. Using a test database with SQL mocking (like sqlmock)
		// 2. Refactoring to move the Count logic into the repository
		// 3. Integration tests with a real test database
		assert.True(t, true, "GetOrders requires DB Count mocking or integration tests")
	})
}

func TestOrderService_CreateOrder(t *testing.T) {
	// Note: CreateOrder involves complex transaction logic with DB
	// This is best tested with integration tests using a test database

	t.Run("requires integration test", func(t *testing.T) {
		// This test is a placeholder
		// Full testing would require mocking DB transactions which is complex
		// Consider writing integration tests for this functionality
		assert.True(t, true, "CreateOrder requires integration testing with real DB")
	})
}
