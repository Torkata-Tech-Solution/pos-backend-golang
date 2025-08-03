package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductCategoryService interface {
	GetProductCategories(c *fiber.Ctx, params *validation.QueryParams) ([]model.ProductCategory, int64, error)
	GetProductCategoryByID(c *fiber.Ctx, id string) (*model.ProductCategory, error)
	GetProductCategoriesByBusinessID(c *fiber.Ctx, businessID string) ([]model.ProductCategory, error)
	CreateProductCategory(c *fiber.Ctx, req *validation.CreateProductCategory) (*model.ProductCategory, error)
	UpdateProductCategory(c *fiber.Ctx, id string, req *validation.UpdateProductCategory) (*model.ProductCategory, error)
	DeleteProductCategory(c *fiber.Ctx, id string) error
}

type ProductService interface {
	GetProducts(c *fiber.Ctx, params *validation.QueryParams) ([]model.Product, int64, error)
	GetProductByID(c *fiber.Ctx, id string) (*model.Product, error)
	GetProductsByBusinessID(c *fiber.Ctx, businessID string) ([]model.Product, error)
	GetProductsByCategoryID(c *fiber.Ctx, categoryID string) ([]model.Product, error)
	CreateProduct(c *fiber.Ctx, req *validation.CreateProduct) (*model.Product, error)
	UpdateProduct(c *fiber.Ctx, id string, req *validation.UpdateProduct) (*model.Product, error)
	DeleteProduct(c *fiber.Ctx, id string) error
}

// Product Category Service Implementation
type productCategoryService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewProductCategoryService(db *gorm.DB, validate *validator.Validate) ProductCategoryService {
	return &productCategoryService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *productCategoryService) GetProductCategories(c *fiber.Ctx, params *validation.QueryParams) ([]model.ProductCategory, int64, error) {
	var categories []model.ProductCategory
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&categories).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search product categories: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&categories).Error; err != nil {
		s.Log.Errorf("Failed to get product categories: %+v", err)
		return nil, 0, err
	}

	return categories, totalResults, nil
}

func (s *productCategoryService) GetProductCategoryByID(c *fiber.Ctx, id string) (*model.ProductCategory, error) {
	category := new(model.ProductCategory)
	result := s.DB.WithContext(c.Context()).First(category, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product category not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get product category by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return category, nil
}

func (s *productCategoryService) GetProductCategoriesByBusinessID(c *fiber.Ctx, businessID string) ([]model.ProductCategory, error) {
	var categories []model.ProductCategory
	result := s.DB.WithContext(c.Context()).Where("business_id = ?", businessID).Find(&categories)
	if result.Error != nil {
		s.Log.Errorf("Failed to get product categories by business ID %s: %+v", businessID, result.Error)
		return nil, result.Error
	}
	return categories, nil
}

func (s *productCategoryService) CreateProductCategory(c *fiber.Ctx, req *validation.CreateProductCategory) (*model.ProductCategory, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	category := &model.ProductCategory{
		Name:        req.Name,
		Description: &req.Description,
		BusinessID:  uuid.MustParse(req.BusinessID),
	}

	result := s.DB.WithContext(c.Context()).Create(category)
	if result.Error != nil {
		s.Log.Errorf("Failed to create product category: %+v", result.Error)
	}

	return category, result.Error
}

func (s *productCategoryService) UpdateProductCategory(c *fiber.Ctx, id string, req *validation.UpdateProductCategory) (*model.ProductCategory, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Name == "" && req.Description == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	category := &model.ProductCategory{
		Name:        req.Name,
		Description: &req.Description,
	}

	result := s.DB.WithContext(c.Context()).Model(&model.ProductCategory{}).Where("id = ?", id).Updates(category)
	if result.Error != nil {
		s.Log.Errorf("Failed to update product category with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product category not found")
	}

	category, err := s.GetProductCategoryByID(c, id)
	if err != nil {
		return nil, err
	}

	return category, result.Error
}

func (s *productCategoryService) DeleteProductCategory(c *fiber.Ctx, id string) error {
	category := new(model.ProductCategory)

	result := s.DB.WithContext(c.Context()).Delete(category, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete product category with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Product category not found")
	}

	return nil
}

// Product Service Implementation
type productService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewProductService(db *gorm.DB, validate *validator.Validate) ProductService {
	return &productService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *productService) GetProducts(c *fiber.Ctx, params *validation.QueryParams) ([]model.Product, int64, error) {
	var products []model.Product
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("Category").Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&products).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search products: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&products).Error; err != nil {
		s.Log.Errorf("Failed to get products: %+v", err)
		return nil, 0, err
	}

	return products, totalResults, nil
}

func (s *productService) GetProductByID(c *fiber.Ctx, id string) (*model.Product, error) {
	product := new(model.Product)
	result := s.DB.WithContext(c.Context()).Preload("Category").First(product, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get product by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return product, nil
}

func (s *productService) GetProductsByBusinessID(c *fiber.Ctx, businessID string) ([]model.Product, error) {
	var products []model.Product
	result := s.DB.WithContext(c.Context()).Preload("Category").Where("business_id = ?", businessID).Find(&products)
	if result.Error != nil {
		s.Log.Errorf("Failed to get products by business ID %s: %+v", businessID, result.Error)
		return nil, result.Error
	}
	return products, nil
}

func (s *productService) GetProductsByCategoryID(c *fiber.Ctx, categoryID string) ([]model.Product, error) {
	var products []model.Product
	result := s.DB.WithContext(c.Context()).Preload("Category").Where("category_id = ?", categoryID).Find(&products)
	if result.Error != nil {
		s.Log.Errorf("Failed to get products by category ID %s: %+v", categoryID, result.Error)
		return nil, result.Error
	}
	return products, nil
}

func (s *productService) CreateProduct(c *fiber.Ctx, req *validation.CreateProduct) (*model.Product, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	product := &model.Product{
		Image:       &req.Image,
		Name:        req.Name,
		Description: &req.Description,
		Price:       req.Price,
		CategoryID:  uuid.MustParse(req.CategoryID),
		BusinessID:  uuid.MustParse(req.BusinessID),
	}

	result := s.DB.WithContext(c.Context()).Create(product)
	if result.Error != nil {
		s.Log.Errorf("Failed to create product: %+v", result.Error)
	}

	return product, result.Error
}

func (s *productService) UpdateProduct(c *fiber.Ctx, id string, req *validation.UpdateProduct) (*model.Product, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Name == "" && req.Description == "" && req.Price == 0 && req.CategoryID == "" && req.Image == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	updateFields := make(map[string]interface{})

	if req.Image != "" {
		updateFields["image"] = req.Image
	}
	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Description != "" {
		updateFields["description"] = req.Description
	}
	if req.Price > 0 {
		updateFields["price"] = req.Price
	}
	if req.CategoryID != "" {
		updateFields["category_id"] = uuid.MustParse(req.CategoryID)
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Product{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		s.Log.Errorf("Failed to update product with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	product, err := s.GetProductByID(c, id)
	if err != nil {
		return nil, err
	}

	return product, result.Error
}

func (s *productService) DeleteProduct(c *fiber.Ctx, id string) error {
	product := new(model.Product)

	result := s.DB.WithContext(c.Context()).Delete(product, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete product with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	return nil
}
