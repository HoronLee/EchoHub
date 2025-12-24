package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"

	commonModel "github.com/HoronLee/EchoHub/internal/model/common"
)

var (
	// 信息样式（每行）
	infoStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.AdaptiveColor{
			Light: "236", Dark: "252",
		})

	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{
			Light: "#4338ca", Dark: "#FF7F7F",
		})

	// 高亮样式
	highlight = lipgloss.NewStyle().
			Bold(false).
			Italic(true).
			Foreground(lipgloss.AdaptiveColor{
			Light: "#7c3aed", Dark: "#53b7f5ff",
		})

	// 外框
	boxStyle = lipgloss.NewStyle().
			Bold(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#fb5151ff")).
			Padding(1, 1).
			Margin(1, 1)
)

const (
	banner = `
 ooooooooooo           oooo                  ooooo ooooo            oooo
 888    88   ooooooo   888ooooo    ooooooo   888   888 oooo  oooo   888ooooo
 888ooo8   888     888 888   888 888     888 888ooo888  888   888   888    888
 888    oo 888         888   888 888     888 888   888  888   888   888    888
o888ooo8888  88ooo888 o888o o888o  88ooo88  o888o o888o  888o88 8o o888ooo88   `
)

// GetLogoBanner 获取Logo横幅
func GetLogoBanner() string {
	lines := strings.Split(banner, "\n")
	var rendered []string

	colors := []string{
		"#FF7F7F", // 珊瑚红
		"#FFB347", // 桃橙色
		"#FFEB9C", // 金黄色
		"#B8E6B8", // 薄荷绿
		"#87CEEB", // 天空蓝
		"#DDA0DD", // 梅花紫
		"#F0E68C", // 卡其色
	}

	for i, line := range lines {
		color := lipgloss.Color(colors[i%len(colors)])
		style := lipgloss.NewStyle().Foreground(color)
		rendered = append(rendered, style.Render(line))
	}
	gradientBanner := lipgloss.JoinVertical(lipgloss.Left, rendered...)

	full := lipgloss.JoinVertical(lipgloss.Left,
		gradientBanner,
	)

	return full
}

// PrintCLIBanner 打印CLI横幅
func PrintCLIBanner() {
	banner := GetLogoBanner()

	if _, err := fmt.Fprintln(os.Stdout, banner); err != nil {
		fmt.Fprintf(os.Stderr, "failed to print banner: %v\n", err)
	}
}

// PrintCLIInfo 打印CLI信息
func PrintCLIInfo(title, msg string) {
	// 使用 lipgloss 渲染 CLI 信息
	fmt.Fprintln(os.Stdout, infoStyle.Render(titleStyle.Render(title)+": "+highlight.Render(msg)))
}

// CLIInfoItem 定义了一个CLI信息项，包含标题和消息
type CLIInfoItem struct {
	Title string
	Msg   string
}

// GetCLIPrintWithBox 获取带边框的CLI信息打印内容
func GetCLIPrintWithBox(items ...CLIInfoItem) string {
	if len(items) == 0 {
		return ""
	}

	var content string
	for i, item := range items {
		line := infoStyle.Render(titleStyle.Render(item.Title) + ": " + highlight.Render(item.Msg))
		if i > 0 {
			content += "\n"
		}
		content += line
	}

	boxedContent := boxStyle.Render(content)
	return boxedContent
}

// PrintCLIWithBox 打印带边框的CLI信息
func PrintCLIWithBox(items ...CLIInfoItem) {
	fmt.Fprintln(os.Stdout, GetCLIPrintWithBox(items...))
}

// ClearScreen 清屏函数，根据操作系统执行不同的清屏命令
func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls") // Windows 清屏命令
	} else {
		cmd = exec.Command("clear") // Linux/macOS 清屏命令
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// GetEchoHubInfo 获取EchoHub信息
func GetEchoHubInfo() string {
	content := lipgloss.JoinVertical(lipgloss.Left,
		infoStyle.Render("📦 "+titleStyle.Render("Version")+": "+highlight.Render(commonModel.Version)),
		infoStyle.Render("🧙 "+titleStyle.Render("Author")+": "+highlight.Render("HoronLee")),
		infoStyle.Render("👉 "+titleStyle.Render("Website")+": "+highlight.Render("https://horonlee.com/")),
		infoStyle.Render("👉 "+titleStyle.Render("GitHub")+": "+highlight.Render("https://github.com/HoronLee/EchoHub")),
	)

	full := lipgloss.JoinVertical(lipgloss.Left,
		boxStyle.Render(content),
	)

	return full
}
