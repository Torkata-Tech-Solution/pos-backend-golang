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

type SaleCouponService interface {
	GetSaleCoupons(c *fiber.Ctx, params *validation.QueryParams) ([]model.SaleCoupon, int64, error)
	GetSaleCouponByID(c *fiber.Ctx, id string) (*model.SaleCoupon, error)
	GetSaleCouponsBySaleID(c *fiber.Ctx, saleID string) ([]model.SaleCoupon, error)
	GetSaleCouponsByCouponID(c *fiber.Ctx, couponID string) ([]model.SaleCoupon, error)
	CreateSaleCoupon(c *fiber.Ctx, req *validation.CreateSaleCoupon) (*model.SaleCoupon, error)
	DeleteSaleCoupon(c *fiber.Ctx, id string) error
	DeleteSaleCouponBySaleAndCoupon(c *fiber.Ctx, saleID, couponID string) error
}

type saleCouponService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewSaleCouponService(db *gorm.DB, validate *validator.Validate) SaleCouponService {
	return &saleCouponService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *saleCouponService) GetSaleCoupons(c *fiber.Ctx, params *validation.QueryParams) ([]model.SaleCoupon, int64, error) {
	var saleCoupons []model.SaleCoupon
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Preload("Sale").Preload("Coupon").Order("created_at desc")

	result := query.Model(&model.SaleCoupon{}).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to count sale coupons: %+v", result.Error)
		return nil, 0, result.Error
	}

	if err := query.Offset(offset).Limit(params.Limit).Find(&saleCoupons).Error; err != nil {
		s.Log.Errorf("Failed to get sale coupons: %+v", err)
		return nil, 0, err
	}

	return saleCoupons, totalResults, nil
}

func (s *saleCouponService) GetSaleCouponByID(c *fiber.Ctx, id string) (*model.SaleCoupon, error) {
	saleCoupon := new(model.SaleCoupon)
	result := s.DB.WithContext(c.Context()).Preload("Sale").Preload("Coupon").First(saleCoupon, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Sale coupon not found")
	}
	if result.Error != nil {
		s.Log.Errorf("Failed to get sale coupon by ID %s: %+v", id, result.Error)
		return nil, result.Error
	}
	return saleCoupon, nil
}

func (s *saleCouponService) GetSaleCouponsBySaleID(c *fiber.Ctx, saleID string) ([]model.SaleCoupon, error) {
	var saleCoupons []model.SaleCoupon
	result := s.DB.WithContext(c.Context()).Preload("Coupon").Where("sale_id = ?", saleID).Find(&saleCoupons)
	if result.Error != nil {
		s.Log.Errorf("Failed to get sale coupons by sale ID %s: %+v", saleID, result.Error)
		return nil, result.Error
	}
	return saleCoupons, nil
}

func (s *saleCouponService) GetSaleCouponsByCouponID(c *fiber.Ctx, couponID string) ([]model.SaleCoupon, error) {
	var saleCoupons []model.SaleCoupon
	result := s.DB.WithContext(c.Context()).Preload("Sale").Where("coupon_id = ?", couponID).Find(&saleCoupons)
	if result.Error != nil {
		s.Log.Errorf("Failed to get sale coupons by coupon ID %s: %+v", couponID, result.Error)
		return nil, result.Error
	}
	return saleCoupons, nil
}

func (s *saleCouponService) CreateSaleCoupon(c *fiber.Ctx, req *validation.CreateSaleCoupon) (*model.SaleCoupon, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Check if sale exists
	var sale model.Sale
	if err := s.DB.WithContext(c.Context()).First(&sale, "id = ?", req.SaleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Sale not found")
		}
		return nil, err
	}

	// Check if coupon exists
	var coupon model.Coupon
	if err := s.DB.WithContext(c.Context()).First(&coupon, "id = ?", req.CouponID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Coupon not found")
		}
		return nil, err
	}

	// Check if relationship already exists
	var existingSaleCoupon model.SaleCoupon
	result := s.DB.WithContext(c.Context()).Where("sale_id = ? AND coupon_id = ?", req.SaleID, req.CouponID).First(&existingSaleCoupon)
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusConflict, "Coupon is already applied to this sale")
	}

	saleCoupon := &model.SaleCoupon{
		SaleID:   uuid.MustParse(req.SaleID),
		CouponID: uuid.MustParse(req.CouponID),
	}

	result = s.DB.WithContext(c.Context()).Create(saleCoupon)
	if result.Error != nil {
		s.Log.Errorf("Failed to create sale coupon: %+v", result.Error)
		return nil, result.Error
	}

	return saleCoupon, nil
}

func (s *saleCouponService) DeleteSaleCoupon(c *fiber.Ctx, id string) error {
	result := s.DB.WithContext(c.Context()).Delete(&model.SaleCoupon{}, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete sale coupon: %+v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Sale coupon not found")
	}
	return nil
}

func (s *saleCouponService) DeleteSaleCouponBySaleAndCoupon(c *fiber.Ctx, saleID, couponID string) error {
	result := s.DB.WithContext(c.Context()).Where("sale_id = ? AND coupon_id = ?", saleID, couponID).Delete(&model.SaleCoupon{})
	if result.Error != nil {
		s.Log.Errorf("Failed to delete sale coupon by sale and coupon: %+v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Sale coupon relationship not found")
	}
	return nil
}
