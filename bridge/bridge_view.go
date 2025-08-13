package bridge

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m BridgeModel) View() string {
	var b strings.Builder
	//b.WriteString(TitleStyle.Render(" Robotech ") + "\n\n")
	//b.WriteString("日志区:\n")
	/*logsLen := len(BridgeLogs)
	start := 0
	if logsLen > 20 {
		start = logsLen - 20
	}
	for _, line := range BridgeLogs {
		b.WriteString(LogStyle.Render(line) + "\n")
	}*/
	b.WriteString(fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView()))
	if _, ok := subModels[m.choice]; ok {
		b.WriteString(subModels[m.choice].View())
	} else {
		b.WriteString("\n选择模式:\n")
		for i, choice := range m.choices {
			cursor := "  "
			if m.cursor == i {
				cursor = CursorStyle.Render("> ")
				b.WriteString(cursor + choice + "\n")
			} else {
				b.WriteString("  " + choice + "\n")
			}
		}
		b.WriteString(HelpStyle.Render("\nq: exit\n"))
	}

	return b.String()
}

func (m BridgeModel) headerView() string {
	title := viewporttitleStyle.Render("Bridge Paper")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m BridgeModel) footerView() string {
	info := viewportinfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
