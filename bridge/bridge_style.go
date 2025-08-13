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

	viewporttitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	viewportinfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return viewporttitleStyle.BorderStyle(b)
	}()
)

var (
	subModels map[string]tea.Model = make(map[string]tea.Model)
	LogChan   chan string
)

type BridgeMsg struct {
	Cmd int
}
