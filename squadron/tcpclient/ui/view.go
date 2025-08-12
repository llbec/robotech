package ui

import (
	"robotech/bridge"
	"strings"
)

func (m TcpClientModel) View() string {
	var b strings.Builder

	b.WriteString("\n------------------------------------------\n")

	switch m.state {
	case StateInputHost:
		b.WriteString("\n请输入 IP:Port :\n")
		b.WriteString(m.hostInput + "_\n")
	case StateEditMsg:
	case StateRunning:
		b.WriteString(bridge.OptionStyle.Render("\nTCP Server " + m.client.Host() + "\n"))
		b.WriteString(bridge.HelpStyle.Render("\ne:edit message and send • s:stop • q: exit\n"))
	default:
	}

	return b.String()
}
