package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Setting struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	OutletID  uuid.UUID `gorm:"not null" json:"outlet_id"`
	Key       string    `gorm:"not null" json:"key"`
	Value     string    `gorm:"type:text;not null" json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Outlet *Outlet `gorm:"foreignKey:outlet_id;references:id" json:"-"`
}

func (setting *Setting) BeforeCreate(_ *gorm.DB) error {
	setting.ID = uuid.New()
	return nil
}
