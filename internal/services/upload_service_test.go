package services

import (
	"errors"
	"mime/multipart"
	"testing"

	"github.com/JihadRinaldi/go-shop/internal/mocks"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadService_UploadProductImage(t *testing.T) {
	t.Run("success - first image", func(t *testing.T) {
		mockProvider := new(mocks.MockUploadProvider)
		mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

		service := &UploadService{
			provider:   mockProvider,
			uploadRepo: mockUploadRepo,
		}

		productID := uint(1)
		file := &multipart.FileHeader{
			Filename: "test-image.jpg",
			Size:     1024,
		}

		expectedURL := "http://example.com/products/1/test-image.jpg"

		mockProvider.On("UploadFile", file, mock.AnythingOfType("string")).Return(expectedURL, nil).Once()
		mockUploadRepo.On("GetProductImages", productID).Return([]models.ProductImage{}, nil).Once()
		mockUploadRepo.On("CreateProductImage", mock.AnythingOfType("*models.ProductImage")).Return(nil).Once()

		result, err := service.UploadProductImage(productID, file)

		assert.NoError(t, err)
		assert.Equal(t, expectedURL, result)
		mockProvider.AssertExpectations(t)
		mockUploadRepo.AssertExpectations(t)
	})

	t.Run("success - additional image", func(t *testing.T) {
		mockProvider := new(mocks.MockUploadProvider)
		mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

		service := &UploadService{
			provider:   mockProvider,
			uploadRepo: mockUploadRepo,
		}

		productID := uint(1)
		file := &multipart.FileHeader{
			Filename: "test-image2.png",
			Size:     2048,
		}

		expectedURL := "http://example.com/products/1/test-image2.png"
		existingImages := []models.ProductImage{
			{
				ID:        1,
				ProductID: productID,
				URL:       "http://example.com/products/1/existing.jpg",
				IsPrimary: true,
			},
		}

		mockProvider.On("UploadFile", file, mock.AnythingOfType("string")).Return(expectedURL, nil).Once()
		mockUploadRepo.On("GetProductImages", productID).Return(existingImages, nil).Once()
		mockUploadRepo.On("CreateProductImage", mock.AnythingOfType("*models.ProductImage")).Return(nil).Once()

		result, err := service.UploadProductImage(productID, file)

		assert.NoError(t, err)
		assert.Equal(t, expectedURL, result)
		mockProvider.AssertExpectations(t)
		mockUploadRepo.AssertExpectations(t)
	})

	t.Run("invalid file extension", func(t *testing.T) {
		mockProvider := new(mocks.MockUploadProvider)
		mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

		service := &UploadService{
			provider:   mockProvider,
			uploadRepo: mockUploadRepo,
		}

		productID := uint(1)
		file := &multipart.FileHeader{
			Filename: "test-document.pdf",
			Size:     1024,
		}

		result, err := service.UploadProductImage(productID, file)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "invalid image extension")
		mockProvider.AssertExpectations(t)
		mockUploadRepo.AssertExpectations(t)
	})

	t.Run("upload provider fails", func(t *testing.T) {
		mockProvider := new(mocks.MockUploadProvider)
		mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

		service := &UploadService{
			provider:   mockProvider,
			uploadRepo: mockUploadRepo,
		}

		productID := uint(1)
		file := &multipart.FileHeader{
			Filename: "test-image.jpg",
			Size:     1024,
		}

		mockProvider.On("UploadFile", file, mock.AnythingOfType("string")).Return("", errors.New("upload failed")).Once()

		result, err := service.UploadProductImage(productID, file)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "upload failed")
		mockProvider.AssertExpectations(t)
		mockUploadRepo.AssertExpectations(t)
	})

	t.Run("database save fails", func(t *testing.T) {
		mockProvider := new(mocks.MockUploadProvider)
		mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

		service := &UploadService{
			provider:   mockProvider,
			uploadRepo: mockUploadRepo,
		}

		productID := uint(1)
		file := &multipart.FileHeader{
			Filename: "test-image.jpg",
			Size:     1024,
		}

		expectedURL := "http://example.com/products/1/test-image.jpg"

		mockProvider.On("UploadFile", file, mock.AnythingOfType("string")).Return(expectedURL, nil).Once()
		mockUploadRepo.On("GetProductImages", productID).Return([]models.ProductImage{}, nil).Once()
		mockUploadRepo.On("CreateProductImage", mock.AnythingOfType("*models.ProductImage")).Return(errors.New("database error")).Once()

		result, err := service.UploadProductImage(productID, file)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockProvider.AssertExpectations(t)
		mockUploadRepo.AssertExpectations(t)
	})

	t.Run("valid file extensions", func(t *testing.T) {
		validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

		for _, ext := range validExtensions {
			mockProvider := new(mocks.MockUploadProvider)
			mockUploadRepo := new(mocks.MockUploadRepositoryInterface)

			service := &UploadService{
				provider:   mockProvider,
				uploadRepo: mockUploadRepo,
			}

			productID := uint(1)
			file := &multipart.FileHeader{
				Filename: "test-image" + ext,
				Size:     1024,
			}

			expectedURL := "http://example.com/products/1/test-image" + ext

			mockProvider.On("UploadFile", file, mock.AnythingOfType("string")).Return(expectedURL, nil).Once()
			mockUploadRepo.On("GetProductImages", productID).Return([]models.ProductImage{}, nil).Once()
			mockUploadRepo.On("CreateProductImage", mock.AnythingOfType("*models.ProductImage")).Return(nil).Once()

			result, err := service.UploadProductImage(productID, file)

			assert.NoError(t, err)
			assert.Equal(t, expectedURL, result)
			mockProvider.AssertExpectations(t)
			mockUploadRepo.AssertExpectations(t)
		}
	})
}
