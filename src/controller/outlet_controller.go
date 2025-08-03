package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OutletController struct {
	OutletService service.OutletService
}

func NewOutletController(outletService service.OutletService) *OutletController {
	return &OutletController{
		OutletService: outletService,
	}
}

// @Tags         Outlet
// @Summary      Get all outlets
// @Description  Get all outlets with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of outlets"    default(10)
// @Param        search   query     string  false  "Search by name or address"
// @Router       /outlets [get]
// @Success      200  {object}  response.SuccessWithPaginatedOutlets
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletController) GetOutlets(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	outlets, totalResults, err := o.OutletService.GetOutlets(c, query)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedOutlets{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Outlets retrieved successfully",
		Results:      outlets,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Outlet
// @Summary      Get outlet by ID
// @Description  Get a specific outlet by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Outlet ID"
// @Router       /outlets/{id} [get]
// @Success      200  {object}  response.SuccessWithOutlet
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Outlet not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletController) GetOutletByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid outlet ID format",
			Errors:  "Outlet ID must be a valid UUID",
		})
	}

	outlet, err := o.OutletService.GetOutletByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithOutlet{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Outlet retrieved successfully",
		Outlet:  *outlet,
	})
}

// @Tags         Outlet
// @Summary      Get outlets by business ID
// @Description  Get all outlets for a specific business
// @Security     BearerAuth
// @Produce      json
// @Param        businessId   path      string  true  "Business ID"
// @Router       /businesses/{businessId}/outlets [get]
// @Success      200  {object}  response.SuccessWithPaginatedOutlets
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletController) GetOutletsByBusinessID(c *fiber.Ctx) error {
	businessID := c.Params("businessId")

	// Validate UUID format
	if _, err := uuid.Parse(businessID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid business ID format",
			Errors:  "Business ID must be a valid UUID",
		})
	}

	outlets, err := o.OutletService.GetOutletsByBusinessID(c, businessID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedOutlets{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Outlets retrieved successfully",
		Results:      outlets,
		Page:         1,
		Limit:        len(outlets),
		TotalPages:   1,
		TotalResults: int64(len(outlets)),
	})
}

// @Tags         Outlet
// @Summary      Create new outlet
// @Description  Create a new outlet
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        outlet  body      validation.CreateOutlet  true  "Outlet data"
// @Router       /outlets [post]
// @Success      201  {object}  response.SuccessWithOutlet
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Outlet already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletController) CreateOutlet(c *fiber.Ctx) error {
	req := new(validation.CreateOutlet)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	outlet, err := o.OutletService.CreateOutlet(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithOutlet{
		Code:    fiber.StatusCreated,
		Status:  "Created",
		Message: "Outlet created successfully",
		Outlet:  *outlet,
	})
}

// @Tags         Outlet
// @Summary      Update outlet
// @Description  Update an existing outlet
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id      path      string                   true  "Outlet ID"
// @Param        outlet  body      validation.UpdateOutlet  true  "Outlet data to update"
// @Router       /outlets/{id} [put]
// @Success      200  {object}  response.SuccessWithOutlet
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Outlet not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletController) UpdateOutlet(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid outlet ID format",
			Errors:  "Outlet ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateOutlet)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	outlet, err := o.OutletService.UpdateOutlet(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithOutlet{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Outlet updated successfully",
		Outlet:  *outlet,
	})
}

// @Tags         Outlet
// @Summary      Delete outlet
// @Description  Delete an outlet by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Outlet ID"
// @Router       /outlets/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Outlet not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletController) DeleteOutlet(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid outlet ID format",
			Errors:  "Outlet ID must be a valid UUID",
		})
	}

	err := o.OutletService.DeleteOutlet(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Outlet deleted successfully",
	})
}
