package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/dto"
	"github.com/JihadRinaldi/go-shop/internal/events"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/JihadRinaldi/go-shop/internal/repositories"
	"github.com/JihadRinaldi/go-shop/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db             *gorm.DB
	config         *config.Config
	eventPublisher events.Publisher
	userRepo       repositories.UserRepositoryInterface
	cartRepo       repositories.CartRepositoryInterface
}

func NewAuthService(db *gorm.DB, config *config.Config, eventPublisher events.Publisher) *AuthService {
	return &AuthService{
		db:             db,
		config:         config,
		eventPublisher: eventPublisher,
		userRepo:       repositories.NewUserRepository(db),
		cartRepo:       repositories.NewCartRepository(db),
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	cart := models.Cart{UserID: user.ID}
	if err := s.cartRepo.Create(&cart); err != nil {
		fmt.Println("Unable to create cart")
	}

	return s.generateAuthResponse(&user)

}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmailAndActive(req.Email, true)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.SecretKey)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// var refreshToken models.RefreshToken
	refreshToken, err := s.userRepo.GetValidRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := s.userRepo.DeleteRefreshTokenByID(refreshToken.ID); err != nil {
		log.Println(err)
		_ = err
	}

	return s.generateAuthResponse(user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.DeleteRefreshToken(refreshToken)
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateToken(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiredAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}

	if err := s.userRepo.CreateRefreshToken(&refreshTokenModel); err != nil {
		return nil, err
	}

	err = s.eventPublisher.Publish("USER_LOGGED_IN", user, nil)
	if err != nil {
		fmt.Println("Failed to publish USER_LOGGED_IN event:", err)
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
