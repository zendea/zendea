package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"zendea/form"
	"zendea/model"
	"zendea/util"
)

// BaseController controller
type BaseController struct {
}

// BindAndValidate bind and validate
func (c *BaseController) BindAndValidate(ctx *gin.Context, obj interface{}) bool {
	if err := form.Bind(ctx, obj); err != nil {
		c.Fail(ctx, &util.CodeError{Code: -1, Message: err.Error()})
		return false
	}
	return true
}

// GetCurrentUser get current user from contexg
func (c *BaseController) GetCurrentUser(ctx *gin.Context) *model.User {
	if currentUser, ok := ctx.Get("CurrentUser"); ok {
		return currentUser.(*model.User)
	}
	return nil
}

// Success output json data
func (c *BaseController) Success(ctx *gin.Context, data interface{}) {
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"success": true,
		"data":    data,
	})
}

// Fail output error
func (c *BaseController) Fail(ctx *gin.Context, error *util.CodeError) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":    error.Code,
		"message": error.Message,
	})
	return
}
