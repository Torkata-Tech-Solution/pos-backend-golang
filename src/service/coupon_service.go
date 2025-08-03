package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CouponService interface {
	GetCoupons(c *fiber.Ctx, params *validation.QueryParams, outletID string) ([]model.Coupon, int64, error)
	GetCouponByID(c *fiber.Ctx, id string) (*model.Coupon, error)
	GetCouponsByOutletID(c *fiber.Ctx, outletID string) ([]model.Coupon, error)
	GetCouponByCode(c *fiber.Ctx, code string, outletID string) (*model.Coupon, error)
	GetActiveCoupons(c *fiber.Ctx, outletID string) ([]model.Coupon, error)
	CreateCoupon(c *fiber.Ctx, req *validation.CreateCoupon) (*model.Coupon, error)
	UpdateCoupon(c *fiber.Ctx, id string, req *validation.UpdateCoupon) (*model.Coupon, error)
	DeleteCoupon(c *fiber.Ctx, id string) error
	ValidateCoupon(c *fiber.Ctx, code string, outletID string) (*model.Coupon, error)
	UseCoupon(c *fiber.Ctx, id string) error
}

type couponService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewCouponService(db *gorm.DB, validate *validator.Validate) CouponService {
	return &couponService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *couponService) GetCoupons(c *fiber.Ctx, params *validation.QueryParams, outletID string) ([]model.Coupon, int64, error) {
	var coupons []model.Coupon
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("Outlet").Order("created_at desc")

	// Apply search filter
	if search := params.Search; search != "" {
		query = query.Where("code LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// Apply outlet filter
	if outletID != "" {
		query = query.Where("outlet_id = ?", outletID)
	}

	// Count total results first
	if err := query.Model(&model.Coupon{}).Count(&totalResults).Error; err != nil {
		s.Log.Errorf("Failed to count coupons: %+v", err)
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(params.Limit).Find(&coupons).Error; err != nil {
		s.Log.Errorf("Failed to get coupons: %+v", err)
		return nil, 0, err
	}

	return coupons, totalResults, nil
}

func (s *couponService) GetCouponByID(c *fiber.Ctx, id string) (*model.Coupon, error) {
	coupon := new(model.Coupon)
	result := s.DB.WithContext(c.Context()).Preload("Outlet").First(coupon, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Coupon not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get coupon by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return coupon, nil
}

func (s *couponService) GetCouponsByOutletID(c *fiber.Ctx, outletID string) ([]model.Coupon, error) {
	var coupons []model.Coupon
	result := s.DB.WithContext(c.Context()).
		Where("outlet_id = ?", outletID).
		Order("created_at desc").
		Find(&coupons)

	if result.Error != nil {
		s.Log.Errorf("Failed to get coupons by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return coupons, nil
}

func (s *couponService) GetCouponByCode(c *fiber.Ctx, code string, outletID string) (*model.Coupon, error) {
	coupon := new(model.Coupon)
	result := s.DB.WithContext(c.Context()).
		Where("code = ? AND outlet_id = ?", code, outletID).
		First(coupon)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Coupon not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get coupon by code %s: %+v", code, result.Error)
		return nil, result.Error
	}
	return coupon, nil
}

func (s *couponService) GetActiveCoupons(c *fiber.Ctx, outletID string) ([]model.Coupon, error) {
	var coupons []model.Coupon
	now := time.Now()

	result := s.DB.WithContext(c.Context()).
		Where("outlet_id = ? AND is_active = ? AND start_date <= ? AND end_date >= ? AND used_count < max_uses",
			outletID, true, now, now).
		Order("created_at desc").
		Find(&coupons)

	if result.Error != nil {
		s.Log.Errorf("Failed to get active coupons by outlet ID %s: %+v", outletID, result.Error)
		return nil, result.Error
	}
	return coupons, nil
}

func (s *couponService) CreateCoupon(c *fiber.Ctx, req *validation.CreateCoupon) (*model.Coupon, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Validate end date is after start date
	if req.EndDate.Before(req.StartDate) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "End date must be after start date")
	}

	isActive := true
	if req.IsActive {
		isActive = req.IsActive
	}

	coupon := &model.Coupon{
		OutletID:      uuid.MustParse(req.OutletID),
		Code:          req.Code,
		Description:   &req.Description,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		MaxUses:       req.MaxUses,
		UsedCount:     0,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		IsActive:      isActive,
	}

	result := s.DB.WithContext(c.Context()).Create(coupon)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Coupon with this code already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create coupon: %+v", result.Error)
	}

	return coupon, result.Error
}

func (s *couponService) UpdateCoupon(c *fiber.Ctx, id string, req *validation.UpdateCoupon) (*model.Coupon, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	updateFields := make(map[string]interface{})

	if req.Code != "" {
		updateFields["code"] = req.Code
	}
	if req.Description != "" {
		updateFields["description"] = req.Description
	}
	if req.DiscountType != "" {
		updateFields["discount_type"] = req.DiscountType
	}
	if req.DiscountValue > 0 {
		updateFields["discount_value"] = req.DiscountValue
	}
	if req.MaxUses > 0 {
		updateFields["max_uses"] = req.MaxUses
	}
	if !req.StartDate.IsZero() {
		updateFields["start_date"] = req.StartDate
	}
	if !req.EndDate.IsZero() {
		updateFields["end_date"] = req.EndDate
	}
	updateFields["is_active"] = req.IsActive

	// Validate dates if both are provided
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() && req.EndDate.Before(req.StartDate) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "End date must be after start date")
	}

	if len(updateFields) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "At least one field must be updated")
	}

	result := s.DB.WithContext(c.Context()).Model(&model.Coupon{}).Where("id = ?", id).Updates(updateFields)
	if result.Error != nil {
		s.Log.Errorf("Failed to update coupon with ID %s: %+v", id, result.Error)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Coupon not found")
	}

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Coupon with this code already exists")
	}

	coupon, err := s.GetCouponByID(c, id)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (s *couponService) DeleteCoupon(c *fiber.Ctx, id string) error {
	coupon := new(model.Coupon)

	result := s.DB.WithContext(c.Context()).Delete(coupon, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete coupon with ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Coupon not found")
	}

	return nil
}

func (s *couponService) ValidateCoupon(c *fiber.Ctx, code string, outletID string) (*model.Coupon, error) {
	coupon, err := s.GetCouponByCode(c, code, outletID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// Check if coupon is active
	if !coupon.IsActive {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Coupon is not active")
	}

	// Check if coupon has started
	if now.Before(coupon.StartDate) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Coupon is not yet valid")
	}

	// Check if coupon has expired
	if now.After(coupon.EndDate) {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Coupon has expired")
	}

	// Check if coupon has reached max uses
	if coupon.UsedCount >= coupon.MaxUses {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Coupon has reached maximum usage limit")
	}

	return coupon, nil
}

func (s *couponService) UseCoupon(c *fiber.Ctx, id string) error {
	result := s.DB.WithContext(c.Context()).Model(&model.Coupon{}).
		Where("id = ?", id).
		Update("used_count", gorm.Expr("used_count + 1"))

	if result.Error != nil {
		s.Log.Errorf("Failed to increment coupon usage for ID %s: %+v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Coupon not found")
	}

	return nil
}
