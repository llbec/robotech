package ui

import (
	"robotech/bridge"
	"robotech/squadron/tcpclient"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *TcpClientModel) handleInitKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEnter:
		m.state = StateInputHost
	}
	return nil
}

func (m *TcpClientModel) handleMenuKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c", "q":
		m.Clear()
		return tea.Quit
	case "s":
		m.state = StateEditMsg
	}
	return nil
}

func (m *TcpClientModel) handleInputHost(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyRunes:
		m.hostInput += string(msg.Runes)
	case tea.KeyBackspace:
		if len(m.hostInput) > 0 {
			m.hostInput = m.hostInput[:len(m.hostInput)-1]
		}
	case tea.KeyEnter:
		if !tcpclient.IsValidHost(m.hostInput) {
			m.Clear()
			return tea.Quit
		}
		m.state = StateRunning
		m.client = tcpclient.NewTCPClient(m.hostInput, bridge.LogChan)
		go m.client.Start()
	case tea.KeyCtrlC:
		m.Clear()
		return tea.Quit
	}
	return nil
}
