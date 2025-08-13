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
		case StateRunning:
			cmd = m.handleMenuKeys(msg)
		case StageSelectFile:
			cmd = m.handleSendFile(msg)
		case StageSendFile:
		}
	default:
		//bridge.LogChan <- fmt.Sprintf("Unhandled message type: %v", msg)
		m.filepicker, cmd = m.filepicker.Update(msg)
	}
	return m, cmd
}
