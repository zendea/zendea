package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"zendea/service"
	"zendea/util"
)

// SigninRequired signin required
func SigninRequired(ctx *gin.Context) {
	user := service.UserService.GetCurrent(ctx)
	if user == nil {
		err := util.ErrorNotLogin
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    err.Code,
			"message": err.Message,
		})

	}
	ctx.Next()
}