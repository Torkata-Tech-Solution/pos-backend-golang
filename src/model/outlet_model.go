package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Outlet struct {
	ID         uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	BusinessID uuid.UUID `gorm:"not null" json:"business_id"`
	Name       string    `gorm:"not null" json:"name"`
	Address    string    `gorm:"not null" json:"address"`
	Phone      *string   `gorm:"uniqueIndex" json:"phone"`
	Email      *string   `gorm:"uniqueIndex" json:"email"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt  time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Business       *Business       `gorm:"foreignKey:business_id;references:id" json:"-"`
	OutletStaff    []OutletStaff   `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Customers      []Customer      `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	PaymentMethods []PaymentMethod `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Tables         []Table         `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Sales          []Sale          `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Settings       []Setting       `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Coupons        []Coupon        `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Printers       []Printer       `gorm:"foreignKey:outlet_id;references:id" json:"-"`
}

func (outlet *Outlet) BeforeCreate(_ *gorm.DB) error {
	outlet.ID = uuid.New()
	return nil
}
