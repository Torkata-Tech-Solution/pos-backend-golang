package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID            uuid.UUID  `gorm:"primaryKey;not null" json:"id"`
	UserID        *uuid.UUID `json:"user_id"`
	Name          string     `gorm:"not null" json:"name"`
	Email         string     `gorm:"uniqueIndex;not null" json:"email"`
	Phone         *string    `json:"phone"`
	Address       *string    `gorm:"type:text" json:"address"`
	LoyaltyPoints int        `gorm:"default:0;not null" json:"loyalty_points"`
	OutletID      uuid.UUID  `gorm:"not null" json:"outlet_id"`
	CreatedAt     time.Time  `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt     time.Time  `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Outlet *Outlet `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	User   *User   `gorm:"foreignKey:user_id;references:id" json:"-"`
	Sales  []Sale  `gorm:"foreignKey:customer_id;references:id" json:"-"`
}

func (customer *Customer) BeforeCreate(_ *gorm.DB) error {
	customer.ID = uuid.New()
	return nil
}
