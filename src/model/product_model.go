package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	Image       *string   `json:"image"`
	Name        string    `gorm:"not null" json:"name"`
	Description *string   `json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CategoryID  uuid.UUID `gorm:"not null" json:"category_id"`
	BusinessID  uuid.UUID `gorm:"not null" json:"business_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt   time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Business  *Business        `gorm:"foreignKey:business_id;references:id" json:"-"`
	Category  *ProductCategory `gorm:"foreignKey:category_id;references:id" json:"-"`
	SaleItems []SaleItem       `gorm:"foreignKey:product_id;references:id" json:"-"`
}

func (product *Product) BeforeCreate(_ *gorm.DB) error {
	product.ID = uuid.New()
	return nil
}
