package mock

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"

	"zendea/model"
	"zendea/dao"
	"zendea/util"
)

var (
	// 头像假数据
	avatars = []string{
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
		"https://cdn.learnku.com/uploads/avatars/7850_1481780622.jpeg!/both/380x380",
	}
)

func userFactory(i int) *factory.Factory {

	password := "123456"

	u := &model.User{
		Username:   util.SqlNullString(""),
		Password:   util.EncodePassword(password),
		Status:     model.StatusOk,
		CreateTime: util.NowTimestamp(),
		UpdateTime: util.NowTimestamp(),
	}

	r := RandInt(0, len(avatars)-1)

	return factory.NewFactory(
		u,
	).Attr("Nickname", func(args factory.Args) (interface{}, error) {
		return fmt.Sprintf("user-%d", i+1), nil
	}).Attr("Avatar", func(args factory.Args) (interface{}, error) {
		return avatars[r], nil
	}).Attr("Description", func(args factory.Args) (interface{}, error) {
		paragraph := randomdata.Paragraph()

		if len(paragraph) >= 70 {
			paragraph = paragraph[:70]
		}
		return paragraph, nil
	}).Attr("Email", func(args factory.Args) (interface{}, error) {
		if i == 0 {
			return util.SqlNullString("1@test.com"), nil
		}
		if i == 1 {
			return util.SqlNullString("2@test.com"), nil
		}
		return util.SqlNullString(randomdata.Email()), nil
	}).Attr("Level", func(args factory.Args) (interface{}, error) {
		if i == 0 {
			return 10, nil
		}
		if i == 1 {
			return 10, nil
		}
		return 0, nil
	})

}

// UserTableSeeder -
func UserTableSeeder(needCleanTable bool) {
	if needCleanTable {
		dropAndCreateTable(&model.User{})
	}

	for i := 0; i < 10; i++ {
		user := userFactory(i).MustCreate().(*model.User)
		fmt.Println("Email:", user.Email)
		if err := dao.UserDao.Create(user); err != nil {
			fmt.Printf("mock user error： %v\n", err)
		}
	}
}
