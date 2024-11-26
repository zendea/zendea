package mock

import (
	"fmt"

	"zendea/model"
	"zendea/dao"
)

// SettingTableSeeder -
func SettingTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Setting{})
	}

	ns := []*model.Setting{
		{
			Key:   "defaultNodeId",
			Value: "1",
		},
		{
			Key:   "siteTitle",
			Value: "Zendea",
		},
		{
			Key:    "siteDescription",
			Value: "小而美的开发者社区",
		},
	}

	for _, n := range ns {
		if err := dao.SettingDao.Create(n); err != nil {
			fmt.Printf("mock category error： %v\n", err)
		}
	}
}