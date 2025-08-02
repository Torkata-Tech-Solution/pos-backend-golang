package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SaleItem struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	SaleID    uuid.UUID `gorm:"not null" json:"sale_id"`
	ProductID uuid.UUID `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"type:numeric(10,2);not null" json:"price"`
	Discount  float64   `gorm:"type:numeric(10,2);default:0;not null" json:"discount"`
	Total     float64   `gorm:"type:numeric(10,2);not null" json:"total"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Sale    *Sale    `gorm:"foreignKey:sale_id;references:id" json:"-"`
	Product *Product `gorm:"foreignKey:product_id;references:id" json:"-"`
}

func (saleItem *SaleItem) BeforeCreate(_ *gorm.DB) error {
	saleItem.ID = uuid.New()
	return nil
}
