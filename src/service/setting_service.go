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

type SettingService interface {
	GetSettings(c *fiber.Ctx, params *validation.QueryParams, outletID string) ([]model.Setting, int64, error)
	GetSettingByID(c *fiber.Ctx, id string) (*model.Setting, error)
	GetSettingsByOutletID(c *fiber.Ctx, outletID string) ([]model.Setting, error)
	GetSettingByKey(c *fiber.Ctx, outletID string, key string) (*model.Setting, error)
	CreateSetting(c *fiber.Ctx, req *validation.CreateSetting) (*model.Setting, error)
	UpdateSetting(c *fiber.Ctx, id string, req *validation.UpdateSetting) (*model.Setting, error)
	DeleteSetting(c *fiber.Ctx, id string) error
	UpsertSetting(c *fiber.Ctx, outletID string, key string, value string) (*model.Setting, error)
}

type settingService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewSettingService(db *gorm.DB, validate *validator.Validate) SettingService {
	return &settingService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *settingService) GetSettings(c *fiber.Ctx, params *validation.QueryParams, outletID string) ([]model.Setting, int64, error) {
	var settings []model.Setting
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("Outlet").Order("key asc")

	if search := params.Search; search != "" {
		query = query.Where("key LIKE ? OR value LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&settings).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search settings: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&settings).Error; err != nil {
		s.Log.Errorf("Failed to get settings: %+v", err)
		return nil, 0, err
	}

	return settings, totalResults, nil
}

func (s *settingService) GetSettingByID(c *fiber.Ctx, id string) (*model.Setting, error) {
	setting := new(model.Setting)
	result := s.DB.WithContext(c.Context()).Preload("Outlet").First(setting, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Setting not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get setting by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return setting, nil
}

func (s *settingService) GetSettingsByOutletID(c *fiber.Ctx, outletID string) ([]model.Setting, error) {
	var settings []model.Setting
	result := s.DB.WithContext(c.Context()).
		Where("outlet_id = ?", outletID).
		Order("key asc").
		Find(&settings)

	if result.Error != nil {
		s.Log.Errorf("Failed to get settings by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return settings, nil
}

func (s *settingService) GetSettingByKey(c *fiber.Ctx, outletID string, key string) (*model.Setting, error) {
	setting := new(model.Setting)
	result := s.DB.WithContext(c.Context()).
		Where("outlet_id = ? AND key = ?", outletID, key).
		First(setting)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Setting not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get setting by key %s: %+v", key, result.Error)
		return nil, result.Error
	}
	return setting, nil
}

func (s *settingService) CreateSetting(c *fiber.Ctx, req *validation.CreateSetting) (*model.Setting, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Check if setting with same key already exists for this outlet
	existingSetting, err := s.GetSettingByKey(c, req.OutletID, req.Key)
	if err == nil && existingSetting != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "Setting with this key already exists for this outlet")
	}

	setting := &model.Setting{
		OutletID: uuid.MustParse(req.OutletID),
		Key:      req.Key,
		Value:    req.Value,
	}

	result := s.DB.WithContext(c.Context()).Create(setting)
	if result.Error != nil {
		s.Log.Errorf("Failed to create setting: %+v", result.Error)
	}

	return setting, result.Error
}

func (s *settingService) UpdateSetting(c *fiber.Ctx, id string, req *validation.UpdateSetting) (*model.Setting, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Value == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Value must be provided")
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Setting{}).
		Where("id = ?", id).
		Update("value", req.Value)

	if result.Error != nil {
		s.Log.Errorf("Failed to update setting with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Setting not found")
	}

	setting, err := s.GetSettingByID(c, id)
	if err != nil {
		return nil, err
	}

	return setting, nil
}

func (s *settingService) DeleteSetting(c *fiber.Ctx, id string) error {
	setting := new(model.Setting)

	result := s.DB.WithContext(c.Context()).Delete(setting, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete setting with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Setting not found")
	}

	return nil
}

func (s *settingService) UpsertSetting(c *fiber.Ctx, outletID string, key string, value string) (*model.Setting, error) {
	// Try to get existing setting
	existingSetting, err := s.GetSettingByKey(c, outletID, key)

	if err != nil && err.Error() != "Setting not found" {
		return nil, err
	}

	if existingSetting != nil {
		// Update existing setting
		updateReq := &validation.UpdateSetting{
			Value: value,
		}
		return s.UpdateSetting(c, existingSetting.ID.String(), updateReq)
	}

	// Create new setting
	createReq := &validation.CreateSetting{
		OutletID: outletID,
		Key:      key,
		Value:    value,
	}
	return s.CreateSetting(c, createReq)
}
