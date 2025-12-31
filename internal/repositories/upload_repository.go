package repositories

import (
	"github.com/JihadRinaldi/go-shop/internal/models"
	"gorm.io/gorm"
)

type UploadRepository struct {
	db *gorm.DB
}

func NewUploadRepository(db *gorm.DB) *UploadRepository {
	return &UploadRepository{db: db}
}

func (r *UploadRepository) CreateProductImage(image *models.ProductImage) error {
	return r.db.Create(image).Error
}

func (r *UploadRepository) GetProductImages(productID uint) ([]models.ProductImage, error) {
	var images []models.ProductImage
	if err := r.db.Where("product_id = ?", productID).Order("is_primary DESC, created_at ASC").Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (r *UploadRepository) DeleteProductImage(id uint) error {
	return r.db.Delete(&models.ProductImage{}, id).Error
}

func (r *UploadRepository) SetPrimaryImage(productID, imageID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Set all images for this product to non-primary
		if err := tx.Model(&models.ProductImage{}).Where("product_id = ?", productID).Update("is_primary", false).Error; err != nil {
			return err
		}

		// Set the specified image as primary
		if err := tx.Model(&models.ProductImage{}).Where("id = ? AND product_id = ?", imageID, productID).Update("is_primary", true).Error; err != nil {
			return err
		}

		return nil
	})
}
