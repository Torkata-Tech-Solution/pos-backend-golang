package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	OutletID  uuid.UUID `gorm:"not null" json:"outlet_id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null" json:"type"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Outlet *Outlet `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Sales  []Sale  `gorm:"foreignKey:payment_method_id;references:id" json:"-"`
}

func (paymentMethod *PaymentMethod) BeforeCreate(_ *gorm.DB) error {
	paymentMethod.ID = uuid.New()
	return nil
}
