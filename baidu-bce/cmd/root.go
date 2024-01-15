package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Used for flags
var configFile string

// The base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "baidu-bce",
	Short: "百度智能云",
	Long: `百度智能云的命令行界面
	
使用SSL证书、DNS（“域名服务 BCD”）、CDN。

等效网页：https://console.bce.baidu.com`,
}

// Add all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file containing access_key and secret_key (default is ./baidu-bce.yaml)")
}

func initConfig() {
	if configFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("baidu-bce")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fail to read config file: %w", err))
	}
}
