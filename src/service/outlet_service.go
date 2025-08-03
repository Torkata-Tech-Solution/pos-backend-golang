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

type OutletService interface {
	GetOutlets(c *fiber.Ctx, params *validation.QueryParams) ([]model.Outlet, int64, error)
	GetOutletByID(c *fiber.Ctx, id string) (*model.Outlet, error)
	GetOutletsByBusinessID(c *fiber.Ctx, businessID string) ([]model.Outlet, error)
	CreateOutlet(c *fiber.Ctx, req *validation.CreateOutlet) (*model.Outlet, error)
	UpdateOutlet(c *fiber.Ctx, id string, req *validation.UpdateOutlet) (*model.Outlet, error)
	DeleteOutlet(c *fiber.Ctx, id string) error
}

type outletService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewOutletService(db *gorm.DB, validate *validator.Validate) OutletService {
	return &outletService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *outletService) GetOutlets(c *fiber.Ctx, params *validation.QueryParams) ([]model.Outlet, int64, error) {
	var outlets []model.Outlet
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR address LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&outlets).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search outlets: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&outlets).Error; err != nil {
		s.Log.Errorf("Failed to get outlets: %+v", err)
		return nil, 0, err
	}

	return outlets, totalResults, nil
}

func (s *outletService) GetOutletByID(c *fiber.Ctx, id string) (*model.Outlet, error) {
	outlet := new(model.Outlet)
	result := s.DB.WithContext(c.Context()).First(outlet, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Outlet not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get outlet by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return outlet, nil
}

func (s *outletService) GetOutletsByBusinessID(c *fiber.Ctx, businessID string) ([]model.Outlet, error) {
	var outlets []model.Outlet
	result := s.DB.WithContext(c.Context()).Where("business_id = ?", businessID).Find(&outlets)
	if result.Error != nil {
		s.Log.Errorf("Failed to get outlets by business ID %s: %+v", businessID, result.Error)
		return nil, result.Error
	}
	return outlets, nil
}

func (s *outletService) CreateOutlet(c *fiber.Ctx, req *validation.CreateOutlet) (*model.Outlet, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	outlet := &model.Outlet{
		BusinessID: uuid.MustParse(req.BusinessID),
		Name:       req.Name,
		Address:    req.Address,
		Phone:      &req.Phone,
		Email:      &req.Email,
	}

	result := s.DB.WithContext(c.Context()).Create(outlet)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Outlet with this phone or email already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create outlet: %+v", result.Error)
	}

	return outlet, result.Error
}

func (s *outletService) UpdateOutlet(c *fiber.Ctx, id string, req *validation.UpdateOutlet) (*model.Outlet, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Name == "" && req.Address == "" && req.Phone == "" && req.Email == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	outlet := &model.Outlet{
		Name:    req.Name,
		Address: req.Address,
		Phone:   &req.Phone,
		Email:   &req.Email,
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Outlet{}).Where("id = ?", id).Updates(outlet)
	if result.Error != nil {
		s.Log.Errorf("Failed to update outlet with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Outlet not found")
	}

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Outlet with this phone or email already exists")
	}

	outlet, err := s.GetOutletByID(c, id)
	if err != nil {
		return nil, err
	}

	return outlet, result.Error
}

func (s *outletService) DeleteOutlet(c *fiber.Ctx, id string) error {
	outlet := new(model.Outlet)

	result := s.DB.WithContext(c.Context()).Delete(outlet, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete outlet with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Outlet not found")
	}

	return nil
}
