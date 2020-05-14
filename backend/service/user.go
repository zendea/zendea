package service

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"

	"zendea/cache"
	"zendea/dao"
	"zendea/form"
	"zendea/model"
	"zendea/util"
	"zendea/util/avatar"
	"zendea/util/log"
	"zendea/util/sqlcnd"
	"zendea/util/uploader"
)

type ScanUserCallback func(users []model.User)

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (s *userService) Get(id int64) *model.User {
	return dao.UserDao.Get(id)
}

func (s *userService) Take(where ...interface{}) *model.User {
	return dao.UserDao.Take(where...)
}

func (s *userService) Find(cnd *sqlcnd.SqlCnd) []model.User {
	return dao.UserDao.Find(cnd)
}

func (s *userService) FindOne(cnd *sqlcnd.SqlCnd) *model.User {
	return dao.UserDao.FindOne(cnd)
}

func (s *userService) List(cnd *sqlcnd.SqlCnd) (list []model.User, paging *sqlcnd.Paging) {
	return dao.UserDao.List(cnd)
}

func (s *userService) Count(cnd *sqlcnd.SqlCnd) int {
	return dao.UserDao.Count(cnd)
}

func (s *userService) Update(dto form.UserUpdateForm) error {
	err := dao.UserDao.Updates(dto.ID, map[string]interface{}{
		"nickname":    dto.Nickname,
		"description": dto.Description,
		"level":       dto.Level,
		"update_time": util.NowTimestamp(),
	})
	cache.UserCache.Invalidate(dto.ID)

	return err
}

