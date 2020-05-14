package controller

import (
	"github.com/gin-gonic/gin"
	"strings"

	"zendea/builder"
	"zendea/cache"
	"zendea/form"
	"zendea/model"
	"zendea/service"
	"zendea/util"
	"zendea/util/sqlcnd"
)

type UserController struct {
	BaseController
}

// GetCurrent get current user
func (c *UserController) GetCurrent(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	ctx.IndentedJSON(200, gin.H{
		"code":    200,
		"success": true,
		"message": "ok",
		"data":    builder.BuildUser(user),
	})
}

// 用户详情
func (c *UserController) Show(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		user := cache.UserCache.Get(gDto.ID)
		if user != nil && user.Status != model.StatusDeleted {
			c.Success(ctx, builder.BuildUser(user))
		} else {
			c.Fail(ctx, util.NewErrorMsg("用户不存在"))
		}
	}
}

// Update 用户资料编辑
func (c *UserController) Update(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, util.ErrorNotLogin)
		return
	}

	var dto form.UserUpdateForm
	if c.BindAndValidate(ctx, &dto) {
		if len(dto.Website) > 0 && util.IsValidateUrl(dto.Website) != nil {
			c.Fail(ctx, util.NewErrorMsg("个人主页地址错误"))
			return
		}
		err := service.UserService.Updates(user.ID, map[string]interface{}{
			"nickname":    dto.Nickname,
			"avatar":      dto.Avatar,
			"website":     dto.Website,
			"description": dto.Description,
		})
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, nil)
	}
}

// GetScoreRank 积分排行
func (c *UserController) GetScoreRank(ctx *gin.Context) {
	userScores := service.UserScoreService.Find(sqlcnd.NewSqlCnd().Desc("score").Limit(10))
	var results []*model.UserInfo
	for _, userScore := range userScores {
		results = append(results, builder.BuildUserDefaultIfNull(userScore.UserId))
	}
	c.Success(ctx, results)
}

// GetScorelogs 用户积分记录
func (c *UserController) GetScorelogs(ctx *gin.Context) {
	page := form.FormValueIntDefault(ctx, "page", 1)
	user := c.GetCurrentUser(ctx)

	logs, paging := service.UserScoreLogService.List(sqlcnd.NewSqlCnd().
		Eq("user_id", user.ID).
		Page(page, 20).Desc("id"))

	c.Success(ctx, gin.H{
		"results": logs,
		"paging":  paging,
	})
}

// GetNotificationsRecent 获取最近3条未读消息
func (c *UserController) GetNotificationsRecent(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var count int64 = 0
	var notifications []model.Notification
	if user != nil {
		count = service.NotificationService.GetUnReadCount(user.ID)
		notifications = service.NotificationService.Find(sqlcnd.NewSqlCnd().Eq("user_id", user.ID).Eq("status", model.NotificationStatusUnread).Limit(3).Desc("id"))
	}
	data := make(map[string]interface{})
	data["count"] = count
	data["notifications"] = builder.BuildNotifications(notifications)
	c.Success(ctx, data)
}

// GetNotifications 用户通知
func (c *UserController) GetNotifications(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	page := form.FormValueIntDefault(ctx, "page", 1)

	messages, paging := service.NotificationService.List(sqlcnd.NewSqlCnd().
		Eq("user_id", user.ID).
		Page(page, 20).Desc("id"))

	// 全部标记为已读
	service.NotificationService.MarkRead(user.ID)

	c.Success(ctx, gin.H{
		"results": builder.BuildNotifications(messages),
		"paging":  paging,
	})
}

// GetFavorites get favorites
func (c *UserController) GetFavorites(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	cursor := form.FormValueInt64Default(ctx, "cursor", 0)

	// 查询列表
	var favorites []model.Favorite
	if cursor > 0 {
		favorites = service.FavoriteService.Find(sqlcnd.NewSqlCnd().Where("user_id = ? and id < ?",
			user.ID, cursor).Desc("id").Limit(20))
	} else {
		favorites = service.FavoriteService.Find(sqlcnd.NewSqlCnd().Where("user_id = ?", user.ID).Desc("id").Limit(20))
	}

	if len(favorites) > 0 {
		cursor = favorites[len(favorites)-1].ID
	}

	c.Success(ctx, gin.H{
		"results": builder.BuildFavorites(favorites),
		"cursor":  cursor,
	})
}

// GetRecentWatchers 关注该用户的人
func (c *UserController) GetRecentWatchers(ctx *gin.Context) {
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		userWatchers := service.UserWatchService.Recent(gDto.ID, 10)
		var users []model.UserInfo
		for _, userWatcher := range userWatchers {
			userInfo := builder.BuildUserById(userWatcher.WatcherID)
			if userInfo != nil {
				users = append(users, *userInfo)
			}
		}
		c.Success(ctx, users)
	}
}

// Watch 关注
func (c *UserController) Watch(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var gDto form.GeneralGetDto
	if c.BindAndValidate(ctx, &gDto) {
		err := service.UserWatchService.Watch(gDto.ID, user.ID)
		if err != nil {
			c.Fail(ctx, util.FromError(err))
			return
		}
		c.Success(ctx, nil)
	}
}

// GetWatched 是否关注了
func (c *UserController) GetWatched(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	userID := form.FormValueInt64Default(ctx, "userId", 0)

	data := map[string]interface{}{}
	if user == nil || userID <= 0 {
		data["watched"] = false
	} else {
		tmp := service.UserWatchService.GetBy(userID, user.ID)
		data["watched"] = tmp != nil
	}
	c.Success(ctx, data)
}

// Delete 取消收藏
func (c *UserController) WatchDelete(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	userID := form.FormValueInt64Default(ctx, "userId", 0)

	tmp := service.UserWatchService.GetBy(userID, user.ID)
	if tmp != nil {
		service.UserWatchService.Delete(tmp.ID)
	}
	c.Success(ctx, nil)
}

// UpdateAvatar 修改头像
func (c *UserController) UpdateAvatar(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	avatar := strings.TrimSpace(ctx.Request.FormValue("avatar"))

	err := service.UserService.UpdateAvatar(user.ID, avatar)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// SetUsername 设置用户名
func (c *UserController) SetUsername(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	username := strings.TrimSpace(ctx.Request.FormValue("username"))

	err := service.UserService.SetUsername(user.ID, username)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// SetEmail 设置邮箱
func (c *UserController) SetEmail(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	email := strings.TrimSpace(ctx.Request.FormValue("email"))

	err := service.UserService.SetEmail(user.ID, email)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// SetPassword 设置密码
func (c *UserController) SetPassword(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)

	var (
		password   = strings.TrimSpace(ctx.Request.FormValue("password"))
		rePassword = strings.TrimSpace(ctx.Request.FormValue("rePassword"))
	)

	err := service.UserService.SetPassword(user.ID, password, rePassword)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}

// ChangePassword 更改密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	var (
		oldPassword = ctx.Request.FormValue("oldPassword")
		password    = ctx.Request.FormValue("password")
		rePassword  = ctx.Request.FormValue("rePassword")
	)
	err := service.UserService.UpdatePassword(user.ID, oldPassword, password, rePassword)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	c.Success(ctx, nil)
}
