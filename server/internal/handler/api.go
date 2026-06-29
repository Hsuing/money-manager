package handler

import (
	"net/http"
	"time"
	"strconv"

	"easyaccounts-server/internal/model"
	"easyaccounts-server/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIController struct {
	DB     *gorm.DB
	DBPath string
}

func NewAPIController(db *gorm.DB, dbPath string) *APIController {
	return &APIController{DB: db, DBPath: dbPath}
}

func getUserID(c *gin.Context) uint {
	v, _ := c.Get("user_id")
	if f, ok := v.(float64); ok {
		return uint(f)
	}
	return 0
}


func (api *APIController) GetDashboard(c *gin.Context) {
	userID := getUserID(c)

	var totalAssets, totalLiabilities, totalExcludeAssets, totalExcludeLiabilities float64
	
	api.DB.Where("user_id = ?", userID).Model(&model.Account{}).Where("account_type = ?", "asset").Select("COALESCE(SUM(balance), 0)").Scan(&totalAssets)
	api.DB.Where("user_id = ?", userID).Model(&model.Account{}).Where("account_type = ?", "liability").Select("COALESCE(SUM(balance), 0)").Scan(&totalLiabilities)
	
	api.DB.Where("user_id = ?", userID).Model(&model.Account{}).Where("account_type = ?", "asset").Select("COALESCE(SUM(exclude_amount), 0)").Scan(&totalExcludeAssets)
	api.DB.Where("user_id = ?", userID).Model(&model.Account{}).Where("account_type = ?", "liability").Select("COALESCE(SUM(exclude_amount), 0)").Scan(&totalExcludeLiabilities)

	netAssets := (totalAssets - totalExcludeAssets) - (totalLiabilities - totalExcludeLiabilities)

	// 支持通过 year 参数查询指定年份，默认为当前年
	targetYear := time.Now().Year()
	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil && y > 2000 && y <= time.Now().Year()+10 {
			targetYear = y
		}
	}

	startOfYear := time.Date(targetYear, 1, 1, 0, 0, 0, 0, time.Local)
	endOfYear := startOfYear.AddDate(1, 0, 0)

	var income, expense float64
	api.DB.Where("transactions.user_id = ?", userID).Model(&model.Transaction{}).
		Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
		Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "inflow", "income", "收入", false, startOfYear, endOfYear).
		Select("COALESCE(SUM(amount), 0)").Scan(&income)

	api.DB.Where("transactions.user_id = ?", userID).Model(&model.Transaction{}).
		Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
		Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "outflow", "expense", "支出", false, startOfYear, endOfYear).
		Select("COALESCE(SUM(amount), 0)").Scan(&expense)

	// 当前月统计（历史年份用12月，当前年用当前月）
	currentMonth := time.Now().Month()
	if targetYear < time.Now().Year() {
		currentMonth = time.December
	}
	currentMonthStart := time.Date(targetYear, currentMonth, 1, 0, 0, 0, 0, time.Local)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, 0)

	var monthIncome, monthExpense float64
	api.DB.Where("transactions.user_id = ?", userID).Model(&model.Transaction{}).
		Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
		Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "inflow", "income", "收入", false, currentMonthStart, currentMonthEnd).
		Select("COALESCE(SUM(amount), 0)").Scan(&monthIncome)

	api.DB.Where("transactions.user_id = ?", userID).Model(&model.Transaction{}).
		Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
		Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "outflow", "expense", "支出", false, currentMonthStart, currentMonthEnd).
		Select("COALESCE(SUM(amount), 0)").Scan(&monthExpense)

	// 按月循环：历史年份循环12个月，当前年循环到当前月
	loopMonths := int(time.Now().Month())
	if targetYear < time.Now().Year() {
		loopMonths = 12
	}
	var monthlyStats []map[string]interface{}
	for m := 1; m <= loopMonths; m++ {
		mStart := time.Date(targetYear, time.Month(m), 1, 0, 0, 0, 0, time.Local)
		mEnd := mStart.AddDate(0, 1, 0)
		var mInc, mExp float64
		api.DB.Where("transactions.user_id = ?", userID).Model(&model.Transaction{}).
			Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
			Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "inflow", "income", "收入", false, mStart, mEnd).
			Select("COALESCE(SUM(amount), 0)").Scan(&mInc)
		api.DB.Where("transactions.user_id = ?", userID).Model(&model.Transaction{}).
			Joins("LEFT JOIN transaction_types ON transactions.type = transaction_types.name").
			Where("(transaction_types.flow = ? OR transactions.type = ? OR transactions.type = ?) AND (transaction_types.is_exclude IS NULL OR transaction_types.is_exclude = ?) AND record_date >= ? AND record_date < ?", "outflow", "expense", "支出", false, mStart, mEnd).
			Select("COALESCE(SUM(amount), 0)").Scan(&mExp)
		monthlyStats = append(monthlyStats, map[string]interface{}{
			"month":   m,
			"income":  mInc,
			"expense": mExp,
			"balance": mInc - mExp,
		})
	}

	c.JSON(http.StatusOK, response.Response{
		Code: 0,
		Msg:  "success",
		Data: gin.H{
			"total_assets":  totalAssets,
			"net_assets":    netAssets,
			"year":          targetYear,
			"year_income":   income,
			"year_expense":  expense,
			"year_balance":  income - expense,
			"month":         int(currentMonth),
			"month_income":  monthIncome,
			"month_expense": monthExpense,
			"month_balance": monthIncome - monthExpense,
			"monthly_stats": monthlyStats,
		},
	})
}

