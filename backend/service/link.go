package service

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"

	"zendea/dao"
	"zendea/model"
	"zendea/form"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var LinkService = newLinkService()

func newLinkService() *linkService {
	return &linkService{}
}

type linkService struct {
}

func (s *linkService) Get(id int64) *model.Link {
	return dao.LinkDao.Get(id)
}

func (s *linkService) Find(cnd *sqlcnd.SqlCnd) []model.Link {
	return dao.LinkDao.Find(cnd)
}

func (s *linkService) List(cnd *sqlcnd.SqlCnd) (list []model.Link, paging *sqlcnd.Paging) {
	return dao.LinkDao.List(cnd)
}

func (s *linkService) Create(dto form.LinkCreateForm) (*model.Link, error) {
	link := &model.Link{
		Title:   dto.Title,
		Url:     dto.URL,
		Summary: dto.Summary,
		Logo:    dto.Logo,
		CreateTime: util.NowTimestamp(),
	}
	err := dao.Tx(func(tx *gorm.DB) error {
		err := dao.LinkDao.Create(link)
		if err != nil {
			return err
		}
		return nil
	})
	return link, err
}

func (s *linkService) Update(dto form.LinkUpdateForm) error {
	err := dao.LinkDao.Updates(dto.ID, map[string]interface{}{
		"title":       dto.Title,
		"url":         dto.URL,
		"summary":     dto.Summary,
		"logo":        dto.Logo,
		"status":      dto.Status,
		"update_time": util.NowTimestamp(),
	})
	
	return err
}

func (s *linkService) Delete(id int64) {
	dao.LinkDao.Delete(id)
}

// 提交友情链接
func (s *linkService) Submit(url, title, summary, logo string) (link *model.Link, err error) {
	url = strings.TrimSpace(url)
	title = strings.TrimSpace(title)
	summary = strings.TrimSpace(summary)
	logo = strings.TrimSpace(logo)

	if len(url) == 0 {
		return nil, errors.New("网址不能为空")
	}
	if len(title) == 0 {
		return nil, errors.New("标题不能为空")
	}

	link = &model.Link{
		Url:        url,
		Title:      title,
		Summary:    summary,
		Logo:       logo,
		Status:     model.StatusPending,
		CreateTime: util.NowTimestamp(),
	}

	err = dao.Tx(func(tx *gorm.DB) error {
		err := dao.LinkDao.Create(link)
		if err != nil {
			return err
		}
		return nil
	})

	return
}
