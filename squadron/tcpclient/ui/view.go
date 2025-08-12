package ui

import (
	"robotech/bridge"
	"strings"
)

func (m TcpClientModel) View() string {
	var b strings.Builder
	b.WriteString(bridge.TitleStyle.Render(" Robotech ") + "\n\n")

	switch m.state {
	case StateInputHost:
		b.WriteString("请输入 IP:Port :\n")
		b.WriteString(m.hostInput + "_\n")
	case StateEditMsg:
	case StateRunning:
		b.WriteString(bridge.OptionStyle.Render("TCP Server " + m.tcpClient.Host() + "\n\n"))
	default:
	}

	b.WriteString("日志区:\n")

	for _, line := range bridge.BridgeLogs {
		b.WriteString(bridge.LogStyle.Render(line) + "\n")
	}
	b.WriteString("\n")
	b.WriteString(bridge.HelpStyle.Render("\nh: set Host • r: run • e:edit message and send • s:stop • q: exit\n"))

	return b.String()
}
