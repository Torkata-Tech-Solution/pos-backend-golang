package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Business struct {
	ID        uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	Domain    string    `gorm:"uniqueIndex;not null" json:"domain"`
	Name      string    `gorm:"not null" json:"name"`
	Address   string    `gorm:"not null" json:"address"`
	Phone     *string   `gorm:"uniqueIndex" json:"phone"`
	Email     *string   `gorm:"uniqueIndex" json:"email"`
	Website   *string   `json:"website"`
	Logo      *string   `json:"logo"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	BusinessUsers     []BusinessUser    `gorm:"foreignKey:business_id;references:id" json:"-"`
	Outlets           []Outlet          `gorm:"foreignKey:business_id;references:id" json:"-"`
	ProductCategories []ProductCategory `gorm:"foreignKey:business_id;references:id" json:"-"`
	Products          []Product         `gorm:"foreignKey:business_id;references:id" json:"-"`
}

func (business *Business) BeforeCreate(_ *gorm.DB) error {
	business.ID = uuid.New()
	return nil
}
