package ui

import (
	"robotech/squadron/tcpclient"

	tea "github.com/charmbracelet/bubbletea"
)

func (m TcpClientModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Clear()
			return m, tea.Quit
		case "h":
			m.state = StateInputHost
		case "r":
			m.state = StateRunning
			if m.tcpClient != nil {
				if !m.tcpClient.IsConnected() {
					go m.tcpClient.Start()
				}
			} else {
				m.tcpClient = tcpclient.NewTCPClient(m.hostInput, nil)
			}
			go m.tcpClient.Start()
		case "e":
			m.state = StateEditMsg
		case "s":
			m.state = StateInit
			if m.tcpClient != nil {
				m.tcpClient.Close()
			}
		}
	}
	return m, nil
}
