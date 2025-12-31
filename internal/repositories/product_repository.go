package repositories

import (
	"github.com/JihadRinaldi/go-shop/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAll(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	query := r.db.Preload("Category").Preload("Images")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetByCategoryID(categoryID uint, limit, offset int) ([]models.Product, error) {
	var products []models.Product
	query := r.db.Preload("Category").Preload("Images").Where("category_id = ?", categoryID)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetBySKU(sku string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Category").Preload("Images").Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", quantity).Error
}
