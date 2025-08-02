package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Coupon struct {
	ID            uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	OutletID      uuid.UUID `gorm:"not null" json:"outlet_id"`
	Code          string    `gorm:"uniqueIndex;not null" json:"code"`
	Description   *string   `gorm:"type:text" json:"description"`
	DiscountType  string    `gorm:"not null" json:"discount_type"`
	DiscountValue float64   `gorm:"type:numeric(10,2);not null" json:"discount_value"`
	MaxUses       int       `gorm:"default:1;not null" json:"max_uses"`
	UsedCount     int       `gorm:"default:0;not null" json:"used_count"`
	StartDate     time.Time `gorm:"not null" json:"start_date"`
	EndDate       time.Time `gorm:"not null" json:"end_date"`
	IsActive      bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt     time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt     time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Outlet      *Outlet      `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	SaleCoupons []SaleCoupon `gorm:"foreignKey:coupon_id;references:id" json:"-"`
}

func (coupon *Coupon) BeforeCreate(_ *gorm.DB) error {
	coupon.ID = uuid.New()
	return nil
}
