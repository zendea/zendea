package main

import (
	"os"
	"runtime"

	"github.com/urfave/cli"
	"zendea/config"
	"zendea/cmd"
)

const APP_VER = "0.0.3-dev"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config.AppName = "Zendea"
}

func main() {
	app := cli.NewApp()
	app.Name = "Zendea"
	app.Usage = "A free, open-source, self-hosted forum software written in Go."
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdWeb,
	}
	
	// default configuration flags
	defaultFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "./app.yaml",
			Usage: "Custom configuration file path",
		},
		cli.VersionFlag,
	}

	// Set the default to be equivalent to cmdWeb and add the default flags
	app.Flags = append(app.Flags, cmd.CmdWeb.Flags...)
	app.Flags = append(app.Flags, defaultFlags...)
	app.Action = cmd.CmdWeb.Action

	app.Run(os.Args)
}