func (api *APIController) GetTransactions(c *gin.Context) {
	userID := getUserID(c)

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	query := api.DB.Where("user_id = ?", userID).Preload("Account").Preload("Category")

	if yearStr != "" && monthStr != "" {
		year, _ := strconv.Atoi(yearStr)
		month, _ := strconv.Atoi(monthStr)
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		endDate := startDate.AddDate(0, 1, 0)
		query = query.Where("record_date >= ? AND record_date < ?", startDate, endDate)
	}

	var transactions []model.Transaction
	query.Order("record_date DESC").Find(&transactions)
	response.Success(c, transactions)
}

func (api *APIController) AddTransaction(c *gin.Context) {
	userID := getUserID(c)

	var tx model.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx.UserID = userID

	
	if tx.RecordDate.IsZero() {
		tx.RecordDate = time.Now()
	}

	dbTx := api.DB.Begin()

	if err := dbTx.Create(&tx).Error; err != nil {
		dbTx.Rollback()
		response.Error(c, http.StatusInternalServerError, 500, "Failed to create transaction")
		return
	}

	var account model.Account
	if err := dbTx.Where("user_id = ?", userID).First(&account, tx.AccountID).Error; err != nil {
		dbTx.Rollback()
		response.Error(c, http.StatusInternalServerError, 500, "Account not found")
		return
	}

	// Retrieve the flow type (handle legacy fallback)
	flow := ""
	if tx.Type == "income" || tx.Type == "expense" {
		flow = tx.Type
		if flow == "income" { flow = "inflow" }
		if flow == "expense" { flow = "outflow" }
	} else {
		var tType model.TransactionType
		if err := dbTx.Where("user_id = ?", userID).Where("name = ?", tx.Type).First(&tType).Error; err == nil {
			flow = tType.Flow
		}
	}

	if flow == "outflow" {
		account.Balance -= tx.Amount
	} else if flow == "inflow" {
		account.Balance += tx.Amount
	}

	if err := dbTx.Where("user_id = ?", userID).Save(&account).Error; err != nil {
		dbTx.Rollback()
		response.Error(c, http.StatusInternalServerError, 500, "Failed to update account balance")
		return
	}

	dbTx.Commit()

	api.DB.Where("user_id = ?", userID).Preload("Account").Preload("Category").First(&tx, tx.ID)
	response.Success(c, tx)
}

func (api *APIController) DeleteTransaction(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	dbTx := api.DB.Begin()

	var tx model.Transaction
	if err := dbTx.Where("user_id = ?", userID).First(&tx, id).Error; err != nil {
		dbTx.Rollback()
		response.Error(c, http.StatusNotFound, 404, "Transaction not found")
		return
	}

	var account model.Account
	if err := dbTx.Where("user_id = ?", userID).First(&account, tx.AccountID).Error; err == nil {
		// Reverse balance
		flow := ""
		if tx.Type == "income" || tx.Type == "expense" {
			flow = tx.Type
			if flow == "income" { flow = "inflow" }
			if flow == "expense" { flow = "outflow" }
		} else {
			var tType model.TransactionType
			if err := dbTx.Where("user_id = ?", userID).Where("name = ?", tx.Type).First(&tType).Error; err == nil {
				flow = tType.Flow
			}
		}

		if flow == "outflow" {
			account.Balance += tx.Amount
		} else if flow == "inflow" {
			account.Balance -= tx.Amount
		}
		dbTx.Where("user_id = ?", userID).Save(&account)
	}

	if err := dbTx.Where("user_id = ?", userID).Delete(&model.Transaction{}, id).Error; err != nil {
		dbTx.Rollback()
		response.Error(c, http.StatusInternalServerError, 500, "Failed to delete transaction")
		return
	}

	dbTx.Commit()
	response.Success(c, gin.H{"status": "deleted"})
}

func (api *APIController) GetCategories(c *gin.Context) {
	userID := getUserID(c)

	var categories []model.Category
	api.DB.Where("user_id = ?", userID).Find(&categories)
	response.Success(c, categories)
}

func (api *APIController) GetAccounts(c *gin.Context) {
	userID := getUserID(c)

	var accounts []model.Account
	api.DB.Where("user_id = ?", userID).Find(&accounts)
	response.Success(c, accounts)
}
