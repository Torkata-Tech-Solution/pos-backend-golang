package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentMethodController struct {
	PaymentMethodService service.PaymentMethodService
}

func NewPaymentMethodController(paymentMethodService service.PaymentMethodService) *PaymentMethodController {
	return &PaymentMethodController{
		PaymentMethodService: paymentMethodService,
	}
}

// @Tags         Payment Method
// @Summary      Get all payment methods
// @Description  Get all payment methods with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of payment methods"    default(10)
// @Param        search   query     string  false  "Search by name or type"
// @Router       /payment-methods [get]
// @Success      200  {object}  response.SuccessWithPaginatedPaymentMethods
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PaymentMethodController) GetPaymentMethods(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	paymentMethods, totalResults, err := p.PaymentMethodService.GetPaymentMethods(c, query)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedPaymentMethods{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Payment methods retrieved successfully",
		Results:      paymentMethods,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Payment Method
// @Summary      Get payment method by ID
// @Description  Get a specific payment method by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Payment Method ID"
// @Router       /payment-methods/{id} [get]
// @Success      200  {object}  response.SuccessWithPaymentMethod
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Payment method not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PaymentMethodController) GetPaymentMethodByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid payment method ID format",
			Errors:  "Payment method ID must be a valid UUID",
		})
	}

	paymentMethod, err := p.PaymentMethodService.GetPaymentMethodByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaymentMethod{
		Code:          fiber.StatusOK,
		Status:        "OK",
		Message:       "Payment method retrieved successfully",
		PaymentMethod: *paymentMethod,
	})
}

// @Tags         Payment Method
// @Summary      Create new payment method
// @Description  Create a new payment method
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        paymentMethod  body      validation.CreatePaymentMethod  true  "Payment method data"
// @Router       /payment-methods [post]
// @Success      201  {object}  response.SuccessWithPaymentMethod
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Payment method name already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PaymentMethodController) CreatePaymentMethod(c *fiber.Ctx) error {
	req := new(validation.CreatePaymentMethod)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	paymentMethod, err := p.PaymentMethodService.CreatePaymentMethod(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithPaymentMethod{
		Code:          fiber.StatusCreated,
		Status:        "Created",
		Message:       "Payment method created successfully",
		PaymentMethod: *paymentMethod,
	})
}

// @Tags         Payment Method
// @Summary      Update payment method
// @Description  Update an existing payment method
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id             path      string                          true  "Payment Method ID"
// @Param        paymentMethod  body      validation.UpdatePaymentMethod  true  "Payment method data to update"
// @Router       /payment-methods/{id} [put]
// @Success      200  {object}  response.SuccessWithPaymentMethod
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Payment method not found"
// @Failure      409  {object}  response.ErrorDetails  "Payment method name already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PaymentMethodController) UpdatePaymentMethod(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid payment method ID format",
			Errors:  "Payment method ID must be a valid UUID",
		})
	}

	req := new(validation.UpdatePaymentMethod)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	paymentMethod, err := p.PaymentMethodService.UpdatePaymentMethod(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaymentMethod{
		Code:          fiber.StatusOK,
		Status:        "OK",
		Message:       "Payment method updated successfully",
		PaymentMethod: *paymentMethod,
	})
}

// @Tags         Payment Method
// @Summary      Delete payment method
// @Description  Delete a payment method by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Payment Method ID"
// @Router       /payment-methods/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Payment method not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PaymentMethodController) DeletePaymentMethod(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid payment method ID format",
			Errors:  "Payment method ID must be a valid UUID",
		})
	}

	err := p.PaymentMethodService.DeletePaymentMethod(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Payment method deleted successfully",
	})
}
