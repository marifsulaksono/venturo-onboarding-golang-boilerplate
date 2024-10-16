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
	salesDetail := make([]map[string]interface{}, 0)
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
					"products":       []map[string]interface{}{},
				}
				salesDetail = append(salesDetail, categoryMap)
			}

			categoryMap["category_total"] = categoryMap["category_total"].(float64) + totalSales

			// Check if the product already exists in the category's products list
			products := categoryMap["products"].([]map[string]interface{})
			var productMap map[string]interface{}
			productFound := false

			for _, product := range products {
				if product["product_id"] == productId {
					productMap = product
					productFound = true
					break
				}
			}

			if !productFound {
				// If product does not exist, create a new one
				productMap = map[string]interface{}{
					"product_id":         productId,
					"product_name":       productName,
					"transactions":       initializePeriod(periods),
					"transactions_total": 0.0,
				}
				products = append(products, productMap)
				categoryMap["products"] = products // Update the products array
			}

			// Update product's transaction totals
			productMap["transactions_total"] = productMap["transactions_total"].(float64) + totalSales

			transaction := productMap["transactions"].(map[string]map[string]interface{})
			transaction[saleDate]["total_sales"] = transaction[saleDate]["total_sales"].(float64) + totalSales

			// Add to total per date and grand total
			totalPerDate[saleDate] += totalSales
			total += totalSales
		}
	}
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
	index, err := f.NewSheet(sheet)
	if err != nil {
		return "", err
	}

	f.SetActiveSheet(index)
	f.SetColWidth(sheet, "A", "Z", 15)

	// Set Header
	f.MergeCell(sheet, "A1", "A2")
	f.SetCellValue(sheet, "A1", "Menu")

	dateEndCol := string('B' + len(dates) - 1)
	f.MergeCell(sheet, "B1", dateEndCol+"1")
	f.SetCellValue(sheet, "B1", "Periode")

	for i, date := range dates {
		col := string('B' + i)
		f.SetCellValue(sheet, col+"2", date)
	}

	totalCol := string('B' + len(dates))
	f.MergeCell(sheet, totalCol+"1", totalCol+"2")
	f.SetCellValue(sheet, totalCol+"1", "Total")

	// Set Body
	row := 3
	for _, category := range formatedReport {
		categoryName := category["category_name"].(string)
		categoryTotal := category["category_total"].(float64)
		products := category["products"].([]map[string]interface{})

		productCount := len(products)

		// Merge category_name in the first column, horizontally (A), and not across multiple rows.
		f.MergeCell(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row))
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), categoryName)

		// Set category_total in the total column
		f.MergeCell(sheet, fmt.Sprintf("%s%d", totalCol, row), fmt.Sprintf("%s%d", totalCol, row+productCount-1))
		f.SetCellValue(sheet, fmt.Sprintf("%d%s", row, totalCol), fmt.Sprintf("Rp. %.f", categoryTotal))

		// Loop through products and add them under the category
		for _, product := range products {
			productName := product["product_name"].(string)
			transactionTotal := product["transactions_total"].(float64)
			transactions := product["transactions"].(map[string]map[string]interface{})

			// Set product_name in the next row (row+1), with no merging needed for products.
			row++ // Move to the next row for the next product
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), productName)

			// Set total sales for each date
			for i, date := range dates {
				col := string('B' + i)
				totalSales := 0.0
				if transaction, ok := transactions[date]; ok {
					totalSales = transaction["total_sales"].(float64)
				}
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), totalSales)
			}

			// Set transactions_total
			f.SetCellValue(sheet, fmt.Sprintf("%s%d", totalCol, row), transactionTotal)
		}
		row++
	}

	// Set Grand Total row
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Grand Total")

	// Calculate sum of total_sales and transactions_total
	for i := 0; i < len(dates); i++ {
		col := string('B' + i)
		sumRange := fmt.Sprintf("%s3:%s%d", col, col, row-1) // Ensure correct range
		f.SetCellFormula(sheet, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("SUM(%s)", sumRange))
	}

	// Set sum for transactions_total
	totalSalesRange := fmt.Sprintf("%s3:%s%d", totalCol, totalCol, row-1) // Correct range for transaction totals
	f.SetCellFormula(sheet, fmt.Sprintf("%s%d", totalCol, row), fmt.Sprintf("SUM(%s)", totalSalesRange))

	if err := f.Write(buf); err != nil {
		return "", err
	}

	return "SalesReport.xlsx", nil
}
