package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/urfave/cli"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"zendea/cache"
	"zendea/cron"
	"zendea/dao"
	"zendea/middleware"
	"zendea/router"
	"zendea/util/log"
)

var CmdWeb = cli.Command{
	Name:  "web",
	Usage: "Start Zendea API",
	Description: `A free, open-source, self-hosted forum software written in Go`,
	Action: runWeb,
	Flags:  []cli.Flag{},
}

func runWeb(*cli.Context) {
	// do something

	//1.Set up log level
	zerolog.SetGlobalLevel(zerolog.Level(loglevel))

	//2.Set up configuration
	viper.SetConfigFile(conf)
	content, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read conf file fail: %s", err.Error()))
	}
	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		log.Fatal(fmt.Sprintf("Parse conf file fail: %s", err.Error()))
	}

	//3.Set up run mode
	mode := viper.GetString("mode")
	gin.SetMode(mode)

	//4.Set up database connection
	dao.Setup()

	//5.Set up cache
	cache.Setup()

	////5.Set up cron
	cron.Setup()

	//6.Initialize language
	middleware.InitLang()

	engine := gin.Default()
	router.Setup(engine, cors)
	engine.Run(":" + viper.GetString("base.port"))
}