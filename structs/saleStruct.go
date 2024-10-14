package structs

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Sale) TableName() string {
	return "t_sales"
}

func (SaleDetail) TableName() string {
	return "t_sales_detail"
}

type Sale struct {
	ID         uuid.UUID        `json:"id" gorm:"primaryKey;type:char(36);not null"`
	Total      float64          `json:"amount" gorm:"not null;default:0"`
	CustomerID uuid.UUID        `json:"m_customer_id" gorm:"type:char(36);not null" validate:"required"`
	Customer   CustomerResponse `json:"customer" gorm:"-:migration;foreignKey:CustomerID;references:ID"`
	Details    []SaleDetail     `json:"details" gorm:"-:migration;foreignKey:SaleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt  time.Time        `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time        `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt   `json:"deleted_at,omitempty"`
}

type SaleDetail struct {
	ID              uuid.UUID     `json:"id" gorm:"primaryKey;type:char(36);not null"`
	SaleID          uuid.UUID     `json:"t_sales_id" gorm:"foreignKey:SaleID;references:ID;type:char(36);not null"`
	ProductID       uuid.UUID     `json:"m_product_id" gorm:"type:char(36);not null" validate:"required"`
	Product         Product       `json:"product" gorm:"-:migration;foreignKey:ProductID;references:ID"`
	ProductDetailID uuid.UUID     `json:"m_product_detail_id" gorm:"type:char(36);not null" validate:"required"`
	ProductDetail   ProductDetail `json:"product_detail" gorm:"-:migration;foreignKey:ProductDetailID;references:ID"`
	TotalItem       int           `json:"total_item" gorm:"not null" validate:"required"`
	Price           float64       `json:"price" gorm:"not null" validate:"required"`
	IsAdded         bool          `json:"is_added,omitempty" gorm:"-"`
	IsUpdated       bool          `json:"is_updated,omitempty" gorm:"-"`
}

type SaleRequest struct {
	ID             uuid.UUID    `json:"id" gorm:"primaryKey;type:char(36);not null"`
	Total          float64      `json:"amount" gorm:"not null;default:0"`
	CustomerID     uuid.UUID    `json:"m_customer_id" gorm:"type:char(36);not null" validate:"required"`
	Details        []SaleDetail `json:"product_detail" validate:"required"`
	DeletedDetails []SaleDetail `json:"deleted_detail,omitempty"`
}

func (s *Sale) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}

func (sd *SaleDetail) BeforeCreate(tx *gorm.DB) error {
	sd.ID = uuid.New()
	return nil
}
