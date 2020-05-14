package service

import (
	"errors"

	"github.com/jinzhu/gorm"

	"zendea/dao"
	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var UserWatchService = newUserWatchService()

func newUserWatchService() *userWatchService {
	return &userWatchService{}
}

type userWatchService struct {
}

func (s *userWatchService) Get(id int64) *model.UserWatch {
	return dao.UserWatchDao.Get(id)
}

func (s *userWatchService) Take(where ...interface{}) *model.UserWatch {
	return dao.UserWatchDao.Take(where...)
}

func (s *userWatchService) Find(cnd *sqlcnd.SqlCnd) []model.UserWatch {
	return dao.UserWatchDao.Find(cnd)
}

func (s *userWatchService) FindOne(cnd *sqlcnd.SqlCnd) *model.UserWatch {
	return dao.UserWatchDao.FindOne(cnd)
}

func (s *userWatchService) List(cnd *sqlcnd.SqlCnd) (list []model.UserWatch, paging *sqlcnd.Paging) {
	return dao.UserWatchDao.List(cnd)
}

func (s *userWatchService) Create(t *model.UserWatch) error {
	return dao.UserWatchDao.Create(t)
}

func (s *userWatchService) Update(t *model.UserWatch) error {
	return dao.UserWatchDao.Update(t)
}

func (s *userWatchService) Updates(id int64, columns map[string]interface{}) error {
	return dao.UserWatchDao.Updates(id, columns)
}

func (s *userWatchService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.UserWatchDao.UpdateColumn(id, name, value)
}

func (s *userWatchService) Delete(id int64) {
	dao.UserWatchDao.Delete(id)
}

func (s *userWatchService) GetBy(userID int64, watcherID int64) *model.UserWatch {
	return dao.UserWatchDao.Take("user_id = ? and watcher_id = ?",
		userID, watcherID)
}

// 统计数量
func (s *userWatchService) Count(userId int64) int64 {
	var count int64 = 0
	dao.DB().Model(&model.UserWatch{}).Where("user_id = ?", userId).Count(&count)
	return count
}

// 最近关注
func (s *userWatchService) Recent(userId int64, count int) []model.UserWatch {
	return s.Find(sqlcnd.NewSqlCnd().Eq("user_id", userId).Desc("id").Limit(count))
}

func (s *userWatchService) Watch(userID int64, watcherID int64) error {
	if userID == watcherID {
		return errors.New("不能自己关注自己")
	}
	user := dao.UserDao.Get(userID)
	if user == nil || user.Status != model.StatusOk {
		return errors.New("用户不存在")
	}

	// 判断是否已经点赞了
	userWatch := dao.UserWatchDao.Take("user_id = ? and watcher_id = ?", userID, watcherID)
	if userWatch != nil {
		return errors.New("已关注")
	}

	return dao.Tx(func(tx *gorm.DB) error {
		// 点赞
		userWatch := &model.UserWatch{
			UserID:     userID,
			WatcherID:  watcherID,
			CreateTime: util.NowTimestamp(),
		}
		err := dao.UserWatchDao.Create(userWatch)
		if err != nil {
			return err
		}
		// 发送点赞通知
		NotificationService.SendUserWatchNotification(userWatch)
		return nil
		// return dao.DB().Model(&topic).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
	})
}
