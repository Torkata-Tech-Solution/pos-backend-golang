package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BusinessController struct {
	BusinessService service.BusinessService
}

func NewBusinessController(businessService service.BusinessService) *BusinessController {
	return &BusinessController{
		BusinessService: businessService,
	}
}

// @Tags         Business
// @Summary      Get all businesses
// @Description  Get all businesses with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of businesses"    default(10)
// @Param        search   query     string  false  "Search by name or address"
// @Router       /businesses [get]
// @Success      200  {object}  response.SuccessWithPaginate[model.Business]
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessController) GetBusinesses(c *fiber.Ctx) error {
	query := &validation.QueryUser{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	businesses, totalResults, err := b.BusinessService.GetBusinesses(c, query)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginate[model.Business]{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Businesses retrieved successfully",
		Results:      businesses,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Business
// @Summary      Get business by ID
// @Description  Get a specific business by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Business ID"
// @Router       /businesses/{id} [get]
// @Success      200  {object}  response.SuccessWithBusiness
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Business not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessController) GetBusinessByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid business ID format",
			Errors:  "Business ID must be a valid UUID",
		})
	}

	business, err := b.BusinessService.GetBusinessByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithBusiness{
		Code:     fiber.StatusOK,
		Status:   "OK",
		Message:  "Business retrieved successfully",
		Business: *business,
	})
}

// @Tags         Business
// @Summary      Create new business
// @Description  Create a new business
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        business  body      validation.CreateBusiness  true  "Business data"
// @Router       /businesses [post]
// @Success      201  {object}  response.SuccessWithBusiness
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Business domain already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessController) CreateBusiness(c *fiber.Ctx) error {
	req := new(validation.CreateBusiness)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	business, err := b.BusinessService.CreateBusiness(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithBusiness{
		Code:     fiber.StatusCreated,
		Status:   "Created",
		Message:  "Business created successfully",
		Business: *business,
	})
}

// @Tags         Business
// @Summary      Update business
// @Description  Update an existing business
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path      string                     true  "Business ID"
// @Param        business  body      validation.UpdateBusiness  true  "Business data to update"
// @Router       /businesses/{id} [put]
// @Success      200  {object}  response.SuccessWithBusiness
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Business not found"
// @Failure      409  {object}  response.ErrorDetails  "Business domain already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessController) UpdateBusiness(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid business ID format",
			Errors:  "Business ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateBusiness)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	business, err := b.BusinessService.UpdateBusiness(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithBusiness{
		Code:     fiber.StatusOK,
		Status:   "OK",
		Message:  "Business updated successfully",
		Business: *business,
	})
}

// @Tags         Business
// @Summary      Delete business
// @Description  Delete a business by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Business ID"
// @Router       /businesses/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Business not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessController) DeleteBusiness(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid business ID format",
			Errors:  "Business ID must be a valid UUID",
		})
	}

	err := b.BusinessService.DeleteBusiness(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Business deleted successfully",
	})
}
