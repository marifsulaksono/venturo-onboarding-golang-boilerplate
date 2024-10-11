package models

import (
	"errors"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/structs"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

func (um *UserModel) GetAll(limit, offset int) ([]structs.User, int64, error) {
	users := []structs.User{}
	if err := um.db.Select("id", "name", "email", "phone_number", "photo", "user_roles_id", "updated_security", "created_at", "updated_at").
		Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := um.db.Table("m_user").Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (um *UserModel) GetById(id uuid.UUID) (structs.User, error) {
	user := structs.User{}
	err := um.db.Select("id", "name", "email", "phone_number", "photo", "user_roles_id", "updated_security", "created_at", "updated_at").
		Where("deleted_at IS NULL").First(&user, id).Error
	return user, err
}

func (um *UserModel) Create(payload *structs.UserRequest) (structs.User, error) {
	var user structs.User
	hashedPassword, pwErr := helpers.PasswordHash(payload.Password)
	if pwErr != nil {
		return user, pwErr
	}
	user = structs.User{
		Name:            payload.Name,
		Email:           payload.Email,
		PhoneNumber:     payload.PhoneNumber,
		Password:        hashedPassword,
		Photo:           payload.Photo,
		UserRolesId:     payload.UserRolesId,
		UpdatedSecurity: time.Now(),
	}

	res := um.db.Create(&user).Clauses(clause.Returning{
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
		return user, res.Error
	}

	return user, nil
}

func (um *UserModel) Update(payload structs.User) (structs.User, error) {
	user := structs.User{ID: payload.ID}
	hashedPassword, pwErr := helpers.PasswordHash(payload.Password)
	if pwErr != nil {
		return user, pwErr
	}
	payload.Password = hashedPassword
	res := um.db.Model(&user).Clauses(clause.Returning{}).Updates(&payload)
	if res.RowsAffected == 0 {
		return user, errors.New("no rows updated")
	}
	return user, nil
}

func (um *UserModel) Delete(id uuid.UUID) error {
	res := um.db.Delete(&structs.User{}, id)
	if res.RowsAffected == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}
