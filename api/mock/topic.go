package mock

import (
	"fmt"
	"zendea/model"
	"zendea/dao"
	"zendea/util"
	"github.com/Pallinder/go-randomdata"
)


// TopicTableSeeder -
func TopicTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Topic{})
	}

	for i := 0; i < 30; i++ {
		topic := topicFactory()
		if err := dao.TopicDao.Create(topic); err != nil {
			fmt.Printf("mock topic error： %v\n", err)
		}
	}
}

func topicFactory() *model.Topic {
	now := util.NowTimestamp()
	title := randomdata.Country(randomdata.FullCountry)
	paragraph := randomdata.Paragraph()
	
	randUID := int64(RandInt(1, 10))
	randCID := int64(RandInt(1, 4))

	return &model.Topic{
		Title:   title,
		Content: paragraph,
		UserId:  randUID,
		NodeId:  randCID,
		CreateTime: now,
		LastCommentTime: now,
	}
}