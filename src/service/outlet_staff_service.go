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

type OutletStaffService interface {
	GetOutletStaff(c *fiber.Ctx, params *validation.QueryParams, outletID string, role string) ([]model.OutletStaff, int64, error)
	GetOutletStaffByID(c *fiber.Ctx, id string) (*model.OutletStaff, error)
	GetOutletStaffByOutletID(c *fiber.Ctx, outletID string) ([]model.OutletStaff, error)
	GetOutletStaffByUsername(c *fiber.Ctx, username string) (*model.OutletStaff, error)
	CreateOutletStaff(c *fiber.Ctx, req *validation.CreateOutletStaff) (*model.OutletStaff, error)
	UpdateOutletStaff(c *fiber.Ctx, id string, req *validation.UpdateOutletStaff) (*model.OutletStaff, error)
	DeleteOutletStaff(c *fiber.Ctx, id string) error
	ChangePassword(c *fiber.Ctx, id string, req *validation.ChangePasswordOutletStaff) error
}

type outletStaffService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewOutletStaffService(db *gorm.DB, validate *validator.Validate) OutletStaffService {
	return &outletStaffService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *outletStaffService) GetOutletStaff(c *fiber.Ctx, params *validation.QueryParams, outletID string, role string) ([]model.OutletStaff, int64, error) {
	var staff []model.OutletStaff
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("Outlet").Order("created_at desc")

	if search := params.Search; search != "" {
		query = query.Where("name ILIKE ? OR username ILIKE ? OR role ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	result := query.Model(&model.OutletStaff{}).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to count outlet staff: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&staff).Error; err != nil {
		s.Log.Errorf("Failed to get outlet staff: %+v", err)
		return nil, 0, err
	}

	return staff, totalResults, nil
}

func (s *outletStaffService) GetOutletStaffByID(c *fiber.Ctx, id string) (*model.OutletStaff, error) {
	staff := new(model.OutletStaff)
	result := s.DB.WithContext(c.Context()).Preload("Outlet").First(staff, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Outlet staff not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get outlet staff by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return staff, nil
}

func (s *outletStaffService) GetOutletStaffByOutletID(c *fiber.Ctx, outletID string) ([]model.OutletStaff, error) {
	var staff []model.OutletStaff
	result := s.DB.WithContext(c.Context()).Where("outlet_id = ?", outletID).Find(&staff)
	if result.Error != nil {
		s.Log.Errorf("Failed to get outlet staff by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return staff, nil
}

func (s *outletStaffService) GetOutletStaffByUsername(c *fiber.Ctx, username string) (*model.OutletStaff, error) {
	staff := new(model.OutletStaff)
	result := s.DB.WithContext(c.Context()).Preload("Outlet").Where("username = ?", username).First(staff)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Outlet staff not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get outlet staff by username %s: %+v", username, result.Error)
		return nil, result.Error
	}
	return staff, nil
}

func (s *outletStaffService) CreateOutletStaff(c *fiber.Ctx, req *validation.CreateOutletStaff) (*model.OutletStaff, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Check if outlet exists
	var outlet model.Outlet
	if err := s.DB.WithContext(c.Context()).First(&outlet, "id = ?", req.OutletID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Outlet not found")
		}
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.Log.Errorf("Failed to hash password: %+v", err)
		return nil, err
	}

	staff := &model.OutletStaff{
		OutletID: uuid.MustParse(req.OutletID),
		Name:     req.Name,
		Password: hashedPassword,
		Role:     req.Role,
	}

	// Set username if provided
	if req.Username != "" {
		staff.Username = &req.Username
	}

	result := s.DB.WithContext(c.Context()).Create(staff)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Username already exists")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to create outlet staff: %+v", result.Error)
		return nil, result.Error
	}

	return staff, nil
}

func (s *outletStaffService) UpdateOutletStaff(c *fiber.Ctx, id string, req *validation.UpdateOutletStaff) (*model.OutletStaff, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	staff := new(model.OutletStaff)
	result := s.DB.WithContext(c.Context()).First(staff, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Outlet staff not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get outlet staff by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	// Update fields
	if req.Name != "" {
		staff.Name = req.Name
	}
	if req.Username != "" {
		staff.Username = &req.Username
	}
	if req.Role != "" {
		staff.Role = req.Role
	}

	result = s.DB.WithContext(c.Context()).Save(staff)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Username already exists")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to update outlet staff: %+v", result.Error)
		return nil, result.Error
	}

	return staff, nil
}

func (s *outletStaffService) ChangePassword(c *fiber.Ctx, id string, req *validation.ChangePasswordOutletStaff) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	staff := new(model.OutletStaff)
	result := s.DB.WithContext(c.Context()).First(staff, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "Outlet staff not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get outlet staff by ID %s: %+v", id, result.Error)
		return result.Error
	}

	// Verify current password
	if !utils.CheckPasswordHash(req.CurrentPassword, staff.Password) {
		return fiber.NewError(fiber.StatusBadRequest, "Current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		s.Log.Errorf("Failed to hash new password: %+v", err)
		return err
	}

	staff.Password = hashedPassword
	result = s.DB.WithContext(c.Context()).Save(staff)
	if result.Error != nil {
		s.Log.Errorf("Failed to update outlet staff password: %+v", result.Error)
		return result.Error
	}

	return nil
}

func (s *outletStaffService) DeleteOutletStaff(c *fiber.Ctx, id string) error {
	result := s.DB.WithContext(c.Context()).Delete(&model.OutletStaff{}, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete outlet staff: %+v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Outlet staff not found")
	}
	return nil
}
