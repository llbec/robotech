package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	logStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC"))

	optionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))
)

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(" Robotech TCP Client/Server ") + "\n\n")

	switch m.state {
	case StateMenu:
		b.WriteString("选择模式:\n")
		for i, choice := range m.choices {
			cursor := "  "
			if m.cursor == i {
				cursor = cursorStyle.Render("> ")
				b.WriteString(cursor + choice + "\n")
			} else {
				b.WriteString("  " + choice + "\n")
			}
		}
	case StateInputIP:
		b.WriteString("请输入 IP 地址:\n")
		b.WriteString(m.ipInput + "_\n")
	case StateInputPort:
		b.WriteString("请输入端口号:\n")
		b.WriteString(m.portInput + "_\n")
	case StateRunning:
		b.WriteString("日志区:\n")
		logsLen := len(m.logs)
		start := 0
		if logsLen > 20 {
			start = logsLen - 20
		}
		for _, line := range m.logs[start:] {
			b.WriteString(logStyle.Render(line) + "\n")
		}
		b.WriteString("\n")

		b.WriteString(optionStyle.Render("操作选项:\n"))
		options := []string{"发送数据", "关闭连接"}

		for i, opt := range options {
			cursor := "  "
			if m.runningCursor == i {
				cursor = cursorStyle.Render("> ")
				b.WriteString(cursor + opt + "\n")
			} else {
				b.WriteString("  " + opt + "\n")
			}
		}
	case StateSendDataModeSelection:
		b.WriteString(optionStyle.Render("选择发送数据方式 (ESC 返回):\n"))
		b.WriteString("1) 手动输入\n")
		b.WriteString("2) 发送文件\n")
	case StateSendDataManualInput:
		b.WriteString("手动输入发送数据 (Enter 发送, ESC 取消):\n")
		b.WriteString(m.sendInput + "_\n")
	case StateSendDataFileInput:
		b.WriteString("输入文件路径 (Enter 发送全部行，ESC 取消):\n")
		b.WriteString(m.filePath + "_\n")
	}

	return b.String()
}
