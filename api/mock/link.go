package mock

import (
	"fmt"
	"zendea/model"
	"zendea/dao"
	"strconv"
)

func linkFactory(i int) *model.Link {
	index := strconv.Itoa(i)
	return &model.Link{
		Title: "link title " + index,
		Url:  "https://www.baidu.com/" + index,
	}
}

// LinksTableSeeder -
func LinkTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.Link{})
	}

	for i := 0; i < 6; i++ {
		link := linkFactory(i)
		if err := dao.LinkDao.Create(link); err != nil {
			fmt.Printf("mock link error： %v\n", err)
		}
	}
}
