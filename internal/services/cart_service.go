package services

import (
	"errors"

	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/dto"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/JihadRinaldi/go-shop/internal/repositories"
	"gorm.io/gorm"
)

type CartService struct {
	db          *gorm.DB
	config      *config.Config
	cartRepo    repositories.CartRepositoryInterface
	productRepo repositories.ProductRepositoryInterface
}

func NewCartService(db *gorm.DB, config *config.Config) *CartService {
	return &CartService{
		db:          db,
		config:      config,
		cartRepo:    repositories.NewCartRepository(db),
		productRepo: repositories.NewProductRepository(db),
	}
}

func (s *CartService) GetCart(userID uint) (*dto.CartResponse, error) {
	var cart models.Cart
	err := s.db.Preload("CartItems.Product.Category").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return s.toCartResponse(&cart), nil
}

func (s *CartService) AddToCart(userID uint, req dto.AddToCartRequest) (*dto.CartResponse, error) {
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient product stock")
	}

	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		cart = &models.Cart{UserID: userID}
		if err := s.cartRepo.Create(cart); err != nil {
			return nil, err
		}
	}

	// Check if item already exists in cart
	var cartItem models.CartItem
	if err := s.db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID).First(&cartItem).Error; err != nil {
		// Create new cart item
		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		s.db.Create(&cartItem)
	} else {
		// Update existing cart item
		cartItem.Quantity += req.Quantity
		if cartItem.Quantity > product.Stock {
			return nil, errors.New("insufficient stock")
		}
		s.db.Save(&cartItem)
	}

	return s.GetCart(userID)
}

func (s *CartService) UpdateCartItem(userID uint, itemID uint, req dto.UpdateCartItemRequest) (*dto.CartResponse, error) {
	var cartItem models.CartItem
	err := s.db.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ? AND cart_items.id = ?", userID, itemID).
		First(&cartItem).Error
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	product, err := s.productRepo.GetByID(cartItem.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if req.Quantity > product.Stock {
		return nil, errors.New("insufficient product stock")
	}

	cartItem.Quantity = req.Quantity
	s.db.Save(&cartItem)

	return s.GetCart(userID)
}

func (s *CartService) RemoveCartItem(userID uint, itemID uint) error {
	return s.db.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("carts.user_id = ? ANd cart_items.id = ?", userID, itemID).
		Delete(&models.CartItem{}).Error
}

func (s *CartService) toCartResponse(cart *models.Cart) *dto.CartResponse {
	cartItems := make([]dto.CartItemResponse, len(cart.CartItems))
	var total float64

	for i, item := range cart.CartItems {
		subTotal := float64(cart.CartItems[i].Quantity) * cart.CartItems[i].Product.Price
		total += subTotal

		cartItems[i] = dto.CartItemResponse{
			ID: item.ID,
			Product: dto.ProductResponse{
				ID:          item.Product.ID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Category: dto.CategoryResponse{
					ID:   item.Product.Category.ID,
					Name: item.Product.Category.Name,
				},
			},
			Quantity: item.Quantity,
			Subtotal: subTotal,
		}
	}

	return &dto.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CartItems: cartItems,
		Total:     total,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}
