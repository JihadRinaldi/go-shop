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
	"gorm.io/gorm"
)

func TestProductService_GetProduct(t *testing.T) {
	mockProductRepo := new(mocks.MockProductRepositoryInterface)
	mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

	service := &ProductService{
		db:          &gorm.DB{},
		config:      &config.Config{},
		productRepo: mockProductRepo,
		uploadRepo:  mockUploadRepo,
	}

	t.Run("success", func(t *testing.T) {
		productID := uint(1)
		expectedProduct := &models.Product{
			ID:          productID,
			CategoryID:  1,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.0,
			Stock:       10,
			SKU:         "TEST-001",
			IsActive:    true,
			Category: models.Category{
				ID:          1,
				Name:        "Test Category",
				Description: "Test Category Description",
				IsActive:    true,
			},
			Images:    []models.ProductImage{},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockProductRepo.On("GetByID", productID).Return(expectedProduct, nil).Once()

		result, err := service.GetProduct(productID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedProduct.ID, result.ID)
		assert.Equal(t, expectedProduct.Name, result.Name)
		assert.Equal(t, expectedProduct.Price, result.Price)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("product not found", func(t *testing.T) {
		productID := uint(999)

		mockProductRepo.On("GetByID", productID).Return(nil, errors.New("product not found")).Once()

		result, err := service.GetProduct(productID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_CreateProduct(t *testing.T) {
	mockProductRepo := new(mocks.MockProductRepositoryInterface)
	mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

	service := &ProductService{
		db:          &gorm.DB{},
		config:      &config.Config{},
		productRepo: mockProductRepo,
		uploadRepo:  mockUploadRepo,
	}

	t.Run("success", func(t *testing.T) {
		req := &dto.CreateProductRequest{
			CategoryID:  1,
			Name:        "New Product",
			Description: "New Description",
			Price:       150.0,
			Stock:       20,
			SKU:         "NEW-001",
		}

		createdProduct := &models.Product{
			ID:          1,
			CategoryID:  req.CategoryID,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
			SKU:         req.SKU,
			IsActive:    true,
			Category: models.Category{
				ID:   1,
				Name: "Test Category",
			},
			Images: []models.ProductImage{},
		}

		mockProductRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil).Once()
		mockProductRepo.On("GetByID", mock.AnythingOfType("uint")).Return(createdProduct, nil).Once()

		result, err := service.CreateProduct(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.Price, result.Price)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("create fails", func(t *testing.T) {
		req := &dto.CreateProductRequest{
			CategoryID:  1,
			Name:        "New Product",
			Description: "New Description",
			Price:       150.0,
			Stock:       20,
			SKU:         "NEW-001",
		}

		mockProductRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(errors.New("create failed")).Once()

		result, err := service.CreateProduct(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockProductRepo := new(mocks.MockProductRepositoryInterface)
	mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

	service := &ProductService{
		db:          &gorm.DB{},
		config:      &config.Config{},
		productRepo: mockProductRepo,
		uploadRepo:  mockUploadRepo,
	}

	t.Run("success", func(t *testing.T) {
		productID := uint(1)
		isActive := true
		req := &dto.UpdateProductRequest{
			CategoryID:  1,
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       200.0,
			Stock:       30,
			IsActive:    &isActive,
		}

		existingProduct := &models.Product{
			ID:          productID,
			CategoryID:  1,
			Name:        "Old Product",
			Description: "Old Description",
			Price:       100.0,
			Stock:       10,
			SKU:         "TEST-001",
			IsActive:    true,
			Category: models.Category{
				ID:   1,
				Name: "Test Category",
			},
		}

		updatedProduct := &models.Product{
			ID:          productID,
			CategoryID:  req.CategoryID,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
			SKU:         "TEST-001",
			IsActive:    *req.IsActive,
			Category: models.Category{
				ID:   1,
				Name: "Test Category",
			},
		}

		mockProductRepo.On("GetByID", productID).Return(existingProduct, nil).Once()
		mockProductRepo.On("Update", mock.AnythingOfType("*models.Product")).Return(nil).Once()
		mockProductRepo.On("GetByID", productID).Return(updatedProduct, nil).Once()

		result, err := service.UpdateProduct(productID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Name, result.Name)
		assert.Equal(t, req.Price, result.Price)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("product not found", func(t *testing.T) {
		productID := uint(999)
		isActive := true
		req := &dto.UpdateProductRequest{
			CategoryID:  1,
			Name:        "Updated Product",
			Description: "Updated Description",
			Price:       200.0,
			Stock:       30,
			IsActive:    &isActive,
		}

		mockProductRepo.On("GetByID", productID).Return(nil, errors.New("product not found")).Once()

		result, err := service.UpdateProduct(productID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_DeleteProduct(t *testing.T) {
	mockProductRepo := new(mocks.MockProductRepositoryInterface)
	mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

	service := &ProductService{
		db:          &gorm.DB{},
		config:      &config.Config{},
		productRepo: mockProductRepo,
		uploadRepo:  mockUploadRepo,
	}

	t.Run("success", func(t *testing.T) {
		productID := uint(1)

		mockProductRepo.On("Delete", productID).Return(nil).Once()

		err := service.DeleteProduct(productID)

		assert.NoError(t, err)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("delete fails", func(t *testing.T) {
		productID := uint(1)

		mockProductRepo.On("Delete", productID).Return(errors.New("delete failed")).Once()

		err := service.DeleteProduct(productID)

		assert.Error(t, err)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestProductService_AddProductImage(t *testing.T) {
	mockProductRepo := new(mocks.MockProductRepositoryInterface)
	mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

	service := &ProductService{
		db:          &gorm.DB{},
		config:      &config.Config{},
		productRepo: mockProductRepo,
		uploadRepo:  mockUploadRepo,
	}

	t.Run("success - first image", func(t *testing.T) {
		productID := uint(1)
		url := "http://example.com/image.jpg"
		altText := "Test Image"

		mockUploadRepo.On("GetProductImages", productID).Return([]models.ProductImage{}, nil).Once()
		mockUploadRepo.On("CreateProductImage", mock.AnythingOfType("*models.ProductImage")).Return(nil).Once()

		err := service.AddProductImage(productID, url, altText)

		assert.NoError(t, err)
		mockUploadRepo.AssertExpectations(t)
	})

	t.Run("success - additional image", func(t *testing.T) {
		productID := uint(1)
		url := "http://example.com/image2.jpg"
		altText := "Test Image 2"

		existingImages := []models.ProductImage{
			{
				ID:        1,
				ProductID: productID,
				URL:       "http://example.com/image1.jpg",
				IsPrimary: true,
			},
		}

		mockUploadRepo.On("GetProductImages", productID).Return(existingImages, nil).Once()
		mockUploadRepo.On("CreateProductImage", mock.AnythingOfType("*models.ProductImage")).Return(nil).Once()

		err := service.AddProductImage(productID, url, altText)

		assert.NoError(t, err)
		mockUploadRepo.AssertExpectations(t)
	})

	t.Run("create fails", func(t *testing.T) {
		productID := uint(1)
		url := "http://example.com/image.jpg"
		altText := "Test Image"

		mockUploadRepo.On("GetProductImages", productID).Return([]models.ProductImage{}, nil).Once()
		mockUploadRepo.On("CreateProductImage", mock.AnythingOfType("*models.ProductImage")).Return(errors.New("create failed")).Once()

		err := service.AddProductImage(productID, url, altText)

		assert.Error(t, err)
		mockUploadRepo.AssertExpectations(t)
	})
}