func (s *userService) Updates(id int64, columns map[string]interface{}) error {
	err := dao.UserDao.Updates(id, columns)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := dao.UserDao.UpdateColumn(id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

func (s *userService) Delete(id int64) {
	dao.UserDao.Delete(id)
	cache.UserCache.Invalidate(id)
}

// 获取当前登录用户
func (s *userService) GetCurrent(ctx *gin.Context) *model.User {
	userTmp, _ := ctx.Get(viper.GetString("jwt.identity_key"))
	if userTmp == nil {
		return nil
	}

	user := cache.UserCache.Get(userTmp.(model.UserClaims).ID)
	if user == nil || user.Status != model.StatusOk {
		return nil
	}
	return user
}

// Scan 扫描
func (s *userService) Scan(cb ScanUserCallback) {
	var cursor int64
	for {
		list := dao.UserDao.Find(sqlcnd.NewSqlCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		cb(list)
	}
}

// Create 注册
func (s *userService) Create(username, email, nickname, password, rePassword string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	nickname = strings.TrimSpace(nickname)

	// 验证用户名
	if len(username) == 0 {
		return nil, errors.New("用户名不能为空")
	}

	// 验证密码
	err := util.IsValidatePassword(password, rePassword)
	if err != nil {
		return nil, err
	}

	// 验证邮箱
	if len(email) > 0 {
		if err := util.IsValidateEmail(email); err != nil {
			return nil, err
		}
		if dao.UserDao.GetByEmail(email) != nil {
			return nil, errors.New("邮箱：" + email + " 已被占用")
		}
	} else {
		return nil, errors.New("请输入邮箱")
	}

	// 验证用户名
	if len(username) > 0 {
		if err := util.IsValidateUsername(username); err != nil {
			return nil, err
		}
		if s.isUsernameExists(username) {
			return nil, errors.New("用户名：" + username + " 已被占用")
		}
	}

	user := &model.User{
		Username:   util.SqlNullString(username),
		Email:      util.SqlNullString(email),
		Nickname:   nickname,
		Password:   util.EncodePassword(password),
		Status:     model.StatusOk,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}

	err = dao.Tx(func(tx *gorm.DB) error {
		if err := dao.UserDao.Create(user); err != nil {
			return err
		}

		avatarUrl, err := s.HandleAvatar(user.ID, "")
		if err != nil {
			return err
		}

		if err := dao.UserDao.UpdateColumn(user.ID, "avatar", avatarUrl); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return user, nil
}

// SignInByLoginSource 第三方账号登录
func (s *userService) SignInByLoginSource(loginSource *model.LoginSource) (*model.User, error) {
	user := s.Get(loginSource.UserID.Int64)
	if user != nil {
		if user.Status != model.StatusOk {
			return nil, errors.New("用户已被禁用")
		}
		return user, nil
	}

	var website string
	var description string
	if loginSource.TargetType == model.LoginSourceTypeGithub {
		if blog := gjson.Get(loginSource.ExtraData, "blog"); blog.Exists() && len(blog.String()) > 0 {
			website = blog.String()
		} else if htmlUrl := gjson.Get(loginSource.ExtraData, "html_url"); htmlUrl.Exists() && len(htmlUrl.String()) > 0 {
			website = htmlUrl.String()
		}

		description = gjson.Get(loginSource.ExtraData, "bio").String()
	}

	user = &model.User{
		Username:    sql.NullString{},
		Nickname:    loginSource.Nickname,
		Status:      model.StatusOk,
		Website:     website,
		Description: description,
		CreateTime:  util.NowTimestamp(),
		UpdateTime:  util.NowTimestamp(),
	}
	err := dao.Tx(func(tx *gorm.DB) error {
		if err := dao.UserDao.Create(user); err != nil {
			return err
		}

		if err := dao.LoginSourceDao.UpdateColumn(loginSource.ID, "user_id", user.ID); err != nil {
			return err
		}

		avatarUrl, err := s.HandleAvatar(user.ID, loginSource.Avatar)
		if err != nil {
			return err
		}

		if err := dao.UserDao.UpdateColumn(user.ID, "avatar", avatarUrl); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, util.FromError(err)
	}
	cache.UserCache.Invalidate(user.ID)
	return user, nil
}

// HandleAvatar 处理头像，优先级如下：1. 如果第三方登录带有来头像；2. 生成随机默认头像
// avatar: 第三方登录带过来的头像
func (s *userService) HandleAvatar(userId int64, avatarUrl string) (string, error) {
	if len(avatarUrl) > 0 {
		return uploader.CopyImage(avatarUrl)
	}

	avatarBytes, err := avatar.Generate(userId)
	if err != nil {
		return "", err
	}
	return uploader.PutImage(avatarBytes)
}

// isEmailExists 邮箱是否存在
func (s *userService) isEmailExists(email string) bool {
	if len(email) == 0 { // 如果邮箱为空，那么就认为是不存在
		return false
	}
	return dao.UserDao.GetByEmail(email) != nil
}

// isUsernameExists 用户名是否存在
func (s *userService) isUsernameExists(username string) bool {
	return dao.UserDao.GetByUsername(username) != nil
}

// SetAvatar 更新头像
func (s *userService) UpdateAvatar(userId int64, avatar string) error {
	return s.UpdateColumn(userId, "avatar", avatar)
}

// SetUsername 设置用户名
func (s *userService) SetUsername(userId int64, username string) error {
	username = strings.TrimSpace(username)
	if err := util.IsValidateUsername(username); err != nil {
		return err
	}

	user := s.Get(userId)
	if len(user.Username.String) > 0 {
		return errors.New("你已设置了用户名，无法重复设置。")
	}
	if s.isUsernameExists(username) {
		return errors.New("用户名：" + username + " 已被占用")
	}
	return s.UpdateColumn(userId, "username", username)
}

// SetEmail 设置密码
func (s *userService) SetEmail(userId int64, email string) error {
	email = strings.TrimSpace(email)
	if err := util.IsValidateEmail(email); err != nil {
		return err
	}
	if s.isEmailExists(email) {
		return errors.New("邮箱：" + email + " 已被占用")
	}
	return s.UpdateColumn(userId, "email", email)
}

// SetPassword 设置密码
func (s *userService) SetPassword(userId int64, password, rePassword string) error {
	if err := util.IsValidatePassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)
	if len(user.Password) > 0 {
		return errors.New("你已设置了密码，如需修改请前往修改页面。")
	}
	password = util.EncodePassword(password)
	return s.UpdateColumn(userId, "password", password)
}

// UpdatePassword 修改密码
func (s *userService) UpdatePassword(userId int64, oldPassword, password, rePassword string) error {
	if err := util.IsValidatePassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)

	if len(user.Password) == 0 {
		return errors.New("你没设置密码，请先设置密码")
	}

	if !util.ValidatePassword(user.Password, oldPassword) {
		return errors.New("旧密码验证失败")
	}

	return s.UpdateColumn(userId, "password", util.EncodePassword(password))
}

// IncrTopicCount topic_count + 1
func (s *userService) IncrTopicCount(userId int64) int {
	t := dao.UserDao.Get(userId)
	if t == nil {
		return 0
	}
	topicCount := t.TopicCount + 1
	if err := dao.UserDao.UpdateColumn(userId, "topic_count", topicCount); err != nil {
		log.Error(err.Error())
	} else {
		cache.UserCache.Invalidate(userId)
	}
	return topicCount
}

// IncrCommentCount comment_count + 1
func (s *userService) IncrCommentCount(userId int64) int {
	t := dao.UserDao.Get(userId)
	if t == nil {
		return 0
	}
	commentCount := t.CommentCount + 1
	if err := dao.UserDao.UpdateColumn(userId, "comment_count", commentCount); err != nil {
		log.Error(err.Error())
	} else {
		cache.UserCache.Invalidate(userId)
	}
	return commentCount
}

// SyncUserCount 同步用户计数
func (s *userService) SyncUserCount() {
	s.Scan(func(users []model.User) {
		for _, user := range users {
			topicCount := dao.TopicDao.Count(sqlcnd.NewSqlCnd().Eq("user_id", user.ID).Eq("status", model.StatusOk))
			commentCount := dao.CommentDao.Count(sqlcnd.NewSqlCnd().Eq("user_id", user.ID).Eq("status", model.StatusOk))
			_ = dao.UserDao.UpdateColumn(user.ID, "topic_count", topicCount)
			_ = dao.UserDao.UpdateColumn(user.ID, "comment_count", commentCount)
			cache.UserCache.Invalidate(user.ID)
		}
	})
}

var (
	errInvalidAccount = errors.New("账号或密码错误")
	errInvalidCode    = errors.New("请输入正确验证码")
	errAccountLocked  = errors.New("账号已被锁定,请联系管理员")
)

// VerifyAndReturnUserInfo - login and return user info
func (s *userService) VerifyAndReturnUserInfo(username, password string) (bool, error, model.User) {
	var userModel *model.User = nil
	if err := util.IsValidateEmail(username); err == nil { // 如果用户输入的是邮箱
		userModel = dao.UserDao.GetByEmail(username)
	} else {
		userModel = dao.UserDao.GetByUsername(username)
	}

	if userModel == nil {
		return false, errInvalidAccount, model.User{}
	}
	// Account not exits
	if userModel.ID < 1 {
		return false, errInvalidAccount, model.User{}
	}

	if !util.ValidatePassword(userModel.Password, password) {
		log.Error("password wrong: username=%s", userModel.Nickname)
		return false, errInvalidAccount, model.User{}
	}

	return true, nil, *userModel

}
