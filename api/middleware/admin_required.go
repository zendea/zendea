package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"zendea/model"
	"zendea/service"
	"zendea/util"
)

// AdminRequired admin required
func AdminRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := service.UserService.GetCurrent(ctx)
		if user == nil {
			err := util.ErrorNotLogin
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    err.Code,
				"message": err.Message,
			})
			return
		}
		if user.Level != model.UserLevelAdmin {
			err := util.ErrorPermissionDenied
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code":    err.Code,
				"message": err.Message,
			})
			return
		}
		ctx.Next()
	}
}
