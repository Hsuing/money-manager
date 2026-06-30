package repo

import (
	"easyaccounts-server/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// InitSchema automatically migrates the schema, creating tables if they don't exist
func InitSchema(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Account{},
		&model.Category{},
		&model.Transaction{},
		&model.Template{},
		&model.ScheduledTransaction{},
		&model.SystemSetting{},
		&model.TransactionType{},
		&model.User{},
	)
}

// SeedData seeds some initial data if the database is empty
func SeedData(db *gorm.DB) {
	var count int64
	
	// Seed Admin User
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("hx13231228792"), bcrypt.DefaultCost)
		db.Create(&model.User{
			Username:     "hanusers",
			PasswordHash: string(hash),
		})
	}

	// Seed Accounts
	db.Model(&model.Account{}).Count(&count)
	if count == 0 {
		db.Create(&model.Account{Name: "Alipay", Type: "alipay", AccountType: "asset", Balance: 2000})
		db.Create(&model.Account{Name: "Bank Card", Type: "bank", AccountType: "asset", Balance: 15000})
	} else {
		// Legacy upgrade: set missing account_type to 'asset'
		db.Model(&model.Account{}).Where("account_type = '' OR account_type IS NULL").Update("account_type", "asset")
	}

	db.Model(&model.Category{}).Count(&count)
	if count == 0 {
		db.Create(&model.Category{Name: "Food & Dining", Type: "expense", Icon: "dining"})
		db.Create(&model.Category{Name: "Transportation", Type: "expense", Icon: "transport"})
		db.Create(&model.Category{Name: "Salary", Type: "income", Icon: "salary"})
		db.Create(&model.Category{Name: "Investment", Type: "income", Icon: "investment"})
	}

	var typeCount int64
	db.Model(&model.TransactionType{}).Count(&typeCount)
	if typeCount == 0 {
		types := []model.TransactionType{
			{Name: "收入", Flow: "inflow", IsExclude: false, IsDisabled: false},
			{Name: "支出", Flow: "outflow", IsExclude: false, IsDisabled: false},
			{Name: "内部转账", Flow: "transfer", IsExclude: false, IsDisabled: false},
			{Name: "借入", Flow: "inflow", IsExclude: true, IsDisabled: false},
			{Name: "还钱", Flow: "outflow", IsExclude: true, IsDisabled: false},
			{Name: "借出", Flow: "outflow", IsExclude: true, IsDisabled: false},
			{Name: "收钱", Flow: "inflow", IsExclude: true, IsDisabled: false},
		}
		for _, t := range types {
			db.Create(&t)
		}
	}

	// Migrate legacy data: if any user exists, assign orphan records to the first user.
	var firstUser model.User
	if err := db.Order("id asc").First(&firstUser).Error; err == nil && firstUser.ID > 0 {
		db.Model(&model.Account{}).Where("user_id = ? OR user_id IS NULL", 0).Update("user_id", firstUser.ID)
		db.Model(&model.Category{}).Where("user_id = ? OR user_id IS NULL", 0).Update("user_id", firstUser.ID)
		db.Model(&model.Transaction{}).Where("user_id = ? OR user_id IS NULL", 0).Update("user_id", firstUser.ID)
		db.Model(&model.Template{}).Where("user_id = ? OR user_id IS NULL", 0).Update("user_id", firstUser.ID)
		db.Model(&model.ScheduledTransaction{}).Where("user_id = ? OR user_id IS NULL", 0).Update("user_id", firstUser.ID)
		db.Model(&model.TransactionType{}).Where("user_id = ? OR user_id IS NULL", 0).Update("user_id", firstUser.ID)
		db.Model(&model.SystemSetting{}).Where("user_id = ? OR user_id IS NULL", 0).Update("user_id", firstUser.ID)
	}
}
