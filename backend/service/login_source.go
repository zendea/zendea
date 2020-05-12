package service

import (
	"database/sql"
	"strconv"
	"strings"

	"zendea/dao"
	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
	"zendea/oauth/github"
	"zendea/oauth/gitee"
	"zendea/oauth/qq"
)

var LoginSourceService = newLoginSourceService()

func newLoginSourceService() *loginSourceService {
	return &loginSourceService{}
}

type loginSourceService struct {
}

func (s *loginSourceService) Get(id int64) *model.LoginSource {
	return dao.LoginSourceDao.Get(id)
}

func (s *loginSourceService) List(cnd *sqlcnd.SqlCnd) (list []model.LoginSource, paging *sqlcnd.Paging) {
	return dao.LoginSourceDao.List(cnd)
}

func (s *loginSourceService) Create(t *model.LoginSource) error {
	return dao.LoginSourceDao.Create(t)
}

func (s *loginSourceService) Update(t *model.LoginSource) error {
	return dao.LoginSourceDao.Update(t)
}

func (s *loginSourceService) Updates(id int64, columns map[string]interface{}) error {
	return dao.LoginSourceDao.Updates(id, columns)
}

func (s *loginSourceService) UpdateColumn(id int64, name string, value interface{}) error {
	return dao.LoginSourceDao.UpdateColumn(id, name, value)
}

func (s *loginSourceService) Delete(id int64) {
	dao.LoginSourceDao.Delete(id)
}

func (s *loginSourceService) GetLoginSource(targetType string, targetID string) *model.LoginSource {
	return dao.LoginSourceDao.Take("target_type = ? and target_id = ?", targetType, targetID)
}

func (s *loginSourceService) GetOrCreate(provider, code, state string) (*model.LoginSource, error) {
	if provider == "github" {
		return s.GetOrCreateByGithub(code, state)
	} else if provider == "gitee" {
		return s.GetOrCreateByGitee(code, state)
	}

	return s.GetOrCreateByQQ(code, state)
}

func (s *loginSourceService) GetOrCreateByGithub(code, state string) (*model.LoginSource, error) {
	userInfo, err := github.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetLoginSource(model.LoginSourceTypeGithub, strconv.FormatInt(userInfo.Id, 10))
	if account != nil {
		return account, nil
	}

	nickname := userInfo.Login
	if len(userInfo.Name) > 0 {
		nickname = strings.TrimSpace(userInfo.Name)
	}

	userInfoJson, _ := util.FormatJson(userInfo)
	account = &model.LoginSource{
		UserID:     sql.NullInt64{},
		Avatar:     userInfo.AvatarUrl,
		Nickname:   nickname,
		TargetType:  model.LoginSourceTypeGithub,
		TargetID:    strconv.FormatInt(userInfo.Id, 10),
		ExtraData:  userInfoJson,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *loginSourceService) GetOrCreateByGitee(code, state string) (*model.LoginSource, error) {
	userInfo, err := gitee.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetLoginSource(model.LoginSourceTypeGitee, strconv.FormatInt(userInfo.Id, 10))
	if account != nil {
		return account, nil
	}

	nickname := userInfo.Login
	if len(userInfo.Name) > 0 {
		nickname = strings.TrimSpace(userInfo.Name)
	}

	userInfoJson, _ := util.FormatJson(userInfo)
	account = &model.LoginSource{
		UserID:     sql.NullInt64{},
		Avatar:     userInfo.AvatarUrl,
		Nickname:   nickname,
		TargetType:  model.LoginSourceTypeGitee,
		TargetID:    strconv.FormatInt(userInfo.Id, 10),
		ExtraData:  userInfoJson,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *loginSourceService) GetOrCreateByQQ(code, state string) (*model.LoginSource, error) {
	userInfo, err := qq.GetUserInfoByCode(code, state)
	if err != nil {
		return nil, err
	}

	account := s.GetLoginSource(model.LoginSourceTypeQQ, userInfo.Unionid)
	if account != nil {
		return account, nil
	}

	userInfoJson, _ := util.FormatJson(userInfo)
	account = &model.LoginSource{
		UserID:     sql.NullInt64{},
		Avatar:     userInfo.FigureurlQQ1,
		Nickname:   strings.TrimSpace(userInfo.Nickname),
		TargetType:  model.LoginSourceTypeQQ,
		TargetID:    userInfo.Unionid,
		ExtraData:  userInfoJson,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}
	err = s.Create(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
