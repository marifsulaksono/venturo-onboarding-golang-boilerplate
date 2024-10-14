package structs

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Customer) TableName() string {
	return "m_customer"
}

func (CustomerResponse) TableName() string {
	return "m_customer"
}

type (
	Customer struct {
		ID          uuid.UUID      `json:"id" gorm:"primaryKey;type:char(36);not null"`
		UserID      uuid.UUID      `json:"m_user_id" gorm:"not null"`
		User        User           `json:"users" gorm:"foreignKey:UserID;-:migration"`
		Name        string         `json:"name" gorm:"not null"`
		PhoneNumber string         `json:"phone_number" gorm:"not null"`
		Photo       string         `json:"photo_url,omitempty"`
		Address     string         `json:"address" gorm:"not null"`
		CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
		UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
		DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
		CreatedBy   uuid.UUID      `json:"created_by,omitempty" gorm:"type:char(36)"`
		UpdatedBy   uuid.UUID      `json:"updated_by,omitempty" gorm:"type:char(36)"`
		DeletedBy   uuid.UUID      `json:"deleted_by,omitempty" gorm:"type:char(36)"`
	}

	CustomerRequest struct {
		ID          uuid.UUID `json:"id"`
		UserID      uuid.UUID `json:"m_user_id"`
		UserRolesId string    `json:"user_roles_id" gorm:"type:char(36)" validate:"required"`
		Email       string    `json:"email" validate:"required,email"`
		Password    string    `json:"password" validate:"required"`
		Name        string    `json:"name" validate:"required"`
		Address     string    `json:"address" validate:"required"`
		Photo       string    `json:"photo_url"`
		PhoneNumber string    `json:"phone_number" validate:"e164"`
	}

	CustomerResponse struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name" validate:"required"`
		Address     string    `json:"address"`
		Photo       string    `json:"photo_url"`
		PhoneNumber string    `json:"phone_number" validate:"e164"`
	}
)

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}
