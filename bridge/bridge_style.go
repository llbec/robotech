package bridge

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	LogStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC"))

	OptionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))
	HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

var (
	BridgeLogs                      = []string{}
	subModels  map[string]tea.Model = make(map[string]tea.Model)
)

func AppendLog(s string) {
	if len(BridgeLogs) >= 200 {
		BridgeLogs = BridgeLogs[1:] // 保持日志长度不超过200
	}
	BridgeLogs = append(BridgeLogs, s)
}
