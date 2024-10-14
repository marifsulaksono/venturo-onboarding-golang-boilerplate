package models

import (
	"errors"
	"simple-crud-rnd/structs"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SaleModel struct {
	db *gorm.DB
}

func NewSaleModel(db *gorm.DB) *SaleModel {
	return &SaleModel{
		db: db,
	}
}

func (sm *SaleModel) GetAll(limit, offset int) ([]structs.Sale, int64, error) {
	sales := []structs.Sale{}
	err := sm.db.Preload("Customer").Select("id, total, customer_id, created_at, updated_at").
		Limit(limit).
		Offset(offset).
		Find(&sales).Error

	if err != nil {
		return nil, 0, err
	}

	var count int64
	if err := sm.db.Table("t_sales").Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return sales, count, nil
}

func (sm *SaleModel) GetById(id uuid.UUID) (structs.Sale, error) {
	sale := structs.Sale{}
	err := sm.db.Select("id", "total", "customer_id", "created_at", "updated_at").
		Where("deleted_at IS NULL").Preload("Customer").Preload("Details.Product").
		Preload("Details.ProductDetail").First(&sale, id).Error
	return sale, err
}

func (sm *SaleModel) Create(payload *structs.SaleRequest) (*structs.Sale, error) {
	sale := structs.Sale{
		CustomerID: payload.CustomerID,
	}

	// Begin transaction
	tx := sm.db.Begin()

	res := tx.Create(&sale).Clauses(clause.Returning{
		Columns: []clause.Column{
			{Name: "id"},
			{Name: "total"},
			{Name: "customer_id"},
		},
	})

	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	if len(payload.Details) > 0 {
		var detailSales []structs.SaleDetail
		for _, details := range payload.Details {
			detail := structs.SaleDetail{
				SaleID:          sale.ID,
				ProductID:       details.ProductID,
				ProductDetailID: details.ProductDetailID,
				Price:           details.Price,
				TotalItem:       details.TotalItem,
			}

			sale.Total += detail.Price * float64(detail.TotalItem)
			detailSales = append(detailSales, detail)
		}

		if err := tx.Create(&detailSales).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Model(&structs.Sale{}).Where("id = ?", sale.ID).Update("total", &sale.Total).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	createdSale, err := sm.GetById(sale.ID)
	if err != nil {
		return nil, err
	}

	return &createdSale, nil
}

func (sm *SaleModel) Update(payload *structs.SaleRequest) (*structs.Sale, error) {
	existSale := structs.Sale{}
	if err := sm.db.Select("id", "customer_id").Where("deleted_at IS NULL").First(&existSale, payload.ID).Error; err != nil {
		return nil, errors.New("data not found")
	}
	// Begin transaction
	tx := sm.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Handle sale details (insert, update, delete)
	for _, detail := range payload.Details {
		if detail.IsAdded {
			// Add new sale detail
			createDetail := structs.SaleDetail{
				SaleID:          payload.ID,
				ProductID:       detail.ProductID,
				ProductDetailID: detail.ProductDetailID,
				TotalItem:       detail.TotalItem,
				Price:           detail.Price,
			}
			if err := tx.Create(&createDetail).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else if detail.IsUpdated {
			// Update existing sale detail
			if err := tx.Model(&structs.SaleDetail{}).Where("id = ?", detail.ID).Updates(&detail).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		payload.Total += detail.Price * float64(detail.TotalItem)
	}

	for _, detail := range payload.DeletedDetails {
		if err := tx.Where("id = ?", detail.ID).Delete(&structs.SaleDetail{}).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Update sale
	existSale.CustomerID = payload.CustomerID
	existSale.Total = payload.Total
	existSale.UpdatedAt = time.Now()
	if err := tx.Model(&structs.Sale{}).Where("id = ?", payload.ID).Updates(&existSale).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	updatedSale, err := sm.GetById(payload.ID)
	if err != nil {
		return nil, err
	}

	return &updatedSale, nil
}

func (sm *SaleModel) Delete(id uuid.UUID) error {
	res := sm.db.Delete(&structs.Sale{}, id)
	if res.RowsAffected == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
