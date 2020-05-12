package admin

import (
	"strconv"
	"github.com/gin-gonic/gin"

	"zendea/builder"
	"zendea/controller"
	"zendea/util"
	"zendea/util/sqlcnd"
	"zendea/form"
	"zendea/service"
)

// UserScoreController user score controller
type UserScoreController struct {
	controller.BaseController
}

// Show 显示积分
func (c *UserScoreController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		userScore := service.UserService.Get(gDto.ID)
		if userScore == nil {
			c.Fail(ctx, util.NewErrorMsg("User score not found, id=" + strconv.FormatInt(gDto.ID, 10)))
			return
		}
		c.Success(ctx, userScore)
	}
}

// List 显示积分列表
func (c *UserScoreController) List(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	limit := form.FormValueIntDefault(ctx, "limit", 20)
	userId := ctx.Request.FormValue("userId")

	conditions := sqlcnd.NewSqlCnd()
	if len(userId) > 0 {
		conditions.Eq("user_id", userId)
	}

	list, paging := service.UserScoreService.List(conditions.Page(page, limit).Desc("id"))

	var results []map[string]interface{}
	for _, userScore := range list {
		item := util.StructToMap(userScore)
		item["user"] = builder.BuildUserDefaultIfNull(userScore.UserId)
		results = append(results, item)
	}

	c.Success(ctx, &sqlcnd.PageResult{Results: results, Page: paging})
}