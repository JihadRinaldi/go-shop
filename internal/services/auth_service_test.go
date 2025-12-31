package services

import (
	"errors"
	"testing"
	"time"

	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/dto"
	"github.com/JihadRinaldi/go-shop/internal/mocks"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/JihadRinaldi/go-shop/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAuthService_Register(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepositoryInterface)
	mockCartRepo := new(mocks.MockCartRepositoryInterface)
	mockPublisher := new(mocks.MockPublisher)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:           "test-secret-key",
			ExpireIn:            15 * time.Minute,
			RefreshTokenExpires: 7 * 24 * time.Hour,
		},
	}

	service := &AuthService{
		config:         cfg,
		eventPublisher: mockPublisher,
		userRepo:       mockUserRepo,
		cartRepo:       mockCartRepo,
	}

	t.Run("success", func(t *testing.T) {
		req := &dto.RegisterRequest{
			Email:     "newuser@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
		}

		mockUserRepo.On("GetByEmail", req.Email).Return(nil, gorm.ErrRecordNotFound).Once()
		mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()
		mockCartRepo.On("Create", mock.AnythingOfType("*models.Cart")).Return(nil).Once()
		mockUserRepo.On("CreateRefreshToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil).Once()
		mockPublisher.On("Publish", "USER_LOGGED_IN", mock.AnythingOfType("*models.User"), mock.Anything).Return(nil).Once()

		result, err := service.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Email, result.User.Email)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		req := &dto.RegisterRequest{
			Email:     "existing@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "1234567890",
		}

		existingUser := &models.User{
			ID:    1,
			Email: req.Email,
		}

		mockUserRepo.On("GetByEmail", req.Email).Return(existingUser, nil).Once()

		result, err := service.Register(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "email already in use")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepositoryInterface)
	mockCartRepo := new(mocks.MockCartRepositoryInterface)
	mockPublisher := new(mocks.MockPublisher)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:           "test-secret-key",
			ExpireIn:            15 * time.Minute,
			RefreshTokenExpires: 7 * 24 * time.Hour,
		},
	}

	service := &AuthService{
		config:         cfg,
		eventPublisher: mockPublisher,
		userRepo:       mockUserRepo,
		cartRepo:       mockCartRepo,
	}

	t.Run("success", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := utils.HashPassword(password)

		req := &dto.LoginRequest{
			Email:    "user@example.com",
			Password: password,
		}

		user := &models.User{
			ID:       1,
			Email:    req.Email,
			Password: hashedPassword,
			IsActive: true,
			Role:     models.UserRoleCustomer,
		}

		mockUserRepo.On("GetByEmailAndActive", req.Email, true).Return(user, nil).Once()
		mockUserRepo.On("CreateRefreshToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil).Once()
		mockPublisher.On("Publish", "USER_LOGGED_IN", mock.AnythingOfType("*models.User"), mock.Anything).Return(nil).Once()

		result, err := service.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Email, result.User.Email)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		req := &dto.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockUserRepo.On("GetByEmailAndActive", req.Email, true).Return(nil, gorm.ErrRecordNotFound).Once()

		result, err := service.Login(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := utils.HashPassword(password)

		req := &dto.LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}

		user := &models.User{
			ID:       1,
			Email:    req.Email,
			Password: hashedPassword,
			IsActive: true,
		}

		mockUserRepo.On("GetByEmailAndActive", req.Email, true).Return(user, nil).Once()

		result, err := service.Login(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_Logout(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepositoryInterface)
	mockCartRepo := new(mocks.MockCartRepositoryInterface)
	mockPublisher := new(mocks.MockPublisher)

	cfg := &config.Config{}

	service := &AuthService{
		config:         cfg,
		eventPublisher: mockPublisher,
		userRepo:       mockUserRepo,
		cartRepo:       mockCartRepo,
	}

	t.Run("success", func(t *testing.T) {
		refreshToken := "valid-refresh-token"

		mockUserRepo.On("DeleteRefreshToken", refreshToken).Return(nil).Once()

		err := service.Logout(refreshToken)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("delete fails", func(t *testing.T) {
		refreshToken := "invalid-token"

		mockUserRepo.On("DeleteRefreshToken", refreshToken).Return(errors.New("delete failed")).Once()

		err := service.Logout(refreshToken)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}
