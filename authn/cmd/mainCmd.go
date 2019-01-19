package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// InitCmd represents the base command when called without any subcommands
	MainCmd = &cobra.Command{
		Use:   "mainCmd",
		Short: "",
		Long:  "",
	}
	confPath string
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := MainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	MainCmd.PersistentFlags().StringVarP(&confPath, "conf", "c", "", "path to config")
}
