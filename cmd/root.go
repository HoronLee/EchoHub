package cmd

import (
	"os"

	"github.com/HoronLee/EchoHub/internal/cli"
	"github.com/HoronLee/EchoHub/internal/config"
	commonModel "github.com/HoronLee/EchoHub/internal/model/common"
	"github.com/spf13/cobra"
)

// 初始化版本信息到 cli 包
func init() {
	cli.Version = commonModel.Version
	cli.BuildTime = commonModel.BuildTime
	cli.GitCommit = commonModel.GitCommit
}

var configPath string

// rootCmd 是 EchoHub 的根命令
// 默认启动CLI With TUI
var rootCmd = &cobra.Command{
	Use:   "echohub",
	Short: "基于Echo、Gorm、Viper、Wire、Cobra的HTTP快速开发框架",
	Long:  `基于Echo、Gorm、Viper、Wire、Cobra的HTTP快速开发框架`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.LoadAppConfig(configPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoTui()
	},
}

// serveCmd 是启动 EchoHub 服务的命令
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 EchoHub HTTP 服务",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoServeWithBlock()
	},
}

// tuiCmd 是启动 EchoHub TUI 的命令
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "启动 EchoHub TUI",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoTui()
	},
}

// versionCmd 是查看当前版本信息的命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看当前版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoVersion()
	},
}

// infoCmd 是查看当前信息的命令
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "查看当前信息",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoEchoHubInfo()
	},
}

// helloCmd 是输出 EchoHub Logo 的命令
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "输出 EchoHub Logo",
	Run: func(cmd *cobra.Command, args []string) {
		cli.DoHello()
	},
}

// Execute 是根命令的入口函数
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
