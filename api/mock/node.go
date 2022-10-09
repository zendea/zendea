package mock

import (
	"fmt"

	"zendea/model"
	"zendea/dao"
)

// NodeTableSeeder -
func NodeTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Node{})
	}

	ns := []*model.Node{
		{
			Name:        "分享",
			Description: "分享创造，分享发现",
			Status:      0,
			TopicCount:  0,
		},
		{
			Name:        "教程",
			Description: "开发技巧、推荐扩展包等",
			Status:      0,
			TopicCount:  0,
		},
		{
			Name:        "问答",
			Description: "请保持友善，互帮互助",
			Status:      0,
			TopicCount:  0,
		},
		{
			Name:        "公告",
			Description: "站点公告",
			Status:      0,
			TopicCount:  0,
		},
	}

	for _, n := range ns {
		if err := dao.NodeDao.Create(n); err != nil {
			fmt.Printf("mock category error： %v\n", err)
		}
	}
}