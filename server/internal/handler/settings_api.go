package handler

import (
	"net/http"
	"os"

	"easyaccounts-server/internal/model"
	"easyaccounts-server/pkg/response"
	"github.com/gin-gonic/gin"
)

// --- Account CRUD ---

func (api *APIController) AddAccount(c *gin.Context) {
	userID := getUserID(c)

	var account model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	account.UserID = userID

	if err := api.DB.Create(&account).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to create account")
		return
	}
	response.Success(c, account)
}

func (api *APIController) UpdateAccount(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	var account model.Account
	if err := api.DB.Where("user_id = ?", userID).First(&account, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, 404, "Account not found")
		return
	}
	if err := c.ShouldBindJSON(&account); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	account.UserID = userID

	api.DB.Save(&account)
	response.Success(c, account)
}

func (api *APIController) DeleteAccount(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	if err := api.DB.Where("user_id = ?", userID).Delete(&model.Account{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to delete account")
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

// --- Category CRUD ---

func (api *APIController) AddCategory(c *gin.Context) {
	userID := getUserID(c)

	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	category.UserID = userID

	if err := api.DB.Create(&category).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to create category")
		return
	}
	response.Success(c, category)
}

func (api *APIController) UpdateCategory(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	var category model.Category
	if err := api.DB.Where("user_id = ?", userID).First(&category, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, 404, "Category not found")
		return
	}
	if err := c.ShouldBindJSON(&category); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	category.UserID = userID

	api.DB.Save(&category)
	response.Success(c, category)
}

func (api *APIController) DeleteCategory(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")

	// Check if this category has children
	var childCount int64
	api.DB.Where("user_id = ?", userID).Model(&model.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		response.Error(c, http.StatusBadRequest, 400, "该分类下还有子分类，请先删除子分类")
		return
	}

	if err := api.DB.Where("user_id = ?", userID).Delete(&model.Category{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to delete category")
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}

// --- Template CRUD ---

func (api *APIController) GetTemplates(c *gin.Context) {
	userID := getUserID(c)

	var templates []model.Template
	api.DB.Where("user_id = ?", userID).Preload("Account").Preload("Category").Find(&templates)
	response.Success(c, templates)
}

func (api *APIController) AddTemplate(c *gin.Context) {
	userID := getUserID(c)

	var template model.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	template.UserID = userID

	if err := api.DB.Create(&template).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to create template")
		return
	}
	api.DB.Where("user_id = ?", userID).Preload("Account").Preload("Category").First(&template, template.ID)
	response.Success(c, template)
}

func (api *APIController) UpdateTemplate(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	var template model.Template
	if err := api.DB.Where("user_id = ?", userID).First(&template, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, 404, "Template not found")
		return
	}
	if err := c.ShouldBindJSON(&template); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	template.UserID = userID

	api.DB.Save(&template)
	api.DB.Where("user_id = ?", userID).Preload("Account").Preload("Category").First(&template, template.ID)
	response.Success(c, template)
}

func (api *APIController) DeleteTemplate(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	api.DB.Where("user_id = ?", userID).Delete(&model.Template{}, id)
	response.Success(c, gin.H{"status": "deleted"})
}

// --- ScheduledTransaction CRUD ---

func (api *APIController) GetSchedules(c *gin.Context) {
	userID := getUserID(c)

	var schedules []model.ScheduledTransaction
	api.DB.Where("user_id = ?", userID).Preload("Account").Preload("Category").Find(&schedules)
	response.Success(c, schedules)
}

func (api *APIController) AddSchedule(c *gin.Context) {
	userID := getUserID(c)

	var schedule model.ScheduledTransaction
	if err := c.ShouldBindJSON(&schedule); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	schedule.UserID = userID

	if err := api.DB.Create(&schedule).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to create schedule")
		return
	}
	api.DB.Where("user_id = ?", userID).Preload("Account").Preload("Category").First(&schedule, schedule.ID)
	response.Success(c, schedule)
}

func (api *APIController) UpdateSchedule(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	var schedule model.ScheduledTransaction
	if err := api.DB.Where("user_id = ?", userID).First(&schedule, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, 404, "Schedule not found")
		return
	}
	if err := c.ShouldBindJSON(&schedule); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	schedule.UserID = userID

	api.DB.Save(&schedule)
	api.DB.Where("user_id = ?", userID).Preload("Account").Preload("Category").First(&schedule, schedule.ID)
	response.Success(c, schedule)
}

func (api *APIController) DeleteSchedule(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	api.DB.Where("user_id = ?", userID).Delete(&model.ScheduledTransaction{}, id)
	response.Success(c, gin.H{"status": "deleted"})
}

// --- SystemSettings ---

func (api *APIController) GetSettings(c *gin.Context) {
	userID := getUserID(c)

	var settings []model.SystemSetting
	api.DB.Where("user_id = ?", userID).Find(&settings)
	// Convert to map
	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}
	response.Success(c, settingsMap)
}

func (api *APIController) UpdateSettings(c *gin.Context) {
	userID := getUserID(c)

	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}

	for k, v := range body {
		var setting model.SystemSetting
		result := api.DB.Where("user_id = ?", userID).Where("key = ?", k).First(&setting)
		if result.Error == nil {
			setting.Value = v
			api.DB.Save(&setting)
		} else {
			api.DB.Create(&model.SystemSetting{Key: k, Value: v, UserID: userID})
		}
	}
	response.Success(c, gin.H{"status": "updated"})
}



// --- DB Backup ---

func (api *APIController) DownloadBackup(c *gin.Context) {
	_ = getUserID(c)

	dbPath := api.DBPath
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		response.Error(c, http.StatusNotFound, 404, "Database file not found")
		return
	}
	c.FileAttachment(dbPath, "easyaccounts_backup.db")
}

// --- TransactionType CRUD ---

func (api *APIController) GetTransactionTypes(c *gin.Context) {
	userID := getUserID(c)

	var types []model.TransactionType
	api.DB.Where("user_id = ?", userID).Find(&types)
	response.Success(c, types)
}

func (api *APIController) AddTransactionType(c *gin.Context) {
	userID := getUserID(c)

	var tType model.TransactionType
	if err := c.ShouldBindJSON(&tType); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tType.UserID = userID

	if err := api.DB.Create(&tType).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to create transaction type")
		return
	}
	response.Success(c, tType)
}

func (api *APIController) UpdateTransactionType(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	var tType model.TransactionType
	if err := api.DB.Where("user_id = ?", userID).First(&tType, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, 404, "Transaction type not found")
		return
	}
	if err := c.ShouldBindJSON(&tType); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tType.UserID = userID

	api.DB.Save(&tType)
	response.Success(c, tType)
}

func (api *APIController) DeleteTransactionType(c *gin.Context) {
	userID := getUserID(c)

	id := c.Param("id")
	if err := api.DB.Where("user_id = ?", userID).Delete(&model.TransactionType{}, id).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to delete transaction type")
		return
	}
	response.Success(c, gin.H{"status": "deleted"})
}
