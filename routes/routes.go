package routes

import (
	"fmt"
	"log"
	"simple-crud-rnd/config"
	"simple-crud-rnd/controllers"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type APIVersionOne struct {
	e          *echo.Echo
	db         *gorm.DB
	cfg        *config.Config
	api        *echo.Group
	assetsPath string
}

func InitVersionOne(e *echo.Echo, db *gorm.DB, cfg *config.Config) *APIVersionOne {
	return &APIVersionOne{
		e,
		db,
		cfg,
		e.Group("/api/v1"),
		fmt.Sprintf("%s/%s", cfg.HTTP.Domain, cfg.HTTP.AssetEndpoint),
	}
}

func (av *APIVersionOne) UserAndAuth() {
	userModel := models.NewUserModel(av.db)
	imageHelper, err := helpers.NewImageHelper(av.cfg.AssetStorage.Path, "profile_photos")
	if err != nil {
		log.Fatal("Failed to initiate an image helper:", err)
	}

	userController := controllers.NewUserController(av.db, userModel, av.cfg, imageHelper, av.assetsPath)
	// authController := controllers.NewAuthController(av.db, userModel, av.cfg)

	auth := av.api.Group("/auth")
	// auth.POST("/login", authController.Login)
	auth.POST("/signup", userController.Create)

	// user := av.api.Group("/users", echojwt.WithConfig(av.cfg.JWT.Config))
	user := av.api.Group("/users")

	user.GET("", userController.Index)
	user.POST("", userController.Create)
	user.GET("/:id", userController.GetById)
	user.PUT("", userController.Update)
	user.DELETE("/:id", userController.Delete)
}

func (av *APIVersionOne) Customer() {
	customerModel := models.NewCustomerModel(av.db)
	imageHelper, err := helpers.NewImageHelper(av.cfg.AssetStorage.Path, "profile_photos")
	if err != nil {
		log.Fatal("Failed to initiate an image helper:", err)
	}
	customerController := controllers.NewCustomerController(av.db, customerModel, av.cfg, imageHelper, av.assetsPath)

	// customer := av.api.Group("/customers", echojwt.WithConfig(av.cfg.JWT.Config))
	customer := av.api.Group("/customers")
	customer.GET("", customerController.Index)
	customer.POST("", customerController.Create)
	customer.GET("/:id", customerController.GetById)
	customer.PUT("", customerController.Update)
	customer.DELETE("/:id", customerController.Delete)
}

func (av *APIVersionOne) ProductCategory() {
	ProductCategoryModel := models.NewProductCategoryModel(av.db)
	productCategoryController := controllers.NewProductCategoryController(av.db, ProductCategoryModel, av.cfg)

	// ProductCategory := av.api.Group("/productCategories", echojwt.WithConfig(av.cfg.JWT.Config))
	productCategory := av.api.Group("/productCategories")
	productCategory.GET("", productCategoryController.Index)
	productCategory.POST("", productCategoryController.Create)
	productCategory.GET("/:id", productCategoryController.GetById)
	productCategory.PUT("", productCategoryController.Update)
	productCategory.DELETE("/:id", productCategoryController.Delete)
}

func (av *APIVersionOne) Product() {
	productModel := models.NewProductModel(av.db)
	imageHelper, err := helpers.NewImageHelper(av.cfg.AssetStorage.Path, "product_photos")
	if err != nil {
		log.Fatal("Failed to initiate an image helper:", err)
	}
	productController := controllers.NewProductController(av.db, productModel, av.cfg, imageHelper, av.assetsPath)

	// product := av.api.Group("/products", echojwt.WithConfig(av.cfg.JWT.Config))
	product := av.api.Group("/products")
	product.GET("", productController.Index)
	product.POST("", productController.Create)
	product.GET("/:id", productController.GetById)
	product.PUT("", productController.Update)
	product.DELETE("/:id", productController.Delete)
}
