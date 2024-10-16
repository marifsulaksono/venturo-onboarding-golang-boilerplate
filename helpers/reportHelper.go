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

// func ReformatSalesReport(sales []structs.Sale, periods []string) []map[string]interface{} {
// 	salesDetail := make([]map[string]interface{}, 0) // Use a slice to hold the category maps
// 	totalPerDate := make(map[string]float64)
// 	total := 0.0

// 	for _, sale := range sales {
// 		saleDate := sale.CreatedAt.Format("2006-01-02")
// 		for _, detail := range sale.Details {
// 			if detail.Product == nil {
// 				continue // Skip if no product relation
// 			}

// 			categoryId := detail.Product.ProductCategoryID
// 			categoryName := detail.Product.Category.Name
// 			productId := detail.ProductID
// 			productName := detail.Product.Name
// 			totalSales := float64(detail.Price) * float64(detail.TotalItem)

// 			// Check if the category already exists in salesDetail
// 			var categoryMap map[string]interface{}
// 			var found bool

// 			for _, category := range salesDetail {
// 				if category["category_id"] == categoryId {
// 					categoryMap = category
// 					found = true
// 					break
// 				}
// 			}

// 			if !found {
// 				// If category does not exist, create a new one
// 				categoryMap = map[string]interface{}{
// 					"category_id":    categoryId,
// 					"category_name":  categoryName,
// 					"category_total": 0.0,
// 					"products":       map[string]map[string]interface{}{},
// 				}
// 				salesDetail = append(salesDetail, categoryMap)
// 			}

// 			categoryMap["category_total"] = categoryMap["category_total"].(float64) + totalSales

// 			if categoryMap["products"].(map[string]map[string]interface{})[productId.String()] == nil {
// 				categoryMap["products"].(map[string]map[string]interface{})[productId.String()] = map[string]interface{}{
// 					"product_id":         productId,
// 					"product_name":       productName,
// 					"transactions":       initializePeriod(periods),
// 					"transactions_total": 0.0,
// 				}
// 			}

// 			product := categoryMap["products"].(map[string]map[string]interface{})[productId.String()]
// 			product["transactions_total"] = product["transactions_total"].(float64) + totalSales

// 			transaction := product["transactions"].(map[string]map[string]interface{})
// 			transaction[saleDate]["total_sales"] = transaction[saleDate]["total_sales"].(float64) + totalSales

// 			// Add to total per date and grand total
// 			totalPerDate[saleDate] += totalSales
// 			total += totalSales
// 		}
// 	}

// 	// The return type is now a slice of maps, each representing a category
// 	return salesDetail
// }

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

// func ExportSalesReport(formatedReport []map[string]interface{}, dates []string, buf *bytes.Buffer) (string, error) {
// 	// Create a new Excel file
// 	f := excelize.NewFile()

// 	sheet := "SalesReport"
// 	f.NewSheet(sheet)

// 	// Set column widths for better readability (adjust as needed)
// 	f.SetColWidth(sheet, "A", "Z", 15)

// 	// Row 1 - Merge "Menu" and "Periode" + Dates
// 	f.MergeCell(sheet, "A1", "A2")
// 	f.SetCellValue(sheet, "A1", "Menu")

// 	// Merge period across the date length
// 	dateEndCol := string('B' + len(dates) - 1)
// 	f.MergeCell(sheet, "B1", dateEndCol+"1")
// 	f.SetCellValue(sheet, "B1", "Periode")

// 	// Fill dates in row 2
// 	for i, date := range dates {
// 		col := string('B' + i)
// 		f.SetCellValue(sheet, col+"2", date)
// 	}

// 	// Merge Total column and set values
// 	totalCol := string('B' + len(dates))
// 	f.MergeCell(sheet, totalCol+"1", totalCol+"2")
// 	f.SetCellValue(sheet, totalCol+"1", "Total")

// 	row := 3
// 	for _, category := range formatedReport {
// 		categoryName := category["category_name"].(string)
// 		categoryTotal := category["category_total"].(float64)
// 		products := category["products"].(map[string]map[string]interface{})

// 		// Merge category_name based on the number of products
// 		productCount := len(products)
// 		endRow := row + productCount - 1
// 		f.MergeCell(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", endRow))
// 		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), categoryName)

// 		// Merge category_total
// 		f.MergeCell(sheet, fmt.Sprintf("%s%d", totalCol, row), fmt.Sprintf("%s%d", totalCol, endRow))
// 		f.SetCellValue(sheet, fmt.Sprintf("%s%d", totalCol, row), categoryTotal)

// 		for _, product := range products {
// 			productName := product["product_name"].(string)
// 			transactionTotal := product["transactions_total"].(float64)
// 			transactions := product["transactions"].(map[string]map[string]interface{})

// 			// Set product_name
// 			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), productName)

