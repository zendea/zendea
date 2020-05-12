package service

import (
	"zendea/dao"
	"zendea/model"
	"zendea/util/sqlcnd"
)

var UserScoreLogService = newUserScoreLogService()

func newUserScoreLogService() *userScoreLogService {
	return &userScoreLogService{}
}

type userScoreLogService struct {
}

func (s *userScoreLogService) Get(id int64) *model.UserScoreLog {
	return dao.UserScoreLogDao.Get(id)
}

func (s *userScoreLogService) Take(where ...interface{}) *model.UserScoreLog {
	return dao.UserScoreLogDao.Take(where...)
}

func (s *userScoreLogService) Find(cnd *sqlcnd.SqlCnd) []model.UserScoreLog {
	return dao.UserScoreLogDao.Find(cnd)
}

func (s *userScoreLogService) FindOne(cnd *sqlcnd.SqlCnd) *model.UserScoreLog {
	return dao.UserScoreLogDao.FindOne(cnd)
}

func (s *userScoreLogService) List(cnd *sqlcnd.SqlCnd) (list []model.UserScoreLog, paging *sqlcnd.Paging) {
	return dao.UserScoreLogDao.List(cnd)
}

func (s *userScoreLogService) Create(t *model.UserScoreLog) error {
	return dao.UserScoreLogDao.Create(t)
}

func (s *userScoreLogService) Update(t *model.UserScoreLog) error {
	return dao.UserScoreLogDao.Update(t)
}

func (s *userScoreLogService) Updates(id int64, columns map[string]interface{}) error {
	return dao.UserScoreLogDao.Updates(id, columns)
}

func (s *userScoreLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.UserScoreLogDao.UpdateColumn(id, name, value)
}

func (s *userScoreLogService) Delete(id int64) {
	dao.UserScoreLogDao.Delete(id)
}
