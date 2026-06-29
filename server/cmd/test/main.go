package main

import (
	"fmt"
	"time"

	"easyaccounts-server/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data/easyaccounts.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	currentYear := time.Now().Year()
	startOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.Local)
	endOfYear := startOfYear.AddDate(1, 0, 0)

	var expense float64
	db.Model(&model.Transaction{}).
		Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
		Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "outflow", "expense", "支出", false, startOfYear, endOfYear).
		Select("COALESCE(SUM(amount), 0)").Scan(&expense)

	fmt.Printf("Year Expense: %v\n", expense)

	currentMonthStart := time.Date(currentYear, time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, 0)

	var monthExpense float64
	db.Model(&model.Transaction{}).
		Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
		Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "outflow", "expense", "支出", false, currentMonthStart, currentMonthEnd).
		Select("COALESCE(SUM(amount), 0)").Scan(&monthExpense)

	fmt.Printf("Month Expense: %v\n", monthExpense)

	for m := 1; m <= int(time.Now().Month()); m++ {
		mStart := time.Date(currentYear, time.Month(m), 1, 0, 0, 0, 0, time.Local)
		mEnd := mStart.AddDate(0, 1, 0)
		var mExp float64
		err = db.Model(&model.Transaction{}).
			Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
			Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "outflow", "expense", "支出", false, mStart, mEnd).
			Select("COALESCE(SUM(amount), 0)").Scan(&mExp).Error
		fmt.Printf("Month %d: Expense: %v, Error: %v\n", m, mExp, err)
	}
}
