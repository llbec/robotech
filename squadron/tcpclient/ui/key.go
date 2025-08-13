package ui

import (
	"os"
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
		return func() tea.Msg { return bridge.BridgeMsg{Cmd: 1} } // 返回到 Bridge
	case "s":
		m.state = StageSelectFile
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
			return func() tea.Msg { return bridge.BridgeMsg{Cmd: 1} } // 返回到 Bridge
		}
		m.state = StateRunning
		m.client = tcpclient.NewTCPClient(m.hostInput, bridge.LogChan)
		go m.client.Start()
	case tea.KeyCtrlC:
		m.Clear()
		return func() tea.Msg { return bridge.BridgeMsg{Cmd: 1} } // 返回到 Bridge
	}
	return nil
}

func (m *TcpClientModel) handleSendFile(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q":
		m.state = StateRunning
	case "s":
		// load file and send
		data, err := os.ReadFile(m.selectedFile)
		if err != nil {
			bridge.LogChan <- "file: " + m.selectedFile + " read error: " + err.Error()
			m.state = StateRunning
			return nil
		}
		m.client.SendHexString(string(data))
		bridge.LogChan <- m.selectedFile + " sent: " + string(data)
		m.state = StateRunning
	default:
		//bridge.LogChan <- fmt.Sprintf("File Picker Path(update): %v", m.filepicker.CurrentDirectory)
		m.filepicker, cmd = m.filepicker.Update(msg)
		if did, path := m.filepicker.DidSelectFile(msg); did {
			m.selectedFile = path
			//m.state = StageSendFile
		}
		if didSelect, _ := m.filepicker.DidSelectDisabledFile(msg); didSelect {
			// Let's clear the selectedFile and display an error.
			m.selectedFile = ""
		}
	}

	return cmd
}
