package service

import (
	"strconv"

	"zendea/util/log"
	"zendea/util/sqlcnd"
)

var StatisticService = newStatisticService()

func newStatisticService() *statisticService {
	return &statisticService{}
}

type statisticService struct {
	running bool // is in running
}

func (s *statisticService) GenerateData() {
	if s.running {
		log.Info("statistic is in building")
		return
	}

	s.running = true
	defer func() {
		s.running = false
	}()

	var (
		statUserCount    = strconv.Itoa(UserService.Count(sqlcnd.NewSqlCnd()))
		statTopicCount   = strconv.Itoa(TopicService.Count(sqlcnd.NewSqlCnd()))
		statCommentCount = strconv.Itoa(CommentService.Count(sqlcnd.NewSqlCnd()))
	)

	SettingService.Set("statUserCount", statUserCount, "社区会员", "社区会员总数")
	SettingService.Set("statTopicCount", statTopicCount, "帖子数", "主题总数")
	SettingService.Set("statCommentCount", statCommentCount, "回帖数", "回帖总数")
}
