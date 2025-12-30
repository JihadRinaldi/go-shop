package services

import (
	"github.com/JihadRinaldi/go-shop/internal/config"
	"github.com/JihadRinaldi/go-shop/internal/dto"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/JihadRinaldi/go-shop/internal/repositories"
	"github.com/JihadRinaldi/go-shop/internal/utils"
	"gorm.io/gorm"
)

type ProductService struct {
	db          *gorm.DB
	config      *config.Config
	productRepo repositories.ProductRepositoryInterface
	uploadRepo  repositories.UploadRepositoryInterface
}

func NewProductService(db *gorm.DB, config *config.Config) *ProductService {
	return &ProductService{
		db:          db,
		config:      config,
		productRepo: repositories.NewProductRepository(db),
		uploadRepo:  repositories.NewUploadRepository(db),
	}
}

func (s *ProductService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *ProductService) GetCategories() ([]*dto.CategoryResponse, error) {
	var categories []models.Category
	if err := s.db.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, err
	}

	var response = make([]*dto.CategoryResponse, len(categories))
	for i := range categories {
		response[i] = &dto.CategoryResponse{
			ID:          categories[i].ID,
			Name:        categories[i].Name,
			Description: categories[i].Description,
			IsActive:    categories[i].IsActive,
			CreatedAt:   categories[i].CreatedAt,
			UpdatedAt:   categories[i].UpdatedAt,
		}
	}

	return response, nil
}

func (s *ProductService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	var category models.Category
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Description = req.Description
	category.IsActive = *req.IsActive
	if err := s.db.Save(&category).Error; err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}

func (s *ProductService) DeleteCategory(id uint) error {
	return s.db.Delete(&models.Category{}, id).Error
}

func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		SKU:         req.SKU,
	}

	if err := s.productRepo.Create(&product); err != nil {
		return nil, err
	}

	return s.GetProduct(product.ID)
}

func (s *ProductService) GetProducts(page, limit int) ([]dto.ProductResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	s.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total)

	products, err := s.productRepo.GetAll(limit, offset)
	if err != nil {
		return nil, nil, err
	}

	response := make([]dto.ProductResponse, len(products))
	for i := range products {
		response[i] = s.convertToProductResponse(&products[i])
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

func (s *ProductService) GetProduct(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := s.convertToProductResponse(product)
	return &response, nil
}

func (s *ProductService) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock
	product.IsActive = *req.IsActive

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return s.GetProduct(product.ID)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.productRepo.Delete(id)
}

func (s *ProductService) AddProductImage(productID uint, url, altText string) error {
	images, _ := s.uploadRepo.GetProductImages(productID)
	isPrimary := len(images) == 0

	image := models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   altText,
		IsPrimary: isPrimary,
	}

	return s.uploadRepo.CreateProductImage(&image)
}

func (s *ProductService) convertToProductResponse(product *models.Product) dto.ProductResponse {
	images := make([]dto.ProductImageResponse, len(product.Images))
	for i := range product.Images {
		images[i] = dto.ProductImageResponse{
			ID:        product.Images[i].ID,
			URL:       product.Images[i].URL,
			AltText:   product.Images[i].AltText,
			IsPrimary: product.Images[i].IsPrimary,
			CreatedAt: product.Images[i].CreatedAt,
		}
	}

	return dto.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		SKU:         product.SKU,
		IsActive:    product.IsActive,
		Category: dto.CategoryResponse{
			ID:          product.Category.ID,
			Name:        product.Category.Name,
			Description: product.Category.Description,
			IsActive:    product.Category.IsActive,
			CreatedAt:   product.Category.CreatedAt,
			UpdatedAt:   product.Category.UpdatedAt,
		},
		Images:    images,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}
}
