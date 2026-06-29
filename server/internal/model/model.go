package model

import "time"

type Account struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`           // e.g., bank, cash, alipay, wechat (Used as icon/subtype)
	AccountType   string    `json:"account_type"`   // asset, liability
	Balance       float64   `json:"balance"`
	ExcludeAmount float64   `json:"exclude_amount"` // 不计入金额
	AccountNumber string    `json:"account_number"` // 卡号/账号
	Remark        string    `json:"remark"`         // 备注
	CreatedAt     time.Time `json:"created_at"`
}

type Category struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UserID    uint   `gorm:"index" json:"user_id"`
	Name      string `json:"name"`
	Type      string `json:"type"` // Bound TransactionType Name (e.g., 借入, 收入)
	ParentID  uint   `json:"parent_id"`
	Icon      string `json:"icon"`
	IsExclude bool   `json:"is_exclude"`
}

type Transaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	AccountID  uint      `json:"account_id"`
	CategoryID uint      `json:"category_id"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"` // income, expense, transfer
	Remark     string    `json:"remark"`
	RecordDate time.Time `json:"record_date"` // effective date of the transaction
	CreatedAt  time.Time `json:"created_at"`

	Account  Account  `gorm:"foreignKey:AccountID" json:"account"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
}

type Template struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Name       string    `json:"name"` // 模板名称
	Type       string    `json:"type"` // income, expense, transfer
	Amount     float64   `json:"amount"`
	CategoryID uint      `json:"category_id"`
	AccountID  uint      `json:"account_id"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`

	Account  Account  `gorm:"foreignKey:AccountID" json:"account"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
}

type ScheduledTransaction struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Name       string    `json:"name"`        // 定时任务名称
	CronRule   string    `json:"cron_rule"`   // 周期规则，例如 "每月1号"
	Type       string    `json:"type"`        // income, expense, transfer
	Amount     float64   `json:"amount"`
	CategoryID uint      `json:"category_id"`
	AccountID  uint      `json:"account_id"`
	Remark     string    `json:"remark"`
	IsActive   bool      `json:"is_active"`   // 是否启用
	CreatedAt  time.Time `json:"created_at"`

	Account  Account  `gorm:"foreignKey:AccountID" json:"account"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
}

type SystemSetting struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	Key       string    `gorm:"primaryKey" json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransactionType struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Name       string    `json:"name"`        // e.g. 借入, 收入
	Flow       string    `json:"flow"`        // inflow, outflow, transfer
	IsExclude  bool      `json:"is_exclude"`  // 不计入总金额
	IsDisabled bool      `json:"is_disabled"` // 禁用
	CreatedAt  time.Time `json:"created_at"`
}

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"` // Don't return password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
}
