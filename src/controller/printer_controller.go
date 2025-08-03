package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PrinterController struct {
	PrinterService service.PrinterService
}

func NewPrinterController(printerService service.PrinterService) *PrinterController {
	return &PrinterController{
		PrinterService: printerService,
	}
}

// @Tags         Printer
// @Summary      Get all printers
// @Description  Get all printers with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page      query     int     false   "Page number"  default(1)
// @Param        limit     query     int     false   "Maximum number of printers"    default(10)
// @Param        search    query     string  false  "Search by name or type"
// @Param        outletId  query     string  false  "Filter by outlet ID"
// @Router       /printers [get]
// @Success      200  {object}  response.SuccessWithPaginatedPrinters
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PrinterController) GetPrinters(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	outletID := c.Query("outletId", "")

	printers, totalResults, err := p.PrinterService.GetPrinters(c, query, outletID)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedPrinters{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Printers retrieved successfully",
		Results:      printers,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Printer
// @Summary      Get printer by ID
// @Description  Get a specific printer by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Printer ID"
// @Router       /printers/{id} [get]
// @Success      200  {object}  response.SuccessWithPrinter
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Printer not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PrinterController) GetPrinterByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid printer ID format",
			Errors:  "Printer ID must be a valid UUID",
		})
	}

	printer, err := p.PrinterService.GetPrinterByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPrinter{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Printer retrieved successfully",
		Printer: *printer,
	})
}

// @Tags         Printer
// @Summary      Get printers by outlet ID
// @Description  Get all printers for a specific outlet
// @Security     BearerAuth
// @Produce      json
// @Param        outletId   path      string  true  "Outlet ID"
// @Router       /outlets/{outletId}/printers [get]
// @Success      200  {object}  response.SuccessWithPaginatedPrinters
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PrinterController) GetPrintersByOutletID(c *fiber.Ctx) error {
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

	printers, err := p.PrinterService.GetPrintersByOutletID(c, outletID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedPrinters{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Printers retrieved successfully",
		Results:      printers,
		Page:         1,
		Limit:        len(printers),
		TotalPages:   1,
		TotalResults: int64(len(printers)),
	})
}

// @Tags         Printer
// @Summary      Create new printer
// @Description  Create a new printer
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        printer  body      validation.CreatePrinter  true  "Printer data"
// @Router       /printers [post]
// @Success      201  {object}  response.SuccessWithPrinter
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Printer name already exists in outlet"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PrinterController) CreatePrinter(c *fiber.Ctx) error {
	req := new(validation.CreatePrinter)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	printer, err := p.PrinterService.CreatePrinter(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithPrinter{
		Code:    fiber.StatusCreated,
		Status:  "Created",
		Message: "Printer created successfully",
		Printer: *printer,
	})
}

// @Tags         Printer
// @Summary      Update printer
// @Description  Update an existing printer
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Printer ID"
// @Param        printer  body      validation.UpdatePrinter  true  "Printer data to update"
// @Router       /printers/{id} [put]
// @Success      200  {object}  response.SuccessWithPrinter
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Printer not found"
// @Failure      409  {object}  response.ErrorDetails  "Printer name already exists in outlet"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PrinterController) UpdatePrinter(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid printer ID format",
			Errors:  "Printer ID must be a valid UUID",
		})
	}

	req := new(validation.UpdatePrinter)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	printer, err := p.PrinterService.UpdatePrinter(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPrinter{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Printer updated successfully",
		Printer: *printer,
	})
}

// @Tags         Printer
// @Summary      Delete printer
// @Description  Delete a printer by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Printer ID"
// @Router       /printers/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Printer not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (p *PrinterController) DeletePrinter(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid printer ID format",
			Errors:  "Printer ID must be a valid UUID",
		})
	}

	err := p.PrinterService.DeletePrinter(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Printer deleted successfully",
	})
}
