package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SettingController struct {
	SettingService service.SettingService
}

func NewSettingController(settingService service.SettingService) *SettingController {
	return &SettingController{
		SettingService: settingService,
	}
}

// @Tags         Setting
// @Summary      Get all settings
// @Description  Get all settings with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page      query     int     false   "Page number"  default(1)
// @Param        limit     query     int     false   "Maximum number of settings"    default(10)
// @Param        search    query     string  false  "Search by key or description"
// @Param        outletId  query     string  false  "Filter by outlet ID"
// @Router       /settings [get]
// @Success      200  {object}  response.SuccessWithPaginatedSettings
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SettingController) GetSettings(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	outletID := c.Query("outletId", "")

	settings, totalResults, err := s.SettingService.GetSettings(c, query, outletID)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedSettings{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Settings retrieved successfully",
		Results:      settings,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Setting
// @Summary      Get setting by ID
// @Description  Get a specific setting by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Setting ID"
// @Router       /settings/{id} [get]
// @Success      200  {object}  response.SuccessWithSetting
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Setting not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SettingController) GetSettingByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid setting ID format",
			Errors:  "Setting ID must be a valid UUID",
		})
	}

	setting, err := s.SettingService.GetSettingByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithSetting{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Setting retrieved successfully",
		Setting: *setting,
	})
}

// @Tags         Setting
// @Summary      Get setting by key
// @Description  Get a specific setting by its key
// @Security     BearerAuth
// @Produce      json
// @Param        key      path      string  true  "Setting Key"
// @Param        outletId query     string  false "Outlet ID (optional)"
// @Router       /settings/key/{key} [get]
// @Success      200  {object}  response.SuccessWithSetting
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Setting not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SettingController) GetSettingByKey(c *fiber.Ctx) error {
	key := c.Params("key")
	outletID := c.Query("outletId", "")

	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Setting key is required",
			Errors:  "Setting key cannot be empty",
		})
	}

	setting, err := s.SettingService.GetSettingByKey(c, key, outletID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithSetting{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Setting retrieved successfully",
		Setting: *setting,
	})
}

// @Tags         Setting
// @Summary      Create new setting
// @Description  Create a new setting
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        setting  body      validation.CreateSetting  true  "Setting data"
// @Router       /settings [post]
// @Success      201  {object}  response.SuccessWithSetting
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Setting key already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SettingController) CreateSetting(c *fiber.Ctx) error {
	req := new(validation.CreateSetting)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	setting, err := s.SettingService.CreateSetting(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithSetting{
		Code:    fiber.StatusCreated,
		Status:  "Created",
		Message: "Setting created successfully",
		Setting: *setting,
	})
}

// @Tags         Setting
// @Summary      Update setting
// @Description  Update an existing setting
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Setting ID"
// @Param        setting  body      validation.UpdateSetting  true  "Setting data to update"
// @Router       /settings/{id} [put]
// @Success      200  {object}  response.SuccessWithSetting
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Setting not found"
// @Failure      409  {object}  response.ErrorDetails  "Setting key already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SettingController) UpdateSetting(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid setting ID format",
			Errors:  "Setting ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateSetting)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	setting, err := s.SettingService.UpdateSetting(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithSetting{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Setting updated successfully",
		Setting: *setting,
	})
}

// @Tags         Setting
// @Summary      Delete setting
// @Description  Delete a setting by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Setting ID"
// @Router       /settings/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Setting not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SettingController) DeleteSetting(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid setting ID format",
			Errors:  "Setting ID must be a valid UUID",
		})
	}

	err := s.SettingService.DeleteSetting(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Setting deleted successfully",
	})
}
