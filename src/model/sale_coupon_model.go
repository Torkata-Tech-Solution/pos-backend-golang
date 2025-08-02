package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SaleCoupon struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	SaleID    uuid.UUID `gorm:"not null" json:"sale_id"`
	CouponID  uuid.UUID `gorm:"not null" json:"coupon_id"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Sale   *Sale   `gorm:"foreignKey:sale_id;references:id" json:"-"`
	Coupon *Coupon `gorm:"foreignKey:coupon_id;references:id" json:"-"`
}

func (saleCoupon *SaleCoupon) BeforeCreate(_ *gorm.DB) error {
	saleCoupon.ID = uuid.New()
	return nil
}
