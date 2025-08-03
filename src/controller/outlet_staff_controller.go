package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OutletStaffController struct {
	OutletStaffService service.OutletStaffService
}

func NewOutletStaffController(outletStaffService service.OutletStaffService) *OutletStaffController {
	return &OutletStaffController{
		OutletStaffService: outletStaffService,
	}
}

// @Tags         Outlet Staff
// @Summary      Get all outlet staff
// @Description  Get all outlet staff with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page      query     int     false   "Page number"  default(1)
// @Param        limit     query     int     false   "Maximum number of staff"    default(10)
// @Param        search    query     string  false  "Search by name or email"
// @Param        outletId  query     string  false  "Filter by outlet ID"
// @Param        role      query     string  false  "Filter by role"
// @Router       /outlet-staff [get]
// @Success      200  {object}  response.SuccessWithPaginatedOutletStaff
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletStaffController) GetOutletStaff(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	outletID := c.Query("outletId", "")
	role := c.Query("role", "")

	staff, totalResults, err := o.OutletStaffService.GetOutletStaff(c, query, outletID, role)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedOutletStaff{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Outlet staff retrieved successfully",
		Results:      staff,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Outlet Staff
// @Summary      Get outlet staff by ID
// @Description  Get a specific outlet staff by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Outlet Staff ID"
// @Router       /outlet-staff/{id} [get]
// @Success      200  {object}  response.SuccessWithOutletStaff
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Outlet staff not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletStaffController) GetOutletStaffByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid outlet staff ID format",
			Errors:  "Outlet staff ID must be a valid UUID",
		})
	}

	staff, err := o.OutletStaffService.GetOutletStaffByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithOutletStaff{
		Code:        fiber.StatusOK,
		Status:      "OK",
		Message:     "Outlet staff retrieved successfully",
		OutletStaff: *staff,
	})
}

// @Tags         Outlet Staff
// @Summary      Get outlet staff by outlet ID
// @Description  Get all staff for a specific outlet
// @Security     BearerAuth
// @Produce      json
// @Param        outletId   path      string  true  "Outlet ID"
// @Router       /outlets/{outletId}/staff [get]
// @Success      200  {object}  response.SuccessWithPaginatedOutletStaff
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletStaffController) GetOutletStaffByOutletID(c *fiber.Ctx) error {
	outletID := c.Params("outletId")

	// Validate UUID format
	if _, err := uuid.Parse(outletID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid outlet ID format",
			Errors:  "Outlet ID must be a valid UUID",
		})
	}

	staff, err := o.OutletStaffService.GetOutletStaffByOutletID(c, outletID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedOutletStaff{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Outlet staff retrieved successfully",
		Results:      staff,
		Page:         1,
		Limit:        len(staff),
		TotalPages:   1,
		TotalResults: int64(len(staff)),
	})
}

// @Tags         Outlet Staff
// @Summary      Create new outlet staff
// @Description  Create a new outlet staff member
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        staff  body      validation.CreateOutletStaff  true  "Outlet staff data"
// @Router       /outlet-staff [post]
// @Success      201  {object}  response.SuccessWithOutletStaff
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "User already assigned to this outlet"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletStaffController) CreateOutletStaff(c *fiber.Ctx) error {
	req := new(validation.CreateOutletStaff)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	staff, err := o.OutletStaffService.CreateOutletStaff(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithOutletStaff{
		Code:        fiber.StatusCreated,
		Status:      "Created",
		Message:     "Outlet staff created successfully",
		OutletStaff: *staff,
	})
}

// @Tags         Outlet Staff
// @Summary      Update outlet staff
// @Description  Update an existing outlet staff member
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      string                         true  "Outlet Staff ID"
// @Param        staff  body      validation.UpdateOutletStaff  true  "Outlet staff data to update"
// @Router       /outlet-staff/{id} [put]
// @Success      200  {object}  response.SuccessWithOutletStaff
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Outlet staff not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletStaffController) UpdateOutletStaff(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid outlet staff ID format",
			Errors:  "Outlet staff ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateOutletStaff)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	staff, err := o.OutletStaffService.UpdateOutletStaff(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithOutletStaff{
		Code:        fiber.StatusOK,
		Status:      "OK",
		Message:     "Outlet staff updated successfully",
		OutletStaff: *staff,
	})
}

// @Tags         Outlet Staff
// @Summary      Delete outlet staff
// @Description  Delete an outlet staff member by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Outlet Staff ID"
// @Router       /outlet-staff/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Outlet staff not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (o *OutletStaffController) DeleteOutletStaff(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid outlet staff ID format",
			Errors:  "Outlet staff ID must be a valid UUID",
		})
	}

	err := o.OutletStaffService.DeleteOutletStaff(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Outlet staff deleted successfully",
	})
}
