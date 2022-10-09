package mock

import (
	"fmt"
	"zendea/model"
	"zendea/dao"
	"zendea/util"
	"github.com/Pallinder/go-randomdata"
)


// CommentTableSeeder -
func CommentTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Comment{})
	}

	for i := 0; i < 300; i++ {
		comment := commentFactory()
		if err := dao.CommentDao.Create(comment); err != nil {
			fmt.Printf("mock comment error： %v\n", err)
		}
	}
}

func commentFactory() *model.Comment {
	paragraph := randomdata.Paragraph()
	
	randUID := int64(RandInt(1, 10))
	randTID := int64(RandInt(1, 30))

	return &model.Comment{
		Content:    paragraph,
		UserId:     randUID,
		EntityType: model.EntityTypeTopic,
		EntityId:   randTID,
		CreateTime: util.NowTimestamp(),
	}
}