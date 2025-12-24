package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HoronLee/EchoHub/internal/config"
	"github.com/HoronLee/EchoHub/internal/di"
	commonModel "github.com/HoronLee/EchoHub/internal/model/common"
	"github.com/HoronLee/EchoHub/internal/server"
	"github.com/HoronLee/EchoHub/internal/tui"
	"github.com/charmbracelet/huh"
)

var s *server.HTTPServer // s 是全局的 EchoHub 服务器实例

// DoServe 启动服务
func DoServe() {
	// 通过Wire初始化服务器
	srv, cleanup, err := di.InitServer(&config.Config)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer cleanup()
	s = srv

	// 启动服务器
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// DoServeWithBlock 阻塞当前线程，直到服务器停止
func DoServeWithBlock() {
	// 通过Wire初始化服务器
	srv, cleanup, err := di.InitServer(&config.Config)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer cleanup()
	s = srv

	// 启动服务器
	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// 阻塞主线程，直到接收到终止信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 创建 context，最大等待 5 秒优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		tui.PrintCLIInfo("❌ 服务停止", "服务器强制关闭")
		os.Exit(1)
	}
	tui.PrintCLIInfo("🎉 停止服务成功", "EchoHub 服务器已停止")
}

// DoStopServe 停止服务
func DoStopServe() {
	if s == nil {
		tui.PrintCLIInfo("⚠️ 停止服务", "EchoHub 服务器未启动")
		return
	}

	// 创建 context，最大等待 5 秒优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil {
		tui.PrintCLIInfo("😭 停止服务失败", err.Error())
		return
	}

	s = nil // 清空全局服务器实例

	tui.PrintCLIInfo("🎉 停止服务成功", "EchoHub 服务器已停止")
}

// DoTui 执行 TUI
func DoTui() {
	// 清除屏幕当前字符
	tui.ClearScreen()
	// 打印 ASCII 风格 Banner
	tui.PrintCLIBanner()

	for {
		// 换行
		fmt.Println()

		var action string
		var options []huh.Option[string]

		if s == nil {
			options = append(options, huh.NewOption("🚀 启动 Web 服务", "serve"))
		} else {
			options = append(options, huh.NewOption("🛑 停止 Web 服务", "stopserve"))
		}

		options = append(options,
			huh.NewOption("🦖 查看信息", "info"),
			huh.NewOption("📌 查看版本", "version"),
			huh.NewOption("❌ 退出", "exit"),
		)

		err := huh.NewSelect[string]().
			Title("欢迎使用 EchoHub TUI .").
			Options(options...).
			Value(&action).
			WithTheme(huh.ThemeCatppuccin()).
			Run()
		if err != nil {
			log.Fatal(err)
		}

		switch action {
		case "serve":
			tui.ClearScreen()
			DoServe()
		case "stopserve":
			tui.ClearScreen()
			DoStopServe()
		case "info":
			tui.ClearScreen()
			DoEchoHubInfo()
		case "version":
			tui.ClearScreen()
			DoVersion()
		case "exit":
			fmt.Println("👋 感谢使用 EchoHub TUI，期待下次再见")
			return
		}
	}
}

// DoVersion 打印版本信息
func DoVersion() {
	item := struct{ Title, Msg string }{
		Title: "📦 当前版本",
		Msg:   "v" + commonModel.Version,
	}
	tui.PrintCLIWithBox(item)
}

// DoEchoHubInfo 打印 EchoHub 信息
func DoEchoHubInfo() {
	if _, err := fmt.Fprintln(os.Stdout, tui.GetEchoHubInfo()); err != nil {
		fmt.Fprintf(os.Stderr, "failed to print EchoHub info: %v\n", err)
	}
}

// DoHello 打印 EchoHub Logo
func DoHello() {
	tui.ClearScreen()
	tui.PrintCLIBanner()
}
