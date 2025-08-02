package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Printer struct {
	ID             uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	OutletID       uuid.UUID `gorm:"not null" json:"outlet_id"`
	Name           string    `gorm:"not null" json:"name"`
	ConnectionType string    `gorm:"not null" json:"connection_type"`
	MacAddress     *string   `json:"mac_address"`
	IPAddress      *string   `json:"ip_address"`
	PaperWidth     *int      `json:"paper_width"`
	DefaultPrinter bool      `gorm:"default:false;not null" json:"default_printer"`
	CreatedAt      time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt      time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Outlet *Outlet `gorm:"foreignKey:outlet_id;references:id" json:"-"`
}

func (printer *Printer) BeforeCreate(_ *gorm.DB) error {
	printer.ID = uuid.New()
	return nil
}
