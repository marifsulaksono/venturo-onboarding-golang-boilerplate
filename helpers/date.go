package helpers

import (
	"simple-crud-rnd/structs"
	"time"
)

// Helper function to generate a list of dates between startDate and endDate
func GenerateDateRange(startDate, endDate string) []string {
	var dates []string
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	for date := start; !date.After(end); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date.Format("2006-01-02"))
	}
	return dates
}

// Helper function to fill missing dates with total_sales = 0
func FillMissingDates(dates []string, transactions []structs.SalesSummaryNew) []structs.SalesSummaryNew {
	dateMap := make(map[string]structs.SalesSummaryNew)
	for _, t := range transactions {
		dateMap[t.DateSales.Format("2006-01-02")] = t
	}

	var result []structs.SalesSummaryNew
	for _, date := range dates {
		if transaction, exists := dateMap[date]; exists {
			result = append(result, transaction)
		} else {
			result = append(result, structs.SalesSummaryNew{
				DateSales:  ParseDate(date),
				TotalSales: 0,
			})
		}
	}
	return result
}

// Helper function to parse a date string
func ParseDate(dateStr string) time.Time {
	date, _ := time.Parse("2006-01-02", dateStr)
	return date
}
