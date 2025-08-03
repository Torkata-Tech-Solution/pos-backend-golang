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

type PaymentMethodService interface {
	GetPaymentMethods(c *fiber.Ctx, params *validation.QueryParams) ([]model.PaymentMethod, int64, error)
	GetPaymentMethodByID(c *fiber.Ctx, id string) (*model.PaymentMethod, error)
	GetPaymentMethodsByOutletID(c *fiber.Ctx, outletID string) ([]model.PaymentMethod, error)
	CreatePaymentMethod(c *fiber.Ctx, req *validation.CreatePaymentMethod) (*model.PaymentMethod, error)
	UpdatePaymentMethod(c *fiber.Ctx, id string, req *validation.UpdatePaymentMethod) (*model.PaymentMethod, error)
	DeletePaymentMethod(c *fiber.Ctx, id string) error
}

type paymentMethodService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewPaymentMethodService(db *gorm.DB, validate *validator.Validate) PaymentMethodService {
	return &paymentMethodService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *paymentMethodService) GetPaymentMethods(c *fiber.Ctx, params *validation.QueryParams) ([]model.PaymentMethod, int64, error) {
	var paymentMethods []model.PaymentMethod
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("Outlet").Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("name LIKE ? OR type LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	result := query.Find(&paymentMethods).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to search payment methods: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&paymentMethods).Error; err != nil {
		s.Log.Errorf("Failed to get payment methods: %+v", err)
		return nil, 0, err
	}

	return paymentMethods, totalResults, nil
}

func (s *paymentMethodService) GetPaymentMethodByID(c *fiber.Ctx, id string) (*model.PaymentMethod, error) {
	paymentMethod := new(model.PaymentMethod)
	result := s.DB.WithContext(c.Context()).Preload("Outlet").First(paymentMethod, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Payment method not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get payment method by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return paymentMethod, nil
}

func (s *paymentMethodService) GetPaymentMethodsByOutletID(c *fiber.Ctx, outletID string) ([]model.PaymentMethod, error) {
	var paymentMethods []model.PaymentMethod
	result := s.DB.WithContext(c.Context()).
		Where("outlet_id = ?", outletID).
		Order("name asc").
		Find(&paymentMethods)

	if result.Error != nil {
		s.Log.Errorf("Failed to get payment methods by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return paymentMethods, nil
}

func (s *paymentMethodService) CreatePaymentMethod(c *fiber.Ctx, req *validation.CreatePaymentMethod) (*model.PaymentMethod, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	paymentMethod := &model.PaymentMethod{
		OutletID: uuid.MustParse(req.OutletID),
		Name:     req.Name,
		Type:     req.Type,
	}

	result := s.DB.WithContext(c.Context()).Create(paymentMethod)
	if result.Error != nil {
		s.Log.Errorf("Failed to create payment method: %+v", result.Error)
	}

	return paymentMethod, result.Error
}

func (s *paymentMethodService) UpdatePaymentMethod(c *fiber.Ctx, id string, req *validation.UpdatePaymentMethod) (*model.PaymentMethod, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	if req.Name == "" && req.Type == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	updateFields := make(map[string]interface{})
	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Type != "" {
		updateFields["type"] = req.Type
	}

	result := s.DB.WithContext(c.Context()).Model(&model.PaymentMethod{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		s.Log.Errorf("Failed to update payment method with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Payment method not found")
	}

	paymentMethod, err := s.GetPaymentMethodByID(c, id)
	if err != nil {
		return nil, err
	}

	return paymentMethod, nil
}

func (s *paymentMethodService) DeletePaymentMethod(c *fiber.Ctx, id string) error {
	paymentMethod := new(model.PaymentMethod)

	result := s.DB.WithContext(c.Context()).Delete(paymentMethod, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete payment method with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Payment method not found")
	}

	return nil
}
