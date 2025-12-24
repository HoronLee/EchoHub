package cmd

import (
	"github.com/spf13/cobra"
)

// init 函数用于初始化根命令和子命令
func init() {
	// 解决Windows下使用 Cobra 触发 mousetrap 提示
	cobra.MousetrapHelpText = ""

	// 添加全局flag
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "配置文件路径")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(helloCmd)
}
