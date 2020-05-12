package admin

import (
	"strconv"
	"github.com/gin-gonic/gin"

	"zendea/controller"
	"zendea/util"
	"zendea/util/sqlcnd"
	"zendea/form"
	"zendea/model"
	"zendea/cache"
	"zendea/service"
)

// UserController user controller
type UserController struct {
	controller.BaseController
}

// Show show user
func (c *UserController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		user := service.UserService.Get(gDto.ID)
		if user == nil {
			c.Fail(ctx, util.NewErrorMsg("User not found, id=" + strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, c.buildUserItem(user))
	}
}

// Store 创建用户
func (c *UserController) Store(ctx *gin.Context) {
	c.Success(ctx, nil)
}

// Update 更新用户信息
func (c *UserController) Update(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	user := service.UserService.Get(gDto.ID)
	if user == nil {
		c.Fail(ctx, util.NewErrorMsg("User not found, id=" + strconv.FormatInt(gDto.ID, 10)))
		return
	}
	
	var userForm form.UserUpdateForm
	if !c.BindAndValidate(ctx, &userForm) {
		return
	}
	userForm.ID = gDto.ID
	err := service.UserService.Update(userForm)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, user)
}

// Delete delete user
func (c *UserController) Delete(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if !c.BindAndValidate(ctx, &gDto) {
		return
	}
	service.UserService.Delete(gDto.ID)
	c.Success(ctx, nil)
}

// List list users
func (c *UserController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	id := ctx.Request.FormValue("id")
	nickname := ctx.Request.FormValue("nickname")
	username := ctx.Request.FormValue("username")

	conditions := sqlcnd.NewSqlCnd()
	if len(id) > 0 {
		conditions.Eq("id", id)
	}
	if len(username) > 0 {
		conditions.Eq("username", username)
	}
	if len(nickname) > 0 {
		conditions.Like("nickname", nickname)
	}
	list, paging := service.UserService.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, user := range list {
		results = append(results, c.buildUserItem(&user))
	}

	c.Success(ctx, &sqlcnd.PageResult{Results: results, Page: paging})
}

func (c *UserController) buildUserItem(user *model.User) map[string]interface{} {
	score := cache.UserCache.GetScore(user.ID)

	result := make(map[string]interface{})
	result["id"] = user.ID
	result["status"] = user.Status
	result["level"] = user.Level
	result["username"] = user.Username.String
	result["nickname"] = user.Nickname
	result["avatar"] = user.Avatar
	result["email"] =  user.Email.String
	result["score"] =  score
	result["createTime"] = user.CreateTime
	result["updateTime"] = user.UpdateTime

	return result
}