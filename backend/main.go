package main

import (
	"github.com/spf13/cobra"
	"os"
	"zendea/config"
	"zendea/cmd"
)

var rootCmd = &cobra.Command{
	Use:               "zendea",
	Short:             "zendea API server",
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	Long:              `Start zendea API server`,
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
}

func init() {
	config.AppName = "Zendea"
	rootCmd.AddCommand(cmd.WebStart)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
