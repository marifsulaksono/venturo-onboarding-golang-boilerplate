package models

import (
	"errors"
	"fmt"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerModel struct {
	db *gorm.DB
}

func NewCustomerModel(db *gorm.DB) *CustomerModel {
	return &CustomerModel{
		db: db,
	}
}

func (cm *CustomerModel) GetAll(limit, offset int, sort, order string) ([]structs.CustomerResponse, int64, error) {
	customers := []structs.CustomerResponse{}
	if err := cm.db.Model(&structs.Customer{}).Find(&customers).Limit(limit).Offset(offset).Order(fmt.Sprintf("%s %s", sort, order)).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := cm.db.Table("m_customer").Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return customers, count, nil
}

func (cm *CustomerModel) GetAllCount() (int64, error) {
	var count int64
	if err := cm.db.Model(&structs.Customer{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (cm *CustomerModel) GetById(id uuid.UUID) (structs.CustomerResponse, error) {
	customer := structs.CustomerResponse{}
	err := cm.db.Model(&structs.Customer{}).First(&customer, id).Error
	return customer, err
}

func (cm *CustomerModel) Create(payload *structs.CustomerRequest) (structs.CustomerResponse, error) {
	var user structs.User
	var createdCustomer structs.CustomerResponse

	hashedPassword, pwErr := helpers.PasswordHash(payload.Password)
	if pwErr != nil {
		return createdCustomer, pwErr
	}

	user = structs.User{
		Name:            payload.Name,
		Email:           payload.Email,
		Password:        hashedPassword,
		Photo:           payload.Photo,
		PhoneNumber:     payload.PhoneNumber,
		UserRolesId:     payload.UserRolesId,
		UpdatedSecurity: time.Now(),
	}

	tx := cm.db.Begin()

	res := tx.Create(&user).Clauses(clause.Returning{
		Columns: []clause.Column{
			{Name: "id"},
			{Name: "name"},
			{Name: "email"},
			{Name: "phone_number"},
			{Name: "photo"},
			{Name: "updated_security"},
		},
	})
	if res.Error != nil {
		tx.Rollback()
		return createdCustomer, res.Error
	}

	customer := structs.Customer{
		UserID:      user.ID,
		Name:        payload.Name,
		Address:     payload.Address,
		Photo:       payload.Photo,
		PhoneNumber: payload.PhoneNumber,
	}

	cusRes := tx.Create(&customer).Clauses(clause.Returning{
		Columns: []clause.Column{
			{Name: "id"},
			{Name: "name"},
			{Name: "address"},
			{Name: "photo"},
			{Name: "phone_number"},
		},
	})

	if cusRes.Error != nil {
		tx.Rollback()
		return createdCustomer, cusRes.Error
	}

	tx.Commit()

	createdCustomer = structs.CustomerResponse{
		ID:          customer.ID,
		Name:        customer.Name,
		Address:     customer.Address,
		Photo:       customer.Photo,
		PhoneNumber: customer.PhoneNumber,
	}

	return createdCustomer, nil
}

func (cm *CustomerModel) Update(payload *structs.CustomerRequest) (structs.CustomerResponse, error) {
	var updatedCustomer structs.CustomerResponse

	hashedPassword, pwErr := helpers.PasswordHash(payload.Password)
	if pwErr != nil {
		return updatedCustomer, pwErr
	}

	userPayload := structs.User{
		Name:            payload.Name,
		Email:           payload.Email,
		Password:        hashedPassword,
		Photo:           payload.Photo,
		PhoneNumber:     payload.PhoneNumber,
		UserRolesId:     payload.UserRolesId,
		UpdatedSecurity: time.Now(),
	}

	tx := cm.db.Begin()

	res := tx.Model(&structs.User{}).Where("id = ?", payload.UserID).Updates(&userPayload)
	if res.RowsAffected == 0 {
		tx.Rollback()
		return updatedCustomer, errors.New("no rows updated")
	}

	customerPayload := structs.Customer{
		ID:          payload.ID,
		Name:        payload.Name,
		Address:     payload.Address,
		PhoneNumber: payload.PhoneNumber,
		Photo:       payload.Photo,
	}

	upRes := tx.Model(&structs.Customer{}).Where("id = ?", customerPayload.ID).Updates(&customerPayload)
	if upRes.RowsAffected == 0 {
		tx.Rollback()
		return updatedCustomer, errors.New("no rows updated")
	}

	updatedCustomer, err := cm.GetById(customerPayload.ID)
	if err != nil {
		return updatedCustomer, err
	}

	tx.Commit()

	return updatedCustomer, nil
}

func (cm *CustomerModel) Delete(id uuid.UUID) error {
	res := cm.db.Delete(&structs.Customer{}, id)
	if res.RowsAffected == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
