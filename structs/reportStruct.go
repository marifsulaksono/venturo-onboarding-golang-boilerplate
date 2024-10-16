package structs

import "time"

type SalesSummary struct {
	Date            time.Time `json:"date"`
	ProductCategory string    `json:"product_category"`
	ProductName     string    `json:"product_name"`
	TotalAmount     int       `json:"total_amount"`
}

type SalesSummaryNew struct {
	DateSales  time.Time `json:"date_sales"`
	TotalSales int       `json:"total_sales"`
}
