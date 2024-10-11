package models

import (
	"errors"
	"simple-crud-rnd/structs"

	"gorm.io/gorm"
)

type ProductCategoryModel struct {
	db *gorm.DB
}

func NewProductCategoryModel(db *gorm.DB) *ProductCategoryModel {
	return &ProductCategoryModel{
		db: db,
	}
}

func (pcm *ProductCategoryModel) GetAll(limit, offset int) ([]structs.ProductCategory, int64, error) {
	productCategories := []structs.ProductCategory{}
	if err := pcm.db.Select("id", "name", "created_at", "updated_at").
		Limit(limit).Offset(offset).Find(&productCategories).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := pcm.db.Table("m_product_category").Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return productCategories, count, nil
}

func (pcm *ProductCategoryModel) GetById(id int) (structs.ProductCategory, error) {
	productCategory := structs.ProductCategory{}
	err := pcm.db.Select("id", "name", "created_at", "updated_at").
		Where("deleted_at IS NULL").First(&productCategory, id).Error
	return productCategory, err
}

func (pcm *ProductCategoryModel) Create(payload *structs.ProductCategory) (*structs.ProductCategory, error) {
	if err := pcm.db.Create(&payload).Error; err != nil {
		return &structs.ProductCategory{}, err
	}
	return payload, nil
}

func (pcm *ProductCategoryModel) Update(payload *structs.ProductCategory) (*structs.ProductCategory, error) {
	var productCategory structs.ProductCategory
	if err := pcm.db.First(&productCategory, payload.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("data not found")
		}
		return nil, err
	}
	productCategory.Name = payload.Name
	if err := pcm.db.Save(&productCategory).Error; err != nil {
		return nil, err
	}

	return &productCategory, nil
}

func (pm *ProductCategoryModel) Delete(id int) error {
	res := pm.db.Delete(&structs.ProductCategory{}, id)
	if res.RowsAffected == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
