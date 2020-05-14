package service

import (
	"zendea/config"
	"zendea/model"
)

var AppinfoService = newAppinfoService()

func newAppinfoService() *appinfoService {
	return &appinfoService{}
}

type appinfoService struct {
}

func (s *appinfoService) GetAppinfo() *model.AppData {
	return &model.AppData{
		Name:           config.AppName,
		Version:        config.AppVersion,
		UserLevelAdmin: model.UserLevelAdmin,
	}
}
