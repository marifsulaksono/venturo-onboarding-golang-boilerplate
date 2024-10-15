package helpers

import (
	"bytes"
	"fmt"
	"simple-crud-rnd/structs"
	"time"

	"github.com/xuri/excelize/v2"
)

func GetPeriod(startDate, endDate string) (map[string]map[string]interface{}, []string, error) {
	// Parse start and end dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, nil, err
	}

	// Add one day to the end date
	end = end.AddDate(0, 0, 1)

	// Initialize interval and period
	periods := make(map[string]map[string]interface{})
	var dates []string
	for date := start; date.Before(end); date = date.AddDate(0, 0, 1) {
		dateStr := date.Format("2006-01-02")
		periods[dateStr] = map[string]interface{}{
			"date_transaction": dateStr,
			"total_sales":      0.0,
		}
		dates = append(dates, dateStr)
	}

	return periods, dates, nil
}

func ReformatSalesReport(sales []structs.Sale, periods []string) []map[string]interface{} {
	salesDetail := make([]map[string]interface{}, 0) // Use a slice to hold the category maps
	totalPerDate := make(map[string]float64)
	total := 0.0

	for _, sale := range sales {
		saleDate := sale.CreatedAt.Format("2006-01-02")
		for _, detail := range sale.Details {
			if detail.Product == nil {
				continue // Skip if no product relation
			}

			categoryId := detail.Product.ProductCategoryID
			categoryName := detail.Product.Category.Name
			productId := detail.ProductID
			productName := detail.Product.Name
			totalSales := float64(detail.Price) * float64(detail.TotalItem)

			// Check if the category already exists in salesDetail
			var categoryMap map[string]interface{}
			var found bool

			for _, category := range salesDetail {
				if category["category_id"] == categoryId {
					categoryMap = category
					found = true
					break
				}
			}

			if !found {
				// If category does not exist, create a new one
				categoryMap = map[string]interface{}{
					"category_id":    categoryId,
					"category_name":  categoryName,
					"category_total": 0.0,
					"products":       map[string]map[string]interface{}{},
				}
				salesDetail = append(salesDetail, categoryMap)
			}

			categoryMap["category_total"] = categoryMap["category_total"].(float64) + totalSales

			if categoryMap["products"].(map[string]map[string]interface{})[productId.String()] == nil {
				categoryMap["products"].(map[string]map[string]interface{})[productId.String()] = map[string]interface{}{
					"product_id":         productId,
					"product_name":       productName,
					"transactions":       initializePeriod(periods),
					"transactions_total": 0.0,
				}
			}

			product := categoryMap["products"].(map[string]map[string]interface{})[productId.String()]
			product["transactions_total"] = product["transactions_total"].(float64) + totalSales

			transaction := product["transactions"].(map[string]map[string]interface{})
			transaction[saleDate]["total_sales"] = transaction[saleDate]["total_sales"].(float64) + totalSales

			// Add to total per date and grand total
			totalPerDate[saleDate] += totalSales
			total += totalSales
		}
	}

	// The return type is now a slice of maps, each representing a category
	return salesDetail
}

func initializePeriod(periods []string) map[string]map[string]interface{} {
	transactions := make(map[string]map[string]interface{})
	for _, period := range periods {
		transactions[period] = map[string]interface{}{
			"date_transaction": period,
			"total_sales":      0.0,
		}
	}
	return transactions
}

func ExportSalesReport(formatedReport []map[string]interface{}, dates []string, buf *bytes.Buffer) (string, error) {
	// Create a new Excel file
	f := excelize.NewFile()

	sheet := "SalesReport"
	f.NewSheet(sheet)

	// Set column widths for better readability (adjust as needed)
	f.SetColWidth(sheet, "A", "Z", 15)

	// Row 1 - Merge "Menu" and "Periode" + Dates
	f.MergeCell(sheet, "A1", "A2")
	f.SetCellValue(sheet, "A1", "Menu")

	// Merge period across the date length
	dateEndCol := string('C' + len(dates) - 1)
	f.MergeCell(sheet, "C1", dateEndCol+"1")
	f.SetCellValue(sheet, "C1", "Periode")

	// Fill dates in row 2
	for i, date := range dates {
		col := string('C' + i)
		f.SetCellValue(sheet, col+"2", date)
	}

	// Merge Total column and set values
	totalCol := string('C' + len(dates))
	f.MergeCell(sheet, totalCol+"1", totalCol+"2")
	f.SetCellValue(sheet, totalCol+"1", "Total")

	row := 3
	for _, category := range formatedReport {
		categoryName := category["category_name"].(string)
		categoryTotal := category["category_total"].(float64)
		products := category["products"].(map[string]map[string]interface{})

		// Merge category_name based on the number of products
		productCount := len(products)
		endRow := row + productCount - 1
		f.MergeCell(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", endRow))
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), categoryName)

		// Merge category_total
		f.MergeCell(sheet, fmt.Sprintf("%s%d", totalCol, row), fmt.Sprintf("%s%d", totalCol, endRow))
		f.SetCellValue(sheet, fmt.Sprintf("%s%d", totalCol, row), categoryTotal)

		for _, product := range products {
			productName := product["product_name"].(string)
			transactionTotal := product["transactions_total"].(float64)
			transactions := product["transactions"].(map[string]map[string]interface{})

			// Set product_name
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), productName)

			// Set total sales for each date
			for i, date := range dates {
				col := string('C' + i)
				totalSales := transactions[date]["total_sales"].(float64)
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), totalSales)
			}

			// Set transactions_total
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", totalCol, row), transactionTotal)
			row++
		}
	}

	// Set Grand Total row
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Grand Total")

	// Calculate sum of total_sales and transactions_total
	for i := 0; i <= len(dates); i++ {
		col := string('C' + i)
		sumRange := fmt.Sprintf("%s3:%s%d", col, col, row-1)
		f.SetCellFormula(sheet, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("SUM(%s)", sumRange))
	}

	// Set sum for transactions_total
	totalSalesRange := fmt.Sprintf("%s3:%s%d", totalCol, totalCol, row-1)
	f.SetCellFormula(sheet, fmt.Sprintf("%s%d", totalCol, row), fmt.Sprintf("SUM(%s)", totalSalesRange))

	err := f.Write(buf)
	if err != nil {
		return "", err
	}

	return "SalesReport.xlsx", nil
}
