package controllers

import (
	"errors"
	"log"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"
	"simple-crud-rnd/structs"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SaleController struct {
	db    *gorm.DB
	model *models.SaleModel
	cfg   *config.Config
}

func NewSaleController(db *gorm.DB, model *models.SaleModel, cfg *config.Config) *SaleController {
	return &SaleController{db, model, cfg}
}

func (uh *SaleController) Index(c echo.Context) error {
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

	offset := (page - 1) * per_page
	data, total, err := uh.model.GetAll(per_page, offset)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, data, "")
	}

	pagedData := helpers.PageData(data, total)
	return helpers.Response(c, http.StatusOK, pagedData, "")
}

func (uh *SaleController) GetById(c echo.Context) error {
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

func (uh *SaleController) Create(c echo.Context) error {
	var request structs.SaleRequest

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	data, err := uh.model.Create(&request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusCreated, data, "")
}

func (uh *SaleController) Update(c echo.Context) error {
	var request structs.SaleRequest

	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	if err := c.Validate(request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, nil, err.Error())
	}

	data, err := uh.model.Update(&request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, nil, err.Error())
	}

	return helpers.Response(c, http.StatusOK, data, "Transaction updated")
}

func (uh *SaleController) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}
	if err := uh.model.Delete(id); err != nil {
		return helpers.Response(c, http.StatusInternalServerError, err.Error(), "")
	}

	return helpers.Response(c, http.StatusOK, true, "Transaction deleted")
}