// 			// Set total sales for each date
// 			for i, date := range dates {
// 				col := string('B' + i)
// 				totalSales := transactions[date]["total_sales"].(float64)
// 				f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("Rp. %.f", totalSales))
// 			}

// 			// Set transactions_total
// 			f.SetCellValue(sheet, fmt.Sprintf("%s%d", totalCol, row), fmt.Sprintf("Rp. %.f", transactionTotal))
// 			row++
// 		}
// 	}

// 	// Set Grand Total row
// 	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Grand Total")

// 	// Calculate sum of total_sales and transactions_total
// 	for i := 0; i <= len(dates); i++ {
// 		col := string('B' + i)
// 		sumRange := fmt.Sprintf("%s3:%s%d", col, col, row-1)
// 		f.SetCellFormula(sheet, fmt.Sprintf("%s%d", col, row), fmt.Sprintf("SUM(%s)", sumRange))
// 	}

// 	// Set sum for transactions_total
// 	totalSalesRange := fmt.Sprintf("%s3:%s%d", totalCol, totalCol, row-1)
// 	f.SetCellFormula(sheet, fmt.Sprintf("%s%d", totalCol, row), fmt.Sprintf("SUM(%s)", totalSalesRange))

// 	if err := f.Write(buf); err != nil {
// 		return "", err
// 	}

// 	return "SalesReport.xlsx", nil
// }

// ============================ V2

func ExportSalesReport(formattedReport []map[string]interface{}, dates []string, buf *bytes.Buffer) (string, error) {
	// Create a new Excel file
	f := excelize.NewFile()

	sheet := "SalesReport"
	f.SetSheetName("Sheet1", sheet)
	index, err := f.GetSheetIndex(sheet)
	if err != nil {
		return "", err
	}
	f.SetActiveSheet(index)

	// set cell format
	formatRupiah := "Rp#,##0.00"
	rupiahStyle, err := f.NewStyle(&excelize.Style{
		NumFmt:       4,
		CustomNumFmt: &formatRupiah,
	})
	if err != nil {
		fmt.Println("Error creating style:", err)
		return "", err
	}

	// Set Header
	f.SetColWidth(sheet, "A", "Z", 15)
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
	for _, category := range formattedReport {
		categoryName := category["category_name"].(string)
		categoryTotal := category["category_total"].(float64)
		products := category["products"].([]map[string]interface{})

		f.MergeCell(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("G%d", row))
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), categoryName)

		categoryTotalCell := fmt.Sprintf("%s%d", totalCol, row)
		f.SetCellValue(sheet, categoryTotalCell, categoryTotal)
		setRupiahCurrencyFormat(f, sheet, categoryTotalCell, rupiahStyle)

		for _, product := range products {
			row++ // Move to the next row for the next product
			productName := product["product_name"].(string)
			transactionTotal := product["transactions_total"].(float64)
			transactions := product["transactions"].(map[string]map[string]interface{})

			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), productName)

			// Set total sales for each date
			for i, date := range dates {
				col := string('B' + i)
				totalSales := 0.0
				if transaction, ok := transactions[date]; ok {
					totalSales = transaction["total_sales"].(float64)
				}
				totalSalesCell := fmt.Sprintf("%s%d", col, row)
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", col, row), totalSales)
				setRupiahCurrencyFormat(f, sheet, totalSalesCell, rupiahStyle)
			}

			// Set transactions_total
			transactionTotalCell := fmt.Sprintf("%s%d", totalCol, row)
			f.SetCellValue(sheet, transactionTotalCell, transactionTotal)
			setRupiahCurrencyFormat(f, sheet, transactionTotalCell, rupiahStyle)
		}

		// Move row to the next empty row after the products
		row++
	}

	// Set Grand Total row
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Grand Total")

	// Set sum of total_sales and transactions_total
	for i := 0; i < len(dates); i++ {
		col := string('B' + i)
		sumRange := fmt.Sprintf("%s3:%s%d", col, col, row-1) // Ensure correct range
		formulaCell := fmt.Sprintf("%s%d", col, row)
		f.SetCellFormula(sheet, formulaCell, fmt.Sprintf("SUM(%s)", sumRange))
		setRupiahCurrencyFormat(f, sheet, formulaCell, rupiahStyle)
	}

	// Set sum for transactions_total
	startCol := "B"
	endCol := string('B' + len(dates) - 1)
	totalSalesRange := fmt.Sprintf("%s%d:%s%d", startCol, row, endCol, row)

	// Set formula in the total column cell
	formulaCell := fmt.Sprintf("%s%d", totalCol, row) // Total column, same row
	f.SetCellFormula(sheet, formulaCell, fmt.Sprintf("SUM(%s)", totalSalesRange))
	setRupiahCurrencyFormat(f, sheet, formulaCell, rupiahStyle)

	if err := f.Write(buf); err != nil {
		return "", err
	}

	return sheet + ".xlsx", nil
}

