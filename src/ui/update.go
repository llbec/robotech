package ui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case StateMenu:
			cmd := m.handleMenuKeys(msg)
			return m, cmd
		case StateInputIP:
			m.handleInputIP(msg)
		case StateInputPort:
			m.handleInputPort(msg)
		case StateRunning:
			m.handleRunningKeys(msg)
		case StateSendDataModeSelection:
			m.handleSendDataModeSelectionKeys(msg)
		case StateSendDataManualInput:
			m.handleSendDataManualInputKeys(msg)
		case StateSendDataFileInput:
			m.handleSendDataFileInputKeys(msg)
		}
	case string:
		// 收到 TCP 模块日志
		m.AppendLog(msg)
	}
	return m, nil
}
