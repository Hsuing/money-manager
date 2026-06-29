package handler

import (
	"image/color"
	"net/http"

	"easyaccounts-server/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// store must be exported or accessible to auth.go
var CaptchaStore = base64Captcha.DefaultMemStore

func (api *APIController) GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverMath(40, 120, 2, base64Captcha.OptionShowHollowLine, &color.RGBA{R: 240, G: 240, B: 246, A: 255}, nil, nil)
	captcha := base64Captcha.NewCaptcha(driver, CaptchaStore)
	
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "Failed to generate captcha")
		return
	}

	response.Success(c, gin.H{
		"captchaId":  id,
		"captchaImg": b64s,
	})
}