func setRupiahCurrencyFormat(f *excelize.File, sheet, cell string, style int) error {
	err := f.SetCellStyle(sheet, cell, cell, style)
	if err != nil {
		fmt.Println("Error applying style:", err)
		return err
	}

	return nil
}

// =====================  V1

// func ReformatSalesReport(sales []structs.Sale, periods []string) map[int]map[string]interface{} {
// 	salesDetail := make(map[int]map[string]interface{})
// 	totalPerDate := make(map[string]float64)
// 	total := 0.0

// 	for _, sale := range sales {
// 		saleDate := sale.CreatedAt.Format("2006-01-02")
// 		for _, detail := range sale.Details {
// 			if detail.Product == nil {
// 				continue // Skip if no product relation
// 			}

// 			categoryId := detail.Product.ProductCategoryID
// 			categoryName := detail.Product.Category.Name
// 			productId := detail.ProductID
// 			productName := detail.Product.Name
// 			totalSales := float64(detail.Price) * float64(detail.TotalItem)

// 			if salesDetail[categoryId] == nil {
// 				salesDetail[categoryId] = map[string]interface{}{
// 					"category_id":    categoryId,
// 					"category_name":  categoryName,
// 					"category_total": 0.0,
// 					"products":       map[string]map[string]interface{}{},
// 				}
// 			}

// 			category := salesDetail[categoryId]
// 			category["category_total"] = category["category_total"].(float64) + totalSales

// 			if category["products"].(map[string]map[string]interface{})[productId.String()] == nil {
// 				category["products"].(map[string]map[string]interface{})[productId.String()] = map[string]interface{}{
// 					"product_id":         productId,
// 					"product_name":       productName,
// 					"transactions":       initializePeriod(periods),
// 					"transactions_total": 0.0,
// 				}
// 			}

// 			product := category["products"].(map[string]map[string]interface{})[productId.String()]
// 			product["transactions_total"] = product["transactions_total"].(float64) + totalSales

// 			transaction := product["transactions"].(map[string]map[string]interface{})
// 			transaction[saleDate]["total_sales"] = transaction[saleDate]["total_sales"].(float64) + totalSales

// 			// Add to total per date and grand total
// 			totalPerDate[saleDate] = totalPerDate[saleDate] + totalSales
// 			total += totalSales
// 		}
// 	}

// 	// You can process totalPerDate or total further if needed
// 	return salesDetail
// }

// func ExportSalesReport(data []map[string]interface{}, periods []string, buf *bytes.Buffer) (string, error) {
// 	f := excelize.NewFile()

// 	// Set up the headers
// 	f.SetCellValue("Sheet1", "A1", "Menu")
// 	f.SetCellValue("Sheet1", "B1", "Period")
// 	f.SetCellValue("Sheet1", "C1", "Total")

// 	// Set up the date headers dynamically based on the periods
// 	for i, date := range periods {
// 		col := string('D' + i)
// 		f.SetCellValue("Sheet1", col+"1", date)
// 	}

// 	row := 2
// 	grandTotalSales := 0.0
// 	grandTotalTransactions := 0.0

// 	// Iterate through categories and products
// 	for _, categoryData := range data {
// 		categoryName := categoryData["category_name"].(string)
// 		categoryTotal := categoryData["category_total"].(float64)

// 		// Write category info
// 		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), categoryName)
// 		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), "Total")
// 		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), categoryTotal)
// 		row++

// 		// Write products under the category
// 		products := categoryData["products"].(map[string]map[string]interface{})
// 		for _, productData := range products {
// 			productName := productData["product_name"].(string)
// 			transactionsTotal := productData["transactions_total"].(float64)

// 			// Write product name and total sales
// 			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), productName)
// 			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), transactionsTotal)

// 			// Correct type assertion for transactions
// 			transactions := productData["transactions"].(map[string]map[string]interface{})
// 			for i, date := range periods {
// 				col := string('D' + i)
// 				if transaction, ok := transactions[date]; ok {
// 					f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", col, row), transaction["total_sales"])
// 				} else {
// 					f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", col, row), 0)
// 				}
// 			}

// 			grandTotalSales += transactionsTotal
// 			grandTotalTransactions += transactionsTotal
// 			row++
// 		}
// 	}

// 	// Add Grand Total at the end
// 	f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), "Grand Total")
// 	f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), grandTotalSales)
// 	f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), grandTotalTransactions) // Optionally show grand total transactions

// 	// Save the file
// filename := fmt.Sprintf("%s_%s", time.Now().Format("20061021545"), "SalesReport.xlsx")
// if err := f.SaveAs(filename); err != nil {
// 	return "", err
// }

// err := f.Write(buf)
// if err != nil {
// 	return "", err
// }

// 	return filename, nil
// }
