package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BusinessService interface {
	GetBusinesses(c *fiber.Ctx, params *validation.QueryParams) ([]model.Business, int64, error)
	GetBusinessByID(c *fiber.Ctx, id string) (*model.Business, error)
	CreateBusiness(c *fiber.Ctx, req *validation.CreateBusiness) (*model.Business, error)
	UpdateBusiness(c *fiber.Ctx, id string, req *validation.UpdateBusiness) (*model.Business, error)
	DeleteBusiness(c *fiber.Ctx, id string) error
}

type businessService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewBusinessService(db *gorm.DB, validate *validator.Validate) BusinessService {
	return &businessService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *businessService) GetBusinesses(c *fiber.Ctx, params *validation.QueryParams) ([]model.Business, int64, error) {
	var businesses []model.Business
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

	result := query.Find(&businesses).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search businesses: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&businesses).Error; err != nil {
		s.Log.Errorf("Failed to get businesses: %+v", err)
		return nil, 0, err
	}

	return businesses, totalResults, nil
}

func (s *businessService) GetBusinessByID(c *fiber.Ctx, id string) (*model.Business, error) {
	business := new(model.Business)
	result := s.DB.WithContext(c.Context()).First(business, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Business not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get business by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return business, nil
}

func (s *businessService) CreateBusiness(c *fiber.Ctx, req *validation.CreateBusiness) (*model.Business, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	business := &model.Business{
		Domain:  req.Domain,
		Name:    req.Name,
		Address: req.Address,
		Phone:   &req.Phone,
		Email:   &req.Email,
		Website: &req.Website,
		Logo:    &req.Logo,
	}

	result := s.DB.WithContext(c.Context()).Create(business)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Business with this domain already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create business: %+v", result.Error)
	}

	return business, result.Error
}

func (s *businessService) UpdateBusiness(c *fiber.Ctx, id string, req *validation.UpdateBusiness) (*model.Business, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Name == "" && req.Address == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}
	business := &model.Business{
		Name:    req.Name,
		Address: req.Address,
		Phone:   &req.Phone,
		Email:   &req.Email,
		Website: &req.Website,
		Logo:    &req.Logo,
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Business{}).Where("id = ?", id).Updates(business)
	if result.Error != nil {
		s.Log.Errorf("Failed to update business with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Business not found")
	}

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Business with this domain already exists")
	}

	business, err := s.GetBusinessByID(c, id)
	if err != nil {
		return nil, err
	}

	return business, result.Error
}

func (s *businessService) DeleteBusiness(c *fiber.Ctx, id string) error {

	business := new(model.Business)

	result := s.DB.WithContext(c.Context()).Delete(business, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete business with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Business not found")
	}

	return nil
}
