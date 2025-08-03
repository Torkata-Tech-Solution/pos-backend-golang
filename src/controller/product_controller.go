package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductController struct {
	ProductService         service.ProductService
	ProductCategoryService service.ProductCategoryService
}

func NewProductController(productService service.ProductService, productCategoryService service.ProductCategoryService) *ProductController {
	return &ProductController{
		ProductService:         productService,
		ProductCategoryService: productCategoryService,
	}
}

// Product Category Endpoints

// @Tags         Product Category
// @Summary      Get all product categories
// @Description  Get all product categories with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of categories"    default(10)
// @Param        search   query     string  false  "Search by name or description"
// @Router       /product-categories [get]
// @Success      200  {object}  response.SuccessWithPaginatedProductCategories
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) GetProductCategories(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	categories, totalResults, err := p.ProductCategoryService.GetProductCategories(c, query)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedProductCategories{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Product categories retrieved successfully",
		Results:      categories,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Product Category
// @Summary      Get product category by ID
// @Description  Get a specific product category by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Product Category ID"
// @Router       /product-categories/{id} [get]
// @Success      200  {object}  response.SuccessWithProductCategory
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Product category not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) GetProductCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid product category ID format",
			Errors:  "Product category ID must be a valid UUID",
		})
	}

	category, err := p.ProductCategoryService.GetProductCategoryByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithProductCategory{
		Code:            fiber.StatusOK,
		Status:          "OK",
		Message:         "Product category retrieved successfully",
		ProductCategory: *category,
	})
}

// @Tags         Product Category
// @Summary      Create new product category
// @Description  Create a new product category
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        category  body      validation.CreateProductCategory  true  "Product category data"
// @Router       /product-categories [post]
// @Success      201  {object}  response.SuccessWithProductCategory
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) CreateProductCategory(c *fiber.Ctx) error {
	req := new(validation.CreateProductCategory)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	category, err := p.ProductCategoryService.CreateProductCategory(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithProductCategory{
		Code:            fiber.StatusCreated,
		Status:          "Created",
		Message:         "Product category created successfully",
		ProductCategory: *category,
	})
}

// @Tags         Product Category
// @Summary      Update product category
// @Description  Update an existing product category
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path      string                             true  "Product Category ID"
// @Param        category  body      validation.UpdateProductCategory  true  "Product category data to update"
// @Router       /product-categories/{id} [put]
// @Success      200  {object}  response.SuccessWithProductCategory
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Product category not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) UpdateProductCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid product category ID format",
			Errors:  "Product category ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateProductCategory)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	category, err := p.ProductCategoryService.UpdateProductCategory(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithProductCategory{
		Code:            fiber.StatusOK,
		Status:          "OK",
		Message:         "Product category updated successfully",
		ProductCategory: *category,
	})
}

// @Tags         Product Category
// @Summary      Delete product category
// @Description  Delete a product category by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Product Category ID"
// @Router       /product-categories/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Product category not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) DeleteProductCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid product category ID format",
			Errors:  "Product category ID must be a valid UUID",
		})
	}

	err := p.ProductCategoryService.DeleteProductCategory(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product category deleted successfully",
	})
}

// Product Endpoints

// @Tags         Product
// @Summary      Get all products
// @Description  Get all products with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of products"    default(10)
// @Param        search   query     string  false  "Search by name or description"
// @Router       /products [get]
// @Success      200  {object}  response.SuccessWithPaginatedProducts
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) GetProducts(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	products, totalResults, err := p.ProductService.GetProducts(c, query)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedProducts{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Products retrieved successfully",
		Results:      products,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Product
// @Summary      Get product by ID
// @Description  Get a specific product by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Router       /products/{id} [get]
// @Success      200  {object}  response.SuccessWithProduct
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Product not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid product ID format",
			Errors:  "Product ID must be a valid UUID",
		})
	}

	product, err := p.ProductService.GetProductByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithProduct{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product retrieved successfully",
		Product: *product,
	})
}

// @Tags         Product
// @Summary      Create new product
// @Description  Create a new product
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        product  body      validation.CreateProduct  true  "Product data"
// @Router       /products [post]
// @Success      201  {object}  response.SuccessWithProduct
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) CreateProduct(c *fiber.Ctx) error {
	req := new(validation.CreateProduct)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	product, err := p.ProductService.CreateProduct(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithProduct{
		Code:    fiber.StatusCreated,
		Status:  "Created",
		Message: "Product created successfully",
		Product: *product,
	})
}

// @Tags         Product
// @Summary      Update product
// @Description  Update an existing product
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "Product ID"
// @Param        product  body      validation.UpdateProduct  true  "Product data to update"
// @Router       /products/{id} [put]
// @Success      200  {object}  response.SuccessWithProduct
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Product not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid product ID format",
			Errors:  "Product ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateProduct)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	product, err := p.ProductService.UpdateProduct(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithProduct{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product updated successfully",
		Product: *product,
	})
}

// @Tags         Product
// @Summary      Delete product
// @Description  Delete a product by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Router       /products/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Product not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *ProductController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid product ID format",
			Errors:  "Product ID must be a valid UUID",
		})
	}

	err := p.ProductService.DeleteProduct(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Product deleted successfully",
	})
}
