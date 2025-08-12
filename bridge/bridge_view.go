package bridge

import (
	"strings"
)

func (m BridgeModel) View() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(" Robotech ") + "\n\n")
	b.WriteString("选择模式:\n")
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = CursorStyle.Render("> ")
			b.WriteString(cursor + choice + "\n")
		} else {
			b.WriteString("  " + choice + "\n")
		}
	}
	b.WriteString("日志区:\n")
	logsLen := len(BridgeLogs)
	start := 0
	if logsLen > 20 {
		start = logsLen - 20
	}
	for _, line := range BridgeLogs[start:] {
		b.WriteString(LogStyle.Render(line) + "\n")
	}
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("\nq: exit\n"))
	return b.String()
}
