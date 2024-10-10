package controllers

import (
	"errors"
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

type CustomerController struct {
	db          *gorm.DB
	model       *models.CustomerModel
	cfg         *config.Config
	imageHelper *helpers.ImageHelper
	assetPath   string
}

func NewCustomerController(db *gorm.DB, model *models.CustomerModel, cfg *config.Config, imageHelper *helpers.ImageHelper, assetPath string) *CustomerController {
	return &CustomerController{db, model, cfg, imageHelper, assetPath}
}

func (uh *CustomerController) Index(c echo.Context) error {
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		per_page = 10
		log.Println("Failed to parse per_page query parameter. Defaulting to 10")
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
		log.Println("Failed to parse page query parameter. Defaulting to 1")
	}
	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "id"
		log.Println("sort query parameter is empty. Defaulting to \"id\"")
	}
	order := c.QueryParam("order")
	if sort == "" {
		sort = "ASC"
		log.Println("sort query parameter is empty. Defaulting to \"ASC\"")
	}

	offset := (page - 1) * per_page
	data, total, err := uh.model.GetAll(per_page, offset, sort, order)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, "")
	}

	pagedData := helpers.PageData(data, total)
	return helpers.Response(c, http.StatusOK, pagedData, "")
}

func (uh *CustomerController) GetById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}
	data, err := uh.model.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.Response(c, http.StatusNotFound, nil, err.Error())
		}
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}
	return helpers.Response(c, http.StatusOK, data, "")
}

func (uh *CustomerController) Create(c echo.Context) error {
	var request structs.CustomerRequest

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, "")
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if request.Photo != "" {
		photo_url, err := uh.imageHelper.Writer(request.Photo, fmt.Sprintf("%s.png", time.Now().Format("20061021545.000000000")))
		if err != nil {
			return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
		}
		request.Photo = photo_url
	}

	data, err := uh.model.Create(&request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, "")
	}

	return helpers.Response(c, http.StatusOK, data, "")
}

func (uh *CustomerController) Update(c echo.Context) error {
	var request structs.CustomerRequest

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, "")
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if request.Photo != "" {
		photo_url, err := uh.imageHelper.Writer(request.Photo, fmt.Sprintf("%s.png", time.Now().Format("20061021545.000000000")))
		if err != nil {
			return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
		}
		request.Photo = photo_url
	}

	data, err := uh.model.Update(&request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, err.Error(), "")
	}

	return helpers.Response(c, http.StatusOK, data, "Customer updated")
}

func (uh *CustomerController) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}
	if err := uh.model.Delete(id); err != nil {
		return helpers.Response(c, http.StatusInternalServerError, err.Error(), "")
	}

	return helpers.Response(c, http.StatusOK, true, "Customer deleted")
}
