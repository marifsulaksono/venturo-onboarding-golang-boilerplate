package controllers

import (
	"fmt"
	"log"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	"simple-crud-rnd/structs"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController struct {
	db          *gorm.DB
	model       *models.UserModel
	cfg         *config.Config
	imageHelper *helpers.ImageHelper
	assetPath   string
}

func NewUserController(db *gorm.DB, model *models.UserModel, cfg *config.Config, imageHelper *helpers.ImageHelper, assetPath string) *UserController {
	return &UserController{db, model, cfg, imageHelper, assetPath}
}

func (uh *UserController) Index(c echo.Context) error {
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		per_page = 10
		log.Printf("Failed to parse per_page query parameter. Defaulting to %d", per_page)
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
		log.Printf("Failed to parse page query parameter. Defaulting to %d", page)
	}

	offset := (page - 1) * per_page
	data, total, err := uh.model.GetAll(per_page, offset)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, err.Error())
	}
	pagedData := helpers.PageData(data, total)
	return helpers.Response(c, http.StatusOK, pagedData, "")
}

func (uh *UserController) GetById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	data, err := uh.model.GetById(id)
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}
	data.Photo, err = uh.imageHelper.Read(data.Photo)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}
	return helpers.Response(c, http.StatusOK, data, "")
}

func (uh *UserController) Create(c echo.Context) error {
	var request structs.UserRequest

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	photo_url, err := uh.imageHelper.Writer(request.Photo, fmt.Sprintf("%s.png", time.Now().Format("20061021545.000000000")))
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	request.Photo = photo_url
	data, err := uh.model.Create(&request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, "")
	}

	return helpers.Response(c, http.StatusCreated, data, "")
}

func (uh *UserController) Update(c echo.Context) error {
	var request structs.User

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	photo_url, err := uh.imageHelper.Writer(request.Photo, fmt.Sprintf("%s.png", time.Now().Format("20061021545.000000000")))
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	request.Photo = photo_url
	data, err := uh.model.Update(request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusOK, data, "User updated")
}

func (uh *UserController) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}
	if err := uh.model.Delete(id); err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusOK, true, "User deleted")
}
