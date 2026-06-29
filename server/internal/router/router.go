package router

import (
	"easyaccounts-server/internal/handler"
	"easyaccounts-server/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, dbPath string) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	apiCtrl := handler.NewAPIController(db, dbPath)
	
	api := r.Group("/api")
	
	// Public routes
	api.POST("/login", apiCtrl.Login)
	api.POST("/register", apiCtrl.Register)
	api.GET("/captcha", apiCtrl.GetCaptcha)

	// Protected routes
	api.Use(middleware.JWTAuth())
	{
		api.GET("/dashboard", apiCtrl.GetDashboard)
		
		api.GET("/transactions", apiCtrl.GetTransactions)
		api.POST("/transactions", apiCtrl.AddTransaction)
		api.DELETE("/transactions/:id", apiCtrl.DeleteTransaction)
		
		api.GET("/categories", apiCtrl.GetCategories)
		api.POST("/categories", apiCtrl.AddCategory)
		api.PUT("/categories/:id", apiCtrl.UpdateCategory)
		api.DELETE("/categories/:id", apiCtrl.DeleteCategory)
		
		api.GET("/accounts", apiCtrl.GetAccounts)
		api.POST("/accounts", apiCtrl.AddAccount)
		api.PUT("/accounts/:id", apiCtrl.UpdateAccount)
		api.DELETE("/accounts/:id", apiCtrl.DeleteAccount)

		api.GET("/templates", apiCtrl.GetTemplates)
		api.POST("/templates", apiCtrl.AddTemplate)
		api.PUT("/templates/:id", apiCtrl.UpdateTemplate)
		api.DELETE("/templates/:id", apiCtrl.DeleteTemplate)

		api.GET("/schedules", apiCtrl.GetSchedules)
		api.POST("/schedules", apiCtrl.AddSchedule)
		api.PUT("/schedules/:id", apiCtrl.UpdateSchedule)
		api.DELETE("/schedules/:id", apiCtrl.DeleteSchedule)

		api.GET("/transaction_types", apiCtrl.GetTransactionTypes)
		api.POST("/transaction_types", apiCtrl.AddTransactionType)
		api.PUT("/transaction_types/:id", apiCtrl.UpdateTransactionType)
		api.DELETE("/transaction_types/:id", apiCtrl.DeleteTransactionType)

		api.GET("/settings", apiCtrl.GetSettings)
		api.POST("/settings", apiCtrl.UpdateSettings)
		api.GET("/backup", apiCtrl.DownloadBackup)
	}

	return r
}
