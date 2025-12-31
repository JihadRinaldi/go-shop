package repositories

import "github.com/JihadRinaldi/go-shop/internal/models"

type UserRepositoryInterface interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmailAndActive(email string, isActive bool) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error

	CreateRefreshToken(token *models.RefreshToken) error
	GetValidRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokenByID(id uint) error
}

type CartRepositoryInterface interface {
	GetByUserID(userID uint) (*models.Cart, error)
	Create(cart *models.Cart) error
	Update(cart *models.Cart) error
	Delete(id uint) error
}

type ProductRepositoryInterface interface {
	GetByID(id uint) (*models.Product, error)
	GetAll(limit, offset int) ([]models.Product, error)
	GetByCategoryID(categoryID uint, limit, offset int) ([]models.Product, error)
	GetBySKU(sku string) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
}

type OrderRepositoryInterface interface {
	GetByID(id uint) (*models.Order, error)
	GetByUserID(userID uint, limit, offset int) ([]models.Order, error)
	GetAll(limit, offset int) ([]models.Order, error)
	Create(order *models.Order) error
	Update(order *models.Order) error
	UpdateStatus(id uint, status models.OrderStatus) error
	Delete(id uint) error
}

type UploadRepositoryInterface interface {
	CreateProductImage(image *models.ProductImage) error
	GetProductImages(productID uint) ([]models.ProductImage, error)
	DeleteProductImage(id uint) error
	SetPrimaryImage(productID, imageID uint) error
}
