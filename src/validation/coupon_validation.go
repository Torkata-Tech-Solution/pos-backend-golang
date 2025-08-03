package validation

import "time"

// Coupon validations
type CreateCoupon struct {
	OutletID      string    `json:"outlet_id" validate:"required,uuid"`
	Code          string    `json:"code" validate:"required,max=100"`
	Description   string    `json:"description" validate:"omitempty"`
	DiscountType  string    `json:"discount_type" validate:"required,oneof=percentage fixed"`
	DiscountValue float64   `json:"discount_value" validate:"required,min=0"`
	MaxUses       int       `json:"max_uses" validate:"required,min=1"`
	StartDate     time.Time `json:"start_date" validate:"required"`
	EndDate       time.Time `json:"end_date" validate:"required"`
	IsActive      bool      `json:"is_active" validate:"omitempty"`
}

type UpdateCoupon struct {
	Code          string    `json:"code" validate:"omitempty,max=100"`
	Description   string    `json:"description" validate:"omitempty"`
	DiscountType  string    `json:"discount_type" validate:"omitempty,oneof=percentage fixed"`
	DiscountValue float64   `json:"discount_value" validate:"omitempty,min=0"`
	MaxUses       int       `json:"max_uses" validate:"omitempty,min=1"`
	StartDate     time.Time `json:"start_date" validate:"omitempty"`
	EndDate       time.Time `json:"end_date" validate:"omitempty"`
	IsActive      bool      `json:"is_active" validate:"omitempty"`
}

// Sales Coupon validations (junction table)
type CreateSaleCoupon struct {
	SaleID   string `json:"sale_id" validate:"required,uuid"`
	CouponID string `json:"coupon_id" validate:"required,uuid"`
}

