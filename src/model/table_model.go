package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Table struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	OutletID  uuid.UUID `gorm:"not null" json:"outlet_id"`
	Name      string    `gorm:"not null" json:"name"`
	Location  *string   `json:"location"`
	Status    *string   `json:"status"`
	Capacity  int       `gorm:"not null" json:"capacity"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Outlet *Outlet `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	Sales  []Sale  `gorm:"foreignKey:table_id;references:id" json:"-"`
}

func (table *Table) BeforeCreate(_ *gorm.DB) error {
	table.ID = uuid.New()
	return nil
}
