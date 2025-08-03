package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TableController struct {
	TableService service.TableService
}

func NewTableController(tableService service.TableService) *TableController {
	return &TableController{
		TableService: tableService,
	}
}

// @Tags         Table
// @Summary      Get all tables
// @Description  Get all tables with pagination and search functionality
// @Security     BearerAuth
// @Produce      json
// @Param        page      query     int     false   "Page number"  default(1)
// @Param        limit     query     int     false   "Maximum number of tables"    default(10)
// @Param        search    query     string  false  "Search by table number or name"
// @Param        outletId  query     string  false  "Filter by outlet ID"
// @Router       /tables [get]
// @Success      200  {object}  response.SuccessWithPaginatedTables
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (t *TableController) GetTables(c *fiber.Ctx) error {
	query := &validation.QueryParams{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: c.Query("search", ""),
	}

	outletID := c.Query("outletId", "")

	tables, totalResults, err := t.TableService.GetTables(c, query, outletID)
	if err != nil {
		return err
	}

	totalPages := int64(math.Ceil(float64(totalResults) / float64(query.Limit)))

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedTables{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Tables retrieved successfully",
		Results:      tables,
		Page:         query.Page,
		Limit:        query.Limit,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	})
}

// @Tags         Table
// @Summary      Get table by ID
// @Description  Get a specific table by its ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Table ID"
// @Router       /tables/{id} [get]
// @Success      200  {object}  response.SuccessWithTable
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Table not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (t *TableController) GetTableByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid table ID format",
			Errors:  "Table ID must be a valid UUID",
		})
	}

	table, err := t.TableService.GetTableByID(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithTable{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Table retrieved successfully",
		Table:   *table,
	})
}

// @Tags         Table
// @Summary      Get tables by outlet ID
// @Description  Get all tables for a specific outlet
// @Security     BearerAuth
// @Produce      json
// @Param        outletId   path      string  true  "Outlet ID"
// @Router       /outlets/{outletId}/tables [get]
// @Success      200  {object}  response.SuccessWithPaginatedTables
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (t *TableController) GetTablesByOutletID(c *fiber.Ctx) error {
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

	tables, err := t.TableService.GetTablesByOutletID(c, outletID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithPaginatedTables{
		Code:         fiber.StatusOK,
		Status:       "OK",
		Message:      "Tables retrieved successfully",
		Results:      tables,
		Page:         1,
		Limit:        len(tables),
		TotalPages:   1,
		TotalResults: int64(len(tables)),
	})
}

// @Tags         Table
// @Summary      Create new table
// @Description  Create a new table
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        table  body      validation.CreateTable  true  "Table data"
// @Router       /tables [post]
// @Success      201  {object}  response.SuccessWithTable
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      409  {object}  response.ErrorDetails  "Table number already exists in outlet"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (t *TableController) CreateTable(c *fiber.Ctx) error {
	req := new(validation.CreateTable)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	table, err := t.TableService.CreateTable(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessWithTable{
		Code:    fiber.StatusCreated,
		Status:  "Created",
		Message: "Table created successfully",
		Table:   *table,
	})
}

// @Tags         Table
// @Summary      Update table
// @Description  Update an existing table
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id     path      string                  true  "Table ID"
// @Param        table  body      validation.UpdateTable  true  "Table data to update"
// @Router       /tables/{id} [put]
// @Success      200  {object}  response.SuccessWithTable
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Table not found"
// @Failure      409  {object}  response.ErrorDetails  "Table number already exists in outlet"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (t *TableController) UpdateTable(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid table ID format",
			Errors:  "Table ID must be a valid UUID",
		})
	}

	req := new(validation.UpdateTable)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid request body",
			Errors:  err.Error(),
		})
	}

	table, err := t.TableService.UpdateTable(c, id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithTable{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Table updated successfully",
		Table:   *table,
	})
}

// @Tags         Table
// @Summary      Delete table
// @Description  Delete a table by ID
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "Table ID"
// @Router       /tables/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorDetails  "Bad Request"
// @Failure      401  {object}  response.ErrorDetails  "Unauthorized"
// @Failure      404  {object}  response.ErrorDetails  "Table not found"
// @Failure      500  {object}  response.ErrorDetails  "Internal Server Error"
func (t *TableController) DeleteTable(c *fiber.Ctx) error {
	id := c.Params("id")

	// Validate UUID format
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorDetails{
			Code:    fiber.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid table ID format",
			Errors:  "Table ID must be a valid UUID",
		})
	}

	err := t.TableService.DeleteTable(c, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "OK",
		Message: "Table deleted successfully",
	})
}
