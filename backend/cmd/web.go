package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"zendea/util/log"
	"zendea/dao"
	"zendea/cache"
	"zendea/cron"
	"zendea/middleware"
	"zendea/router"
)

var (
	conf   string
	port     string
	loglevel uint8
	cors     bool
	cluster  bool
	//StartCmd : set up restful api server
	WebStart = &cobra.Command{
		Use:     "web",
		Short:   "Start zendea API server",
		Example: "zendea web -c app.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			usage()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	WebStart.PersistentFlags().StringVarP(&conf, "conf", "c", "./app.yaml", "Start server with provided configuration file")
	WebStart.PersistentFlags().StringVarP(&port, "port", "p", "9527", "Tcp port server listening on")
	WebStart.PersistentFlags().Uint8VarP(&loglevel, "loglevel", "l", 0, "Log level")
	WebStart.PersistentFlags().BoolVarP(&cors, "cors", "x", false, "Enable cors headers")
	WebStart.PersistentFlags().BoolVarP(&cluster, "cluster", "s", false, "cluster-alone mode or distributed mod")
}

func usage() {
	usageStr := `
███████╗███████╗███╗   ██╗██████╗ ███████╗ █████╗ 
╚══███╔╝██╔════╝████╗  ██║██╔══██╗██╔════╝██╔══██╗
  ███╔╝ █████╗  ██╔██╗ ██║██║  ██║█████╗  ███████║
 ███╔╝  ██╔══╝  ██║╚██╗██║██║  ██║██╔══╝  ██╔══██║
███████╗███████╗██║ ╚████║██████╔╝███████╗██║  ██║
╚══════╝╚══════╝╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝                                                 
`
	fmt.Printf("%s\n", usageStr)
}

func setup() {
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
}

func run() error {
	engine := gin.Default()
	router.Setup(engine, cors)
	return engine.Run(":" + viper.GetString("base.port"))
}