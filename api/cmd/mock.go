package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/urfave/cli"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"zendea/dao"
	"zendea/util/log"
	"zendea/mock"
)


var CmdMock = cli.Command{
	Name:  "mock",
	Usage: "Mock Data",
	Description: `A free, open-source, self-hosted forum software written in Go`,
	Action: runMock,
	Flags:  []cli.Flag{},
}

func runMock(ctx *cli.Context) {
	//1.Set up log level
	zerolog.SetGlobalLevel(zerolog.Level(0))

	conf := "./app.yaml"
	if ctx.IsSet("conf") {
		conf = ctx.String("conf")
	}

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

	dao.Setup()
	log.Info("run mock\n")
	mock.Mock()
}
