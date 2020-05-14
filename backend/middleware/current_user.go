package middleware

import (
	"github.com/gin-gonic/gin"

	"zendea/service"
)

func CurrentUser(ctx *gin.Context) {
	ctx.Set("CurrentUser", service.UserService.GetCurrent(ctx))
	ctx.Next()
}
