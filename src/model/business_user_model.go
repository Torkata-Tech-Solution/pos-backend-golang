package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BusinessUser struct {
	ID         uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	BusinessID uuid.UUID `gorm:"not null" json:"business_id"`
	UserID     uuid.UUID `gorm:"not null" json:"user_id"`
	Role       string    `gorm:"not null" json:"role"`
	CreatedAt  time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt  time.Time `gorm:"autoCreateTime:milli;autoUpdateTime:milli" json:"-"`

	// Relationships
	Business *Business `gorm:"foreignKey:business_id;references:id" json:"-"`
	User     *User     `gorm:"foreignKey:user_id;references:id" json:"-"`
}

func (businessUser *BusinessUser) BeforeCreate(_ *gorm.DB) error {
	businessUser.ID = uuid.New()
	return nil
}
