package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SaleController struct {
	SaleService service.SaleService
}

func NewSaleController(saleService service.SaleService) *SaleController {
	return &SaleController{
		SaleService: saleService,
	}
}

// @Tags         Sale
// @Summary      Get all sales
// @Description  Get all sales with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page       query     int     false   "Page number"  default(1)
// @Param        limit      query     int     false   "Maximum number of sales"    default(10)
// @Param        search     query     string  false  "Search by invoice number or customer name"
// @Param        outletId   query     string  false  "Filter by outlet ID"
// @Param        customerId query     string  false  "Filter by customer ID"
// @Param        status     query     string  false  "Filter by status (pending, completed, cancelled)"
// @Param        dateFrom   query     string  false  "Filter by date from (YYYY-MM-DD)"
// @Param        dateTo     query     string  false  "Filter by date to (YYYY-MM-DD)"
// @Router       /sales [get]
// @Success      200  {object}  response.SuccessWithPaginatedSales
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SaleController) GetSales(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	filters := &validation.SaleFilters{
		OutletID:      c.Query("outletId", ""),
		CustomerID:    c.Query("customerId", ""),
		PaymentMethod: c.Query("paymentMethod", ""),
		Status:        c.Query("status", ""),
		DateFrom:      c.Query("dateFrom", ""),
		DateTo:        c.Query("dateTo", ""),
	}

	sales, totalResults, err := s.SaleService.GetSales(c, query, filters)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedSales{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Sales retrieved successfully",
		Results:      sales,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Sale
// @Summary      Get sale by ID
// @Description  Get a specific sale by its ID with items details
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Sale ID"
// @Router       /sales/{id} [get]
// @Success      200  {object}  response.SuccessWithSale
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Sale not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SaleController) GetSaleByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid sale ID format",
			Errors:  "Sale ID must be a valid UUID",
		})
	}

	sale, err := s.SaleService.GetSaleByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithSale{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Sale retrieved successfully",
		Sale:    *sale,
	})
}

// @Tags         Sale
// @Summary      Get sale by invoice number
// @Description  Get a specific sale by its invoice number
// @Security     BearerAuth
// @Produce      json
// @Param        invoiceNumber  path      string  true  "Invoice Number"
// @Router       /sales/invoice/{invoiceNumber} [get]
// @Success      200  {object}  response.SuccessWithSale
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Sale not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SaleController) GetSaleByInvoiceNumber(c *fiber.Ctx) error {
	invoiceNumber := c.Params("invoiceNumber")

	if invoiceNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invoice number is required",
			Errors:  "Invoice number cannot be empty",
		})
	}

	sale, err := s.SaleService.GetSaleByInvoiceNumber(c, invoiceNumber)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithSale{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Sale retrieved successfully",
		Sale:    *sale,
	})
}

// @Tags         Sale
// @Summary      Create new sale
// @Description  Create a new sale transaction
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        sale  body      validation.CreateSale  true  "Sale data"
// @Router       /sales [post]
// @Success      201  {object}  response.SuccessWithSale
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Invoice number already exists"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SaleController) CreateSale(c *fiber.Ctx) error {
	req := new(validation.CreateSale)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	sale, err := s.SaleService.CreateSale(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithSale{
		Code:    fiber.StatusCreated,
		Status:  "Created",
		Message: "Sale created successfully",
		Sale:    *sale,
	})
}

// @Tags         Sale
// @Summary      Update sale
// @Description  Update an existing sale
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "Sale ID"
// @Param        sale  body      validation.UpdateSale  true  "Sale data to update"
// @Router       /sales/{id} [put]
// @Success      200  {object}  response.SuccessWithSale
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Sale not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SaleController) UpdateSale(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid sale ID format",
			Errors:  "Sale ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateSale)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	sale, err := s.SaleService.UpdateSale(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithSale{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Sale updated successfully",
		Sale:    *sale,
	})
}

// @Tags         Sale
// @Summary      Update sale status
// @Description  Update the status of a sale (pending, completed, cancelled)
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      string                      true  "Sale ID"
// @Param        status body      validation.UpdateSaleStatus  true  "Sale status data"
// @Router       /sales/{id}/status [patch]
// @Success      200  {object}  response.SuccessWithSale
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Sale not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SaleController) UpdateSaleStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid sale ID format",
			Errors:  "Sale ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateSaleStatus)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	sale, err := s.SaleService.UpdateSaleStatus(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithSale{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Sale status updated successfully",
		Sale:    *sale,
	})
}

// @Tags         Sale
// @Summary      Delete sale
// @Description  Delete a sale by ID (soft delete)
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Sale ID"
// @Router       /sales/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Sale not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SaleController) DeleteSale(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid sale ID format",
			Errors:  "Sale ID must be a valid UUID",
		})
	}

	err := s.SaleService.DeleteSale(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Sale deleted successfully",
	})
}

// @Tags         Sale
// @Summary      Get sales report
// @Description  Get sales report with various aggregations
// @Security     BearerAuth
// @Produce      json
// @Param        outletId  query     string  false  "Filter by outlet ID"
// @Param        dateFrom  query     string  false  "Filter by date from (YYYY-MM-DD)"
// @Param        dateTo    query     string  false  "Filter by date to (YYYY-MM-DD)"
// @Param        groupBy   query     string  false  "Group by (day, week, month, year)"
// @Router       /sales/report [get]
// @Success      200  {object}  response.SuccessWithSalesReport
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (s *SaleController) GetSalesReport(c *fiber.Ctx) error {
	filters := &validation.SalesReportFilters{
		OutletID: c.Query("outletId", ""),
		DateFrom: c.Query("dateFrom", ""),
		DateTo:   c.Query("dateTo", ""),
		GroupBy:  c.Query("groupBy", "day"),
	}

	report, err := s.SaleService.GetSalesReport(c, filters)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithSalesReport{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Sales report retrieved successfully",
		Report:  report,
	})
}
