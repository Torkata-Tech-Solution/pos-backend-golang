package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory struct {
	ID          uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description *string   `json:"description"`
	BusinessID  uuid.UUID `gorm:"not null" json:"business_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt   time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Business *Business `gorm:"foreignKey:business_id;references:id" json:"-"`
	Products []Product `gorm:"foreignKey:category_id;references:id" json:"-"`
}

func (productCategory *ProductCategory) BeforeCreate(_ *gorm.DB) error {
	productCategory.ID = uuid.New()
	return nil
}
