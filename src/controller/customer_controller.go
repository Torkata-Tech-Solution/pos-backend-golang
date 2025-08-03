package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CustomerController struct {
	CustomerService service.CustomerService
}

func NewCustomerController(customerService service.CustomerService) *CustomerController {
	return &CustomerController{
		CustomerService: customerService,
	}
}

// @Tags         Customer
// @Summary      Get all customers
// @Description  Get all customers with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page     query     int     false   "Page number"  default(1)
// @Param        limit    query     int     false   "Maximum number of customers"    default(10)
// @Param        search   query     string  false  "Search by name or email or phone"
// @Router       /customers [get]
// @Success      200  {object}  response.SuccessWithPaginatedCustomers
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CustomerController) GetCustomers(ctx *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", 10),
		Search: ctx.Query("search", ""),
	}

	customers, totalResults, err := c.CustomerService.GetCustomers(ctx, query)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedCustomers{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Customers retrieved successfully",
		Results:      customers,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Customer
// @Summary      Get customer by ID
// @Description  Get a specific customer by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Customer ID"
// @Router       /customers/{id} [get]
// @Success      200  {object}  response.SuccessWithCustomer
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Customer not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CustomerController) GetCustomerByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid customer ID format",
			Errors:  "Customer ID must be a valid UUID",
		})
	}

	customer, err := c.CustomerService.GetCustomerByID(ctx, id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithCustomer{
		Code:     fiber.StatusOK,
		Status:   "OK",
		Message:  "Customer retrieved successfully",
		Customer: *customer,
	})
}

// @Tags         Customer
// @Summary      Create new customer
// @Description  Create a new customer
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        customer  body      validation.CreateCustomer  true  "Customer data"
// @Router       /customers [post]
// @Success      201  {object}  response.SuccessWithCustomer
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Customer email already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CustomerController) CreateCustomer(ctx *fiber.Ctx) error {
	req := new(validation.CreateCustomer)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	customer, err := c.CustomerService.CreateCustomer(ctx, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.SuccessWithCustomer{
		Code:     fiber.StatusCreated,
		Status:   "Created",
		Message:  "Customer created successfully",
		Customer: *customer,
	})
}

// @Tags         Customer
// @Summary      Update customer
// @Description  Update an existing customer
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id        path      string                     true  "Customer ID"
// @Param        customer  body      validation.UpdateCustomer  true  "Customer data to update"
// @Router       /customers/{id} [put]
// @Success      200  {object}  response.SuccessWithCustomer
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Customer not found"
// @Failure      409  {object}  response.ErrorDetails  "Customer email already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CustomerController) UpdateCustomer(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid customer ID format",
			Errors:  "Customer ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateCustomer)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	customer, err := c.CustomerService.UpdateCustomer(ctx, id, req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithCustomer{
		Code:     fiber.StatusOK,
		Status:   "OK",
		Message:  "Customer updated successfully",
		Customer: *customer,
	})
}

// @Tags         Customer
// @Summary      Delete customer
// @Description  Delete a customer by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Customer ID"
// @Router       /customers/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Customer not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (c *CustomerController) DeleteCustomer(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid customer ID format",
			Errors:  "Customer ID must be a valid UUID",
		})
	}

	err := c.CustomerService.DeleteCustomer(ctx, id)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Customer deleted successfully",
	})
}
