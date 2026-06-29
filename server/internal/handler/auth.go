package handler

import (
	"net/http"

	"easyaccounts-server/internal/model"
	"easyaccounts-server/pkg/jwt"
	"easyaccounts-server/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

func (ctrl *APIController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "Invalid request body")
		return
	}

	var user model.User
	if err := ctrl.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.Error(c, http.StatusUnauthorized, 401, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Error(c, http.StatusUnauthorized, 401, "Invalid username or password")
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to generate token")
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=3,alphanum"`
	Password   string `json:"password" binding:"required,min=6"`
	CaptchaId  string `json:"captchaId" binding:"required"`
	CaptchaVal string `json:"captchaVal" binding:"required"`
}

func (ctrl *APIController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "参数不合法（需要用户名、密码及图形验证码）")
		return
	}

	if !CaptchaStore.Verify(req.CaptchaId, req.CaptchaVal, true) {
		response.Error(c, http.StatusBadRequest, 400, "图形验证码错误或已过期")
		return
	}

	var existing model.User
	if err := ctrl.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		response.Error(c, http.StatusConflict, 409, "用户名已被注册")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "密码加密失败")
		return
	}

	user := model.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	if err := ctrl.DB.Create(&user).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "创建用户失败")
		return
	}

	response.Success(c, gin.H{
		"message": "注册成功",
	})
}
