package models

import (
	"errors"
	"simple-crud-rnd/structs"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductModel struct {
	db *gorm.DB
}

func NewProductModel(db *gorm.DB) *ProductModel {
	return &ProductModel{
		db: db,
	}
}

func (pm *ProductModel) GetAll(limit, offset int) ([]structs.Product, int64, error) {
	products := []structs.Product{}
	if err := pm.db.Select("id", "product_category_id", "name", "price", "description", "photo", "is_available", "created_at", "updated_at").
		Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := pm.db.Table("m_product").Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return products, count, nil
}

func (pm *ProductModel) GetById(id uuid.UUID) (structs.ProductResponse, error) {
	product := structs.ProductResponse{}
	err := pm.db.Select("id", "product_category_id", "name", "price", "description", "photo", "is_available", "created_at", "updated_at").
		Where("deleted_at IS NULL").Preload("Category").Preload("Details").First(&product, id).Error
	return product, err
}

func (pm *ProductModel) Create(payload *structs.ProductWithDetailCreateOrUpdate) (*structs.Product, error) {
	product := structs.Product{
		ProductCategoryID: payload.ProductCategoryID,
		Name:              payload.Name,
		Price:             payload.Price,
		Description:       payload.Description,
		Photo:             payload.Photo,
		IsAvailable:       payload.IsAvailable,
	}

	// Begin transaction
	tx := pm.db.Begin()

	res := tx.Create(&product).Clauses(clause.Returning{
		Columns: []clause.Column{
			{Name: "id"},
			{Name: "product_category_id"},
			{Name: "name"},
			{Name: "price"},
			{Name: "description"},
			{Name: "photo"},
			{Name: "is_available"},
		},
	})

	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	if len(payload.Details) > 0 {
		var detailProducts []structs.ProductDetail
		for _, details := range payload.Details {
			detail := structs.ProductDetail{
				ProductID:   product.ID,
				Type:        details.Type,
				Price:       details.Price,
				Description: details.Description,
			}

			detailProducts = append(detailProducts, detail)
		}

		if err := tx.Create(&detailProducts).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (pm *ProductModel) Update(payload *structs.ProductWithDetailCreateOrUpdate) (*structs.ProductResponse, error) {
	existProduct := structs.Product{}
	if err := pm.db.Select("id", "name").Where("deleted_at IS NULL").First(&existProduct, payload.ID).Error; err != nil {
		return nil, errors.New("data not found")
	}
	// Begin transaction
	tx := pm.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update product
	if err := tx.Model(&structs.Product{}).Where("id = ?", payload.ID).Updates(&payload.Product).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Handle product details (insert, update, delete)
	for _, detail := range payload.Details {
		if detail.IsAdded {
			// Add new product detail
			createDetail := structs.ProductDetail{
				ProductID:   payload.ID,
				Type:        detail.Type,
				Description: detail.Description,
				Price:       detail.Price,
			}
			if err := tx.Create(&createDetail).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else if detail.IsUpdated {
			// Update existing product detail
			if err := tx.Model(&structs.ProductDetail{}).Where("id = ?", detail.ID).Updates(&detail).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	for _, detail := range payload.DeletedDetails {
		if err := tx.Where("id = ?", detail.ID).Delete(&structs.ProductDetail{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	product, err := pm.GetById(payload.ID)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (pm *ProductModel) Delete(id uuid.UUID) error {
	res := pm.db.Delete(&structs.Product{}, id)
	if res.RowsAffected == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
