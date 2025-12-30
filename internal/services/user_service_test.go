package services

import (
	"errors"
	"testing"
	"time"

	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/dto"
	"github.com/JihadRinaldi/go-shop/internal/mocks"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetProfile(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepositoryInterface)
	cfg := &config.Config{}

	service := &UserService{
		config:   cfg,
		userRepo: mockUserRepo,
	}

	t.Run("success", func(t *testing.T) {
		userID := uint(1)
		expectedUser := &models.User{
			ID:        userID,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Role:      models.UserRoleCustomer,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserRepo.On("GetByID", userID).Return(expectedUser, nil).Once()

		result, err := service.GetProfile(userID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedUser.ID, result.ID)
		assert.Equal(t, expectedUser.Email, result.Email)
		assert.Equal(t, expectedUser.FirstName, result.FirstName)
		assert.Equal(t, expectedUser.LastName, result.LastName)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uint(999)
		mockUserRepo.On("GetByID", userID).Return(nil, errors.New("user not found")).Once()

		result, err := service.GetProfile(userID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepositoryInterface)
	cfg := &config.Config{}

	service := &UserService{
		config:   cfg,
		userRepo: mockUserRepo,
	}

	t.Run("success", func(t *testing.T) {
		userID := uint(1)
		updateReq := &dto.UpdateProfileRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Phone:     "0987654321",
		}

		existingUser := &models.User{
			ID:        userID,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Role:      models.UserRoleCustomer,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		updatedUser := &models.User{
			ID:        userID,
			Email:     "test@example.com",
			FirstName: updateReq.FirstName,
			LastName:  updateReq.LastName,
			Phone:     updateReq.Phone,
			Role:      models.UserRoleCustomer,
			IsActive:  true,
			CreatedAt: existingUser.CreatedAt,
			UpdatedAt: time.Now(),
		}

		mockUserRepo.On("GetByID", userID).Return(existingUser, nil).Once()
		mockUserRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockUserRepo.On("GetByID", userID).Return(updatedUser, nil).Once()

		result, err := service.UpdateProfile(userID, updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updateReq.FirstName, result.FirstName)
		assert.Equal(t, updateReq.LastName, result.LastName)
		assert.Equal(t, updateReq.Phone, result.Phone)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uint(999)
		updateReq := &dto.UpdateProfileRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Phone:     "0987654321",
		}

		mockUserRepo.On("GetByID", userID).Return(nil, errors.New("user not found")).Once()

		result, err := service.UpdateProfile(userID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("update fails", func(t *testing.T) {
		userID := uint(1)
		updateReq := &dto.UpdateProfileRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Phone:     "0987654321",
		}

		existingUser := &models.User{
			ID:        userID,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
			Role:      models.UserRoleCustomer,
			IsActive:  true,
		}

		mockUserRepo.On("GetByID", userID).Return(existingUser, nil).Once()
		mockUserRepo.On("Update", mock.AnythingOfType("*models.User")).Return(errors.New("update failed")).Once()

		result, err := service.UpdateProfile(userID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockUserRepo.AssertExpectations(t)
	})
}
