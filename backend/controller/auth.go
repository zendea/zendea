package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/dchest/captcha"

	"zendea/builder"
	"zendea/service"
	"zendea/util"
	"zendea/form"
)

// AuthController auth controller
type AuthController struct {
	BaseController
}

// Signup create
func (c *AuthController) Signup(ctx *gin.Context) {
	ref := ctx.Request.FormValue("ref")
	var dto form.AuthSignupForm
	if c.BindAndValidate(ctx, &dto) {
		if !captcha.VerifyString(dto.CaptchaID, dto.CaptchaCode) {
			c.Fail(ctx, util.ErrorCaptchaWrong)
			return
		}

		user, err := service.UserService.Create(dto.Username, dto.Email, dto.Nickname, dto.Password, dto.RePassword)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, gin.H{
			"user":  builder.BuildUser(user),
			"ref":   ref,
		})
	}
}
