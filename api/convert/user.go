package convert

import (
	"strconv"

	"zendea/cache"
	"zendea/model"
	"zendea/util"
	"zendea/util/avatar"
)

func ToUserDefaultIfNull(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	if user == nil {
		user = &model.User{}
		user.ID = id
		user.Username = util.SqlNullString(strconv.FormatInt(id, 10))
		user.Avatar = avatar.DefaultAvatar
		user.CreateTime = util.NowTimestamp()
	}
	return ToUser(user)
}

func ToUserById(id int64) *model.UserInfo {
	user := cache.UserCache.Get(id)
	return ToUser(user)
}

func ToUser(user *model.User) *model.UserInfo {
	if user == nil {
		return nil
	}
	a := user.Avatar
	if len(a) == 0 {
		a = avatar.DefaultAvatar
	}
	levelName := "普通用户"
	if user.Level == model.UserLevelAdmin {
		levelName = "管理员"
	}
	ret := &model.UserInfo{
		Id:           user.ID,
		Username:     user.Username.String,
		Nickname:     user.Nickname,
		Avatar:       a,
		Email:        user.Email.String,
		Level:        user.Level,
		LevelName:    levelName,
		Website:      user.Website,
		Description:  user.Description,
		TopicCount:   user.TopicCount,
		CommentCount: user.CommentCount,
		PasswordSet:  len(user.Password) > 0,
		Status:       user.Status,
		CreateTime:   user.CreateTime,
	}
	if user.Status == model.StatusDeleted {
		ret.Username = "blacklist"
		ret.Nickname = "黑名单用户"
		ret.Avatar = avatar.DefaultAvatar
		ret.Email = ""
		ret.Website = ""
		ret.Description = ""
	} else {
		ret.Score = cache.UserCache.GetScore(user.ID)
	}
	return ret
}

func ToUsers(users []model.User) []model.UserInfo {
	if len(users) == 0 {
		return nil
	}
	var responses []model.UserInfo
	for _, user := range users {
		item := ToUser(&user)
		if item != nil {
			responses = append(responses, *item)
		}
	}
	return responses
}
