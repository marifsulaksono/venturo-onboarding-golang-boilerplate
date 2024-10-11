package structs

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (Product) TableName() string {
	return "m_product"
}

func (ProductCategory) TableName() string {
	return "m_product_category"
}

func (ProductDetail) TableName() string {
	return "m_product_detail"
}

func (ProductDetailRequest) TableName() string {
	return "m_product_detail"
}

func (ProductResponse) TableName() string {
	return "m_product"
}

type ProductCategory struct {
	ID        int            `json:"id" gorm:"primaryKey,autoIncrement"`
	Name      string         `json:"name" gorm:"not null" validate:"required,max=150"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Product struct {
	ID                uuid.UUID      `json:"id" gorm:"primaryKey;type:char(36);not null"`
	ProductCategoryID int            `json:"product_category_id" gorm:"not null" validate:"required"`
	Name              string         `json:"name" gorm:"not null" validate:"required,max=150"`
	Price             float64        `json:"price" gorm:"not null" validate:"required"`
	Description       string         `json:"description,omitempty" validate:"max=65535"`
	Photo             string         `json:"photo_url,omitempty"`
	IsAvailable       bool           `json:"is_available" gorm:"not null" validate:"required"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CreatedBy         uuid.UUID      `json:"created_by,omitempty" gorm:"type:char(36)"`
	UpdatedBy         uuid.UUID      `json:"updated_by,omitempty" gorm:"type:char(36)"`
	DeletedBy         uuid.UUID      `json:"deleted_by,omitempty" gorm:"type:char(36)"`
}

type ProductResponse struct {
	ID                uuid.UUID       `json:"id" gorm:"primaryKey;type:char(36);not null"`
	ProductCategoryID int             `json:"product_category_id" gorm:"not null" validate:"required"`
	Category          ProductCategory `gorm:"foreignKey:ProductCategoryID;references:ID" json:"category,omitempty"`
	Name              string          `json:"name" gorm:"not null" validate:"required,max=150"`
	Price             float64         `json:"price" gorm:"not null" validate:"required"`
	Description       string          `json:"description,omitempty" validate:"max=65535"`
	Photo             string          `json:"photo_url,omitempty"`
	IsAvailable       bool            `json:"is_available" gorm:"not null" validate:"required"`
	Details           []ProductDetail `json:"details" gorm:"-:migration;foreignKey:ProductID;references:ID"`
	CreatedAt         time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt  `json:"deleted_at" gorm:"index"`
	CreatedBy         uuid.UUID       `json:"created_by,omitempty" gorm:"type:char(36)"`
	UpdatedBy         uuid.UUID       `json:"updated_by,omitempty" gorm:"type:char(36)"`
	DeletedBy         uuid.UUID       `json:"deleted_by,omitempty" gorm:"type:char(36)"`
}

type ProductDetail struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id" gorm:"not null" validate:"required"`
	Type        string    `json:"type" gorm:"not null" validate:"required,oneof='Level' 'Topping'"`
	Description string    `json:"description" validate:"required,max=255"`
	Price       float64   `json:"price" gorm:"not null" validate:"required"`
}

type ProductDetailRequest struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id" validate:"required"`
	Type        string    `json:"type" validate:"required,oneof='Level' 'Topping'"`
	Description string    `json:"description" validate:"required,max=255"`
	Price       float64   `json:"price" validate:"required"`
	IsAdded     bool      `json:"is_added,omitempty"`
	IsUpdated   bool      `json:"is_updated,omitempty"`
	IsDeleted   bool      `json:"is_deleted,omitempty"`
}

type ProductWithDetailCreateOrUpdate struct {
	Product
	Details []ProductDetailRequest `json:"details"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

func (pd *ProductDetail) BeforeCreate(tx *gorm.DB) error {
	pd.ID = uuid.New()
	return nil
}
