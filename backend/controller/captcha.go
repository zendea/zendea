package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/dchest/captcha"

	"zendea/util/urls"
)

// CaptchaController captcha controller
type CaptchaController struct {
	BaseController
}

// GetRequest request captcha id and url
func (c *CaptchaController) GetRequest(ctx *gin.Context) {

	captchaID := captcha.NewLen(4)
	captchaURL := urls.AbsUrl("/api/captcha/show/" + captchaID)

	data := make(map[string]interface{})
	data["captchaId"] = captchaID
	data["captchaUrl"] = captchaURL
	
	c.Success(ctx, data)
}

// Show show captcha image
func (c *CaptchaController) Show(ctx *gin.Context) {
	captchaID := ctx.Param("captchaId")
	if captchaID == "" {
		return
	}
	if !captcha.Reload(captchaID) {
		return
	}

	ctx.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Writer.Header().Set("Pragma", "no-cache")
	ctx.Writer.Header().Set("Expires", "10")
	ctx.Writer.Header().Set("Content-Type", "image/png")
	captcha.WriteImage(ctx.Writer, captchaID, captcha.StdWidth, captcha.StdHeight)
}
