package service

import (
	"errors"
	"zendea/dao"
	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
	"zendea/form"
)

var SectionService = newSectionService()

func newSectionService() *sectionService {
	return &sectionService{}
}

type sectionService struct {
}

func (s *sectionService) Get(id int64) *model.Section {
	return dao.SectionDao.Get(id)
}

func (s *sectionService) List(cnd *sqlcnd.SqlCnd) (list []model.Section, paging *sqlcnd.Paging) {
	return dao.SectionDao.List(cnd)
}

func (s *sectionService) Count(cnd *sqlcnd.SqlCnd) int {
	return dao.SectionDao.Count(cnd)
}

func (s *sectionService) Create(dto form.SectionCreateForm) (*model.Section, error) {
	section := &model.Section{
		Name:       dto.Name,
		SortNo:     dto.SortNo,
		CreateTime: util.NowTimestamp(),
	}
	if err := dao.SectionDao.Create(section); err != nil {
		return nil, errors.New("创建分类失败")
	}
	return section, nil
}

func (s *sectionService) Update(dto form.SectionUpdateForm) error {
	err := dao.SectionDao.Updates(dto.ID, map[string]interface{}{
		"name":        dto.Name,
		"sort_no":     dto.SortNo,
		"update_time": util.NowTimestamp(),
	})
	
	return err
}

func (s *sectionService) Delete(id int64) {
	dao.SectionDao.Delete(id)
}

func (s *sectionService) GetSectionNodes(sectionId int64) []model.Node {
	return dao.NodeDao.Find(sqlcnd.NewSqlCnd().Where("section_id = ?", sectionId).Asc("sort_no").Desc("id"))
}

func (s *sectionService) GetSections() []model.Section {
	return dao.SectionDao.Find(sqlcnd.NewSqlCnd().Asc("sort_no").Desc("id"))
}
