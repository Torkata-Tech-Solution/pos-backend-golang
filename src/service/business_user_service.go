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

type BusinessUserService interface {
	GetBusinessUsers(c *fiber.Ctx, params *validation.QueryParams) ([]model.BusinessUser, int64, error)
	GetBusinessUserByID(c *fiber.Ctx, id string) (*model.BusinessUser, error)
	GetBusinessUsersByBusinessID(c *fiber.Ctx, businessID string) ([]model.BusinessUser, error)
	GetBusinessUsersByUserID(c *fiber.Ctx, userID string) ([]model.BusinessUser, error)
	GetBusinessUserByBusinessAndUser(c *fiber.Ctx, businessID, userID string) (*model.BusinessUser, error)
	CreateBusinessUser(c *fiber.Ctx, req *validation.CreateBusinessUser) (*model.BusinessUser, error)
	UpdateBusinessUser(c *fiber.Ctx, id string, req *validation.UpdateBusinessUser) (*model.BusinessUser, error)
	DeleteBusinessUser(c *fiber.Ctx, id string) error
}

type businessUserService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewBusinessUserService(db *gorm.DB, validate *validator.Validate) BusinessUserService {
	return &businessUserService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *businessUserService) GetBusinessUsers(c *fiber.Ctx, params *validation.QueryParams) ([]model.BusinessUser, int64, error) {
	var businessUsers []model.BusinessUser
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("Business").Preload("User").Order("created_at desc")

	if search := params.Search; search != "" {
		query = query.Joins("LEFT JOIN businesses ON business_users.business_id = businesses.id").
			Joins("LEFT JOIN users ON business_users.user_id = users.id").
			Where("businesses.name ILIKE ? OR users.name ILIKE ? OR business_users.role ILIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	result := query.Model(&model.BusinessUser{}).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to count business users: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&businessUsers).Error; err != nil {
		s.Log.Errorf("Failed to get business users: %+v", err)
		return nil, 0, err
	}

	return businessUsers, totalResults, nil
}

func (s *businessUserService) GetBusinessUserByID(c *fiber.Ctx, id string) (*model.BusinessUser, error) {
	businessUser := new(model.BusinessUser)
	result := s.DB.WithContext(c.Context()).Preload("Business").Preload("User").First(businessUser, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Business user not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get business user by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return businessUser, nil
}

func (s *businessUserService) GetBusinessUsersByBusinessID(c *fiber.Ctx, businessID string) ([]model.BusinessUser, error) {
	var businessUsers []model.BusinessUser
	result := s.DB.WithContext(c.Context()).Preload("User").Where("business_id = ?", businessID).Find(&businessUsers)
	if result.Error != nil {
		s.Log.Errorf("Failed to get business users by business ID %s: %+v", businessID, result.Error)
		return nil, result.Error
	}
	return businessUsers, nil
}

func (s *businessUserService) GetBusinessUsersByUserID(c *fiber.Ctx, userID string) ([]model.BusinessUser, error) {
	var businessUsers []model.BusinessUser
	result := s.DB.WithContext(c.Context()).Preload("Business").Where("user_id = ?", userID).Find(&businessUsers)
	if result.Error != nil {
		s.Log.Errorf("Failed to get business users by user ID %s: %+v", userID, result.Error)
		return nil, result.Error
	}
	return businessUsers, nil
}

func (s *businessUserService) GetBusinessUserByBusinessAndUser(c *fiber.Ctx, businessID, userID string) (*model.BusinessUser, error) {
	businessUser := new(model.BusinessUser)
	result := s.DB.WithContext(c.Context()).Preload("Business").Preload("User").
		Where("business_id = ? AND user_id = ?", businessID, userID).First(businessUser)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Business user relationship not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get business user by business ID %s and user ID %s: %+v", businessID, userID, result.Error)
		return nil, result.Error
	}
	return businessUser, nil
}

func (s *businessUserService) CreateBusinessUser(c *fiber.Ctx, req *validation.CreateBusinessUser) (*model.BusinessUser, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Check if business exists
	var business model.Business
	if err := s.DB.WithContext(c.Context()).First(&business, "id = ?", req.BusinessID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Business not found")
		}
		return nil, err
	}

	// Check if user exists
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", req.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "User not found")
		}
		return nil, err
	}

	// Check if relationship already exists
	var existingBusinessUser model.BusinessUser
	result := s.DB.WithContext(c.Context()).Where("business_id = ? AND user_id = ?", req.BusinessID, req.UserID).First(&existingBusinessUser)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusConflict, "User is already associated with this business")
	}

	businessUser := &model.BusinessUser{
		BusinessID: uuid.MustParse(req.BusinessID),
		UserID:     uuid.MustParse(req.UserID),
		Role:       req.Role,
	}

	result = s.DB.WithContext(c.Context()).Create(businessUser)
	if result.Error != nil {
		s.Log.Errorf("Failed to create business user: %+v", result.Error)
		return nil, result.Error
	}

	return businessUser, nil
}

func (s *businessUserService) UpdateBusinessUser(c *fiber.Ctx, id string, req *validation.UpdateBusinessUser) (*model.BusinessUser, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	businessUser := new(model.BusinessUser)
	result := s.DB.WithContext(c.Context()).First(businessUser, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Business user not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get business user by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	// Update role
	if req.Role != "" {
		businessUser.Role = req.Role
	}

	result = s.DB.WithContext(c.Context()).Save(businessUser)
	if result.Error != nil {
		s.Log.Errorf("Failed to update business user: %+v", result.Error)
		return nil, result.Error
	}

	return businessUser, nil
}

func (s *businessUserService) DeleteBusinessUser(c *fiber.Ctx, id string) error {
	result := s.DB.WithContext(c.Context()).Delete(&model.BusinessUser{}, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete business user: %+v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Business user not found")
	}
	return nil
}
