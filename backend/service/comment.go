package service

import (
	"errors"
	"strings"

	"zendea/dao"
	"zendea/model"
	"zendea/form"
	"zendea/util"
	"zendea/util/sqlcnd"
)

var CommentService = newCommentService()

func newCommentService() *commentService {
	return &commentService{}
}

type commentService struct {
}

func (s *commentService) Get(id int64) *model.Comment {
	return dao.CommentDao.Get(id)
}

func (s *commentService) Take(where ...interface{}) *model.Comment {
	return dao.CommentDao.Take(where...)
}

func (s *commentService) Find(cnd *sqlcnd.SqlCnd) []model.Comment {
	return dao.CommentDao.Find(cnd)
}

func (s *commentService) FindOne(cnd *sqlcnd.SqlCnd) *model.Comment {
	return dao.CommentDao.FindOne(cnd)
}

func (s *commentService) List(cnd *sqlcnd.SqlCnd) (list []model.Comment, paging *sqlcnd.Paging) {
	return dao.CommentDao.List(cnd)
}

func (s *commentService) Count(cnd *sqlcnd.SqlCnd) int {
	return dao.CommentDao.Count(cnd)
}

func (s *commentService) Update(dto form.CommentUpdateForm) error {
	err := dao.NodeDao.Updates(dto.ID, map[string]interface{}{
		"status":      dto.Status,
	})
	
	return err
}

func (s *commentService) Delete(id int64) error {
	return dao.CommentDao.UpdateColumn(id, "status", model.StatusDeleted)
}

func (s *commentService) Create(dto form.CommentCreateForm) (*model.Comment, error) {
	if dto.EntityID <= 0 {
		return nil, errors.New("参数EntityId非法")
	}

	comment := &model.Comment{
		UserId:      dto.UserID,
		EntityType:  dto.EntityType,
		EntityId:    dto.EntityID,
		Content:     dto.Content,
		ContentType: model.ContentTypeMarkdown,
		QuoteId:     dto.QuoteID,
		Status:      model.StatusOk,
		CreateTime:  util.NowTimestamp(),
	}
	if err := dao.CommentDao.Create(comment); err != nil {
		return nil, errors.New("创建评论失败")
	}

	// 更新帖子最后回复时间
	if dto.EntityType == model.EntityTypeTopic {
		TopicService.OnComment(dto.EntityID, dto.UserID, util.NowTimestamp())
	}

	// 用户跟帖计数
	UserService.IncrCommentCount(dto.UserID)
	// 获得积分
	UserScoreService.IncrementPostCommentScore(comment)
	// 发送通知
	NotificationService.SendCommentNotification(comment)

	return comment, nil
}
// 发表评论
func (s *commentService) Publish(userId int64, createForm *form.CommentCreateForm) (*model.Comment, error) {
	createForm.Content = strings.TrimSpace(createForm.Content)

	if len(createForm.EntityType) == 0 {
		return nil, errors.New("参数非法")
	}
	if createForm.EntityID <= 0 {
		return nil, errors.New("参数非法")
	}
	if len(createForm.Content) == 0 {
		return nil, errors.New("请输入评论内容")
	}

	contentType := createForm.ContentType
	if contentType == "" {
		contentType = model.ContentTypeMarkdown
	}

	comment := &model.Comment{
		UserId:      userId,
		EntityType:  createForm.EntityType,
		EntityId:    createForm.EntityID,
		Content:     createForm.Content,
		ContentType: contentType,
		QuoteId:     createForm.QuoteID,
		Status:      model.StatusOk,
		CreateTime:  util.NowTimestamp(),
	}
	if err := dao.CommentDao.Create(comment); err != nil {
		return nil, err
	}

	// 更新帖子最后回复时间
	if createForm.EntityType == model.EntityTypeTopic {
		TopicService.OnComment(createForm.EntityID, userId, util.NowTimestamp())
	}

	// 用户跟帖计数
	UserService.IncrCommentCount(userId)
	// 获得积分
	UserScoreService.IncrementPostCommentScore(comment)
	// 发送通知
	NotificationService.SendCommentNotification(comment)

	return comment, nil
}

// 列表
func (s *commentService) GetComments(entityType string, entityId int64, cursor int64) (comments []model.Comment, nextCursor int64) {
	cnd := sqlcnd.NewSqlCnd().
		Eq("entity_type", entityType).
		Eq("entity_id", entityId).
		Eq("status", model.StatusOk).
		Asc("id").Limit(50)

	if cursor > 0 {
		cnd.Gt("id", cursor)
	}
	comments = dao.CommentDao.Find(cnd)
	if len(comments) > 0 {
		nextCursor = comments[len(comments)-1].ID
	} else {
		nextCursor = cursor
	}
	return
}
