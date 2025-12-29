package services

import (
	"errors"

	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/dto"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/JihadRinaldi/go-shop/internal/utils"
	"gorm.io/gorm"
)

type OrderService struct {
	db     *gorm.DB
	config *config.Config
}

func NewOrderService(db *gorm.DB, config *config.Config) *OrderService {
	return &OrderService{
		db:     db,
		config: config,
	}
}

func (s *OrderService) CreateOrder(userID uint) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
		if err := tx.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			return errors.New("cart not found")
		}

		if len(cart.CartItems) == 0 {
			return errors.New("cart is empty")
		}

		var totalAmount float64
		var orderItems []models.OrderItem

		for _, cartItem := range cart.CartItems {
			if cartItem.Product.Stock < cartItem.Quantity {
				return errors.New("insufficient stock for product: " + cartItem.Product.Name)
			}

			itemTotal := float64(cartItem.Quantity) * cartItem.Product.Price
			totalAmount += itemTotal

			orderItems = append(orderItems, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     cartItem.Product.Price,
			})

			cartItem.Product.Stock -= cartItem.Quantity
			if err := tx.Save(&cartItem.Product).Error; err != nil {
				return err
			}

			order := models.Order{
				UserID:      userID,
				Status:      models.OrderStatusPending,
				TotalAmount: totalAmount,
				OrderItems:  orderItems,
			}

			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
				return err
			}

			response := s.toOrderResponse(&order)
			orderResponse = &response
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orderResponse, nil
}

func (s *OrderService) GetOrder(userID uint, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order
	err := s.db.Preload("OrderItems").Where("id = ? AND user_id = ? AND deleted_at IS NULL", orderID, userID).First(&order).Error
	if err != nil {
		return nil, errors.New("order not found")
	}

	resp := s.toOrderResponse(&order)

	return &resp, nil
}

func (s *OrderService) GetOrders(userID uint, page, limit int) ([]dto.OrderResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	var orders []models.Order
	var total int64

	s.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total)

	if err := s.db.Preload("OrderItems.Product.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, nil, err
	}

	response := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		order := &orders[i]
		response[i] = s.toOrderResponse(order)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return response, meta, nil
}

func (s *OrderService) toOrderResponse(order *models.Order) dto.OrderResponse {
	var orderItems []dto.OrderItemResponse

	for _, item := range order.OrderItems {
		var images []dto.ProductImageResponse
		for _, img := range item.Product.Images {
			images = append(images, dto.ProductImageResponse{
				ID:        img.ID,
				URL:       img.URL,
				AltText:   img.AltText,
				CreatedAt: img.CreatedAt,
				IsPrimary: img.IsPrimary,
			})
		}
		orderItems = append(orderItems, dto.OrderItemResponse{
			ID: item.ID,
			Product: dto.ProductResponse{
				ID:          item.Product.ID,
				CategoryID:  item.Product.CategoryID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Stock:       item.Product.Stock,
				SKU:         item.Product.SKU,
				IsActive:    item.Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          item.Product.Category.ID,
					Name:        item.Product.Category.Name,
					Description: item.Product.Category.Description,
					IsActive:    item.Product.Category.IsActive,
					CreatedAt:   item.Product.Category.CreatedAt,
					UpdatedAt:   item.Product.Category.UpdatedAt,
				},
				Images: images,
			},
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		})
	}

	return dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		Status:      string(order.Status),
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItems,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}
