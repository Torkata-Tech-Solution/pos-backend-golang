package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TableService interface {
	GetTables(c *fiber.Ctx, params *validation.QueryParams, outletID string) ([]model.Table, int64, error)
	GetTableByID(c *fiber.Ctx, id string) (*model.Table, error)
	GetTablesByOutletID(c *fiber.Ctx, outletID string) ([]model.Table, error)
	GetAvailableTables(c *fiber.Ctx, outletID string) ([]model.Table, error)
	CreateTable(c *fiber.Ctx, req *validation.CreateTable) (*model.Table, error)
	UpdateTable(c *fiber.Ctx, id string, req *validation.UpdateTable) (*model.Table, error)
	DeleteTable(c *fiber.Ctx, id string) error
	UpdateTableStatus(c *fiber.Ctx, id string, status string) error
}

type tableService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewTableService(db *gorm.DB, validate *validator.Validate) TableService {
	return &tableService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *tableService) GetTables(c *fiber.Ctx, params *validation.QueryParams, outletID string) ([]model.Table, int64, error) {
	var tables []model.Table
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("Outlet").Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR location LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&tables).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search tables: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&tables).Error; err != nil {
		s.Log.Errorf("Failed to get tables: %+v", err)
		return nil, 0, err
	}

	return tables, totalResults, nil
}

func (s *tableService) GetTableByID(c *fiber.Ctx, id string) (*model.Table, error) {
	table := new(model.Table)
	result := s.DB.WithContext(c.Context()).Preload("Outlet").First(table, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Table not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get table by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return table, nil
}

func (s *tableService) GetTablesByOutletID(c *fiber.Ctx, outletID string) ([]model.Table, error) {
	var tables []model.Table
	result := s.DB.WithContext(c.Context()).Where("outlet_id = ?", outletID).Order("name asc").Find(&tables)
	if result.Error != nil {
		s.Log.Errorf("Failed to get tables by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return tables, nil
}

func (s *tableService) GetAvailableTables(c *fiber.Ctx, outletID string) ([]model.Table, error) {
	var tables []model.Table
	result := s.DB.WithContext(c.Context()).
		Where("outlet_id = ? AND status = ?", outletID, "available").
		Order("name asc").
		Find(&tables)

	if result.Error != nil {
		s.Log.Errorf("Failed to get available tables by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return tables, nil
}

func (s *tableService) CreateTable(c *fiber.Ctx, req *validation.CreateTable) (*model.Table, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	status := "available"
	if req.Status != "" {
		status = req.Status
	}

	table := &model.Table{
		OutletID: uuid.MustParse(req.OutletID),
		Name:     req.Name,
		Location: &req.Location,
		Status:   &status,
		Capacity: req.Capacity,
	}

	result := s.DB.WithContext(c.Context()).Create(table)
	if result.Error != nil {
		s.Log.Errorf("Failed to create table: %+v", result.Error)
	}

	return table, result.Error
}

func (s *tableService) UpdateTable(c *fiber.Ctx, id string, req *validation.UpdateTable) (*model.Table, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	updateFields := make(map[string]interface{})

	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Location != "" {
		updateFields["location"] = req.Location
	}
	if req.Status != "" {
		updateFields["status"] = req.Status
	}
	if req.Capacity > 0 {
		updateFields["capacity"] = req.Capacity
	}

	if len(updateFields) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Table{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		s.Log.Errorf("Failed to update table with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Table not found")
	}

	table, err := s.GetTableByID(c, id)
	if err != nil {
		return nil, err
	}

	return table, nil
}

func (s *tableService) DeleteTable(c *fiber.Ctx, id string) error {
	table := new(model.Table)

	result := s.DB.WithContext(c.Context()).Delete(table, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete table with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Table not found")
	}

	return nil
}

func (s *tableService) UpdateTableStatus(c *fiber.Ctx, id string, status string) error {
	// Validate status
	validStatuses := []string{"available", "occupied", "reserved"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid table status. Must be: available, occupied, or reserved")
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Table{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		s.Log.Errorf("Failed to update table status for ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Table not found")
	}

	return nil
}
