package service

import (
	"strings"

	"zendea/cache"
	"zendea/dao"
	"zendea/model"
	"zendea/util/sqlcnd"
)

type ScanTagCallback func(tags []model.Tag) bool

var TagService = newTagService()

func newTagService() *tagService {
	return &tagService{}
}

type tagService struct {
}

func (s *tagService) Get(id int64) *model.Tag {
	return dao.TagDao.Get(id)
}

func (s *tagService) Take(where ...interface{}) *model.Tag {
	return dao.TagDao.Take(where...)
}

func (s *tagService) Find(cnd *sqlcnd.SqlCnd) []model.Tag {
	return dao.TagDao.Find(cnd)
}

func (s *tagService) FindOne(cnd *sqlcnd.SqlCnd) *model.Tag {
	return dao.TagDao.FindOne(cnd)
}

func (s *tagService) List(cnd *sqlcnd.SqlCnd) (list []model.Tag, paging *sqlcnd.Paging) {
	return dao.TagDao.List(cnd)
}

func (s *tagService) Create(t *model.Tag) error {
	return dao.TagDao.Create(t)
}

func (s *tagService) Update(t *model.Tag) error {
	if err := dao.TagDao.Update(t); err != nil {
		return err
	}
	cache.TagCache.Invalidate(t.ID)
	return nil
}

// 自动完成
func (s *tagService) Autocomplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.TagDao.Find(sqlcnd.NewSqlCnd().Where("status = ? and name like ?",
		model.StatusOk, "%"+input+"%").Limit(6))
}

func (s *tagService) GetOrCreate(name string) (*model.Tag, error) {
	return dao.TagDao.GetOrCreate(name)
}

func (s *tagService) GetByName(name string) *model.Tag {
	return dao.TagDao.GetByName(name)
}

func (s *tagService) GetTags() []model.TagResponse {
	list := dao.TagDao.Find(sqlcnd.NewSqlCnd().Where("status = ?", model.StatusOk))

	var tags []model.TagResponse
	for _, tag := range list {
		tags = append(tags, model.TagResponse{TagId: tag.ID, TagName: tag.Name})
	}
	return tags
}

func (s *tagService) GetTagInIds(tagIds []int64) []model.Tag {
	return dao.TagDao.GetTagInIds(tagIds)
}

// 扫描
func (s *tagService) Scan(cb ScanTagCallback) {
	var cursor int64
	for {
		list := dao.TagDao.Find(sqlcnd.NewSqlCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].ID
		if !cb(list) {
			break
		}
	}
}
