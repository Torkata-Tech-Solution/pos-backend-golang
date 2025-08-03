package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BusinessUserController struct {
	BusinessUserService service.BusinessUserService
}

func NewBusinessUserController(businessUserService service.BusinessUserService) *BusinessUserController {
	return &BusinessUserController{
		BusinessUserService: businessUserService,
	}
}

// @Tags         Business User
// @Summary      Get all business users
// @Description  Get all business users with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page       query     int     false   "Page number"  default(1)
// @Param        limit      query     int     false   "Maximum number of business users"    default(10)
// @Param        search     query     string  false  "Search by name or email"
// @Param        businessId query     string  false  "Filter by business ID"
// @Param        role       query     string  false  "Filter by role"
// @Router       /business-users [get]
// @Success      200  {object}  response.SuccessWithPaginatedBusinessUsers
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessUserController) GetBusinessUsers(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	businessUsers, totalResults, err := b.BusinessUserService.GetBusinessUsers(c, query)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedBusinessUsers{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Business users retrieved successfully",
		Results:      businessUsers,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Business User
// @Summary      Get business user by ID
// @Description  Get a specific business user by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Business User ID"
// @Router       /business-users/{id} [get]
// @Success      200  {object}  response.SuccessWithBusinessUser
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Business user not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessUserController) GetBusinessUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid business user ID format",
			Errors:  "Business user ID must be a valid UUID",
		})
	}

	businessUser, err := b.BusinessUserService.GetBusinessUserByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithBusinessUser{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Business user retrieved successfully",
		BusinessUser: *businessUser,
	})
}

// @Tags         Business User
// @Summary      Get business users by business ID
// @Description  Get all users for a specific business
// @Security     BearerAuth
// @Produce      json
// @Param        businessId   path      string  true  "Business ID"
// @Router       /businesses/{businessId}/users [get]
// @Success      200  {object}  response.SuccessWithPaginatedBusinessUsers
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessUserController) GetBusinessUsersByBusinessID(c *fiber.Ctx) error {
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

	businessUsers, err := b.BusinessUserService.GetBusinessUsersByBusinessID(c, businessID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedBusinessUsers{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Business users retrieved successfully",
		Results:      businessUsers,
		Page:         1,
		Limit:        len(businessUsers),
		TotalPages:   1,
		TotalResults: int64(len(businessUsers)),
	})
}

// @Tags         Business User
// @Summary      Create new business user
// @Description  Create a new business user (assign user to business)
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        businessUser  body      validation.CreateBusinessUser  true  "Business user data"
// @Router       /business-users [post]
// @Success      201  {object}  response.SuccessWithBusinessUser
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "User already assigned to this business"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessUserController) CreateBusinessUser(c *fiber.Ctx) error {
	req := new(validation.CreateBusinessUser)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	businessUser, err := b.BusinessUserService.CreateBusinessUser(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithBusinessUser{
		Code:         fiber.StatusCreated,
		Status:       "Created",
		Message:      "Business user created successfully",
		BusinessUser: *businessUser,
	})
}

// @Tags         Business User
// @Summary      Update business user
// @Description  Update an existing business user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id            path      string                         true  "Business User ID"
// @Param        businessUser  body      validation.UpdateBusinessUser  true  "Business user data to update"
// @Router       /business-users/{id} [put]
// @Success      200  {object}  response.SuccessWithBusinessUser
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Business user not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessUserController) UpdateBusinessUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid business user ID format",
			Errors:  "Business user ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateBusinessUser)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	businessUser, err := b.BusinessUserService.UpdateBusinessUser(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithBusinessUser{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Business user updated successfully",
		BusinessUser: *businessUser,
	})
}

// @Tags         Business User
// @Summary      Delete business user
// @Description  Delete a business user by ID (remove user from business)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Business User ID"
// @Router       /business-users/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Business user not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (b *BusinessUserController) DeleteBusinessUser(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid business user ID format",
			Errors:  "Business user ID must be a valid UUID",
		})
	}

	err := b.BusinessUserService.DeleteBusinessUser(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Business user deleted successfully",
	})
}
