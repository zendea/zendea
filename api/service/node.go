package service

import (
	"errors"
	"github.com/jinzhu/gorm"

	"zendea/dao"
	"zendea/form"
	"zendea/model"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var NodeService = newNodeService()

func newNodeService() *nodeService {
	return &nodeService{}
}

type nodeService struct {
}

func (s *nodeService) Get(id int64) *model.Node {
	return dao.NodeDao.Get(id)
}

func (s *nodeService) List(cnd *sqlcnd.SqlCnd) (list []model.Node, paging *sqlcnd.Paging) {
	return dao.NodeDao.List(cnd)
}

func (s *nodeService) Create(dto form.NodeCreateForm) (*model.Node, error) {
	node := &model.Node{
		Name:        dto.Name,
		Description: dto.Description,
		SortNo:      dto.SortNo,
		Status:      dto.Status,
		CreateTime:  util.NowTimestamp(),
	}
	if err := dao.NodeDao.Create(node); err != nil {
		return nil, errors.New("创建节点失败")
	}
	return node, nil
}

func (s *nodeService) Update(dto form.NodeUpdateForm) error {
	err := dao.NodeDao.Updates(dto.ID, map[string]interface{}{
		"name":        dto.Name,
		"description": dto.Description,
		"section_id":  dto.SectionID,
		"sort_no":     dto.SortNo,
		"status":      dto.Status,
		"update_time": util.NowTimestamp(),
	})

	return err
}

func (s *nodeService) Delete(id int64) {
	dao.NodeDao.Delete(id)
}

// 主题数+1
func (s *nodeService) IncrTopicCount(nodeId int64) {
	dao.DB().Model(&model.Node{}).Where("id = ?", nodeId).UpdateColumn("topic_count", gorm.Expr("topic_count + ?", 1))
}

func (s *nodeService) GetRecommendNodes() []model.Node {
	return dao.NodeDao.Find(sqlcnd.NewSqlCnd().Eq("status", model.StatusOk).Asc("sort_no").Desc("id").Limit(3))
}

func (s *nodeService) GetNodes() []model.Node {
	return dao.NodeDao.Find(sqlcnd.NewSqlCnd().Eq("status", model.StatusOk).Asc("sort_no").Desc("id"))
}
