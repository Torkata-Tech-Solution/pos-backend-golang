package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OutletStaff struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	OutletID  uuid.UUID `gorm:"not null" json:"outlet_id"`
	Name      string    `gorm:"not null" json:"name"`
	Username  *string   `gorm:"uniqueIndex" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"not null" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Outlet *Outlet `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Sales  []Sale  `gorm:"foreignKey:outlet_staff_id;references:id" json:"-"`
}

func (outletStaff *OutletStaff) BeforeCreate(_ *gorm.DB) error {
	outletStaff.ID = uuid.New()
	return nil
}
