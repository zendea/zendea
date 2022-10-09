package mock

import (
	"math/rand"
	"zendea/dao"
)

// dropAndCreateTable - 清空表
func dropAndCreateTable(table interface{}) {
	dao.DB().DropTable(table)
	dao.DB().CreateTable(table)
}

// Mock -
func Mock() {
	UserTableSeeder(true)
	NodeTableSeeder(true)
	TopicTableSeeder(true)
	CommentTableSeeder(true)
	LinkTableSeeder(true)
	SettingTableSeeder(true)
}

func RandInt(min, max int) int {
	if min >= max || min < 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}