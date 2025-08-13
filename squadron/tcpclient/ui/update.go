package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m TcpClientModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case StateInit:
			cmd = m.handleInitKeys(msg)
		case StateInputHost:
			cmd = m.handleInputHost(msg)
		case StateEditMsg:
			cmd = m.handleMenuKeys(msg)
		case StateRunning:
			cmd = m.handleMenuKeys(msg)
		}
	}
	//bridge.AppendLog(fmt.Sprintf("TCP Client State: %v, Host Input: %s", m.state, m.hostInput))
	return m, cmd
}
