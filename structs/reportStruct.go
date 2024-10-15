package structs

import "time"

type SalesSummary struct {
	Date            time.Time `json:"date"`
	ProductCategory string    `json:"product_category"`
	ProductName     string    `json:"product_name"`
	TotalAmount     int       `json:"total_amount"`
}
type SummaryOfSales struct {
	CategoryId   int              `json:"category_id"`
	CategoryName string           `json:"category_name"`
	Products     []ProductSummary `json:"products" gorm:"-"`
}

type ProductSummary struct {
	ProductId   int               `json:"product_id"`
	ProductName string            `json:"product_name"`
	Saless      []SalesSummaryNew `json:"sales" gorm:"-"`
	SalessTotal int               `json:"sales_total"`
}

type SalesSummaryNew struct {
	DateSales  time.Time `json:"date_sales"`
	TotalSales int       `json:"total_sales"`
}
