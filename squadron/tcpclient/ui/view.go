package ui

import (
	"robotech/bridge"
	"strings"
)

func (m TcpClientModel) View() string {
	var b strings.Builder
	switch m.state {
	case StateInputHost:
		b.WriteString("\n请输入 IP:Port :\n")
		b.WriteString(m.hostInput + "_\n")
	case StateEditMsg:
	case StateRunning:
		b.WriteString(bridge.OptionStyle.Render("\nTCP Server " + m.client.Host() + "\n"))
		b.WriteString(bridge.HelpStyle.Render("\ns:edit message and send • q: exit\n"))
	default:
	}

	return b.String()
}
