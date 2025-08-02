package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Sale struct {
	ID              uuid.UUID  `gorm:"primaryKey;not null" json:"id"`
	OutletID        uuid.UUID  `gorm:"not null" json:"outlet_id"`
	OutletStaffID   uuid.UUID  `gorm:"not null" json:"outlet_staff_id"`
	CustomerID      *uuid.UUID `json:"customer_id"`
	PaymentMethodID *uuid.UUID `json:"payment_method_id"`
	TableID         uuid.UUID  `gorm:"not null" json:"table_id"`
	InvoiceNumber   string     `gorm:"uniqueIndex;not null" json:"invoice_number"`
	Total           float64    `gorm:"type:numeric(10,2);not null" json:"total"`
	Discount        float64    `gorm:"type:numeric(10,2);default:0;not null" json:"discount"`
	Tax             float64    `gorm:"type:numeric(10,2);default:0;not null" json:"tax"`
	GrandTotal      float64    `gorm:"type:numeric(10,2);not null" json:"grand_total"`
	Status          string     `gorm:"not null" json:"status"`
	SaleDate        time.Time  `gorm:"default:CURRENT_TIMESTAMP;not null" json:"sale_date"`
	Note            *string    `gorm:"type:text" json:"note"`
	CreatedAt       time.Time  `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt       time.Time  `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Outlet        *Outlet        `gorm:"foreignKey:outlet_id;references:id" json:"-"`
	OutletStaff   *OutletStaff   `gorm:"foreignKey:outlet_staff_id;references:id" json:"-"`
	Customer      *Customer      `gorm:"foreignKey:customer_id;references:id" json:"-"`
	PaymentMethod *PaymentMethod `gorm:"foreignKey:payment_method_id;references:id" json:"-"`
	Table         *Table         `gorm:"foreignKey:table_id;references:id" json:"-"`
	SaleItems     []SaleItem     `gorm:"foreignKey:sale_id;references:id" json:"-"`
	SaleCoupons   []SaleCoupon   `gorm:"foreignKey:sale_id;references:id" json:"-"`
}

func (sale *Sale) BeforeCreate(_ *gorm.DB) error {
	sale.ID = uuid.New()
	return nil
}
