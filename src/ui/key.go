package ui

import (
	"fmt"
	"robotech/src/tcp"
	"robotech/src/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleMenuKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter":
		if m.cursor == len(m.choices)-1 {
			return tea.Quit // 退出
		}
		m.mode = ModeTCPServer
		if m.cursor == 1 {
			m.mode = ModeTCPClient
		}
		m.state = StateInputIP
		m.ipInput = ""
	}
	return nil
}

func (m *Model) handleInputIP(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyRunes:
		m.ipInput += string(msg.Runes)
	case tea.KeyBackspace:
		if len(m.ipInput) > 0 {
			m.ipInput = m.ipInput[:len(m.ipInput)-1]
		}
	case tea.KeyEnter:
		if m.ipInput == "" {
			break
		}
		m.state = StateInputPort
		m.portInput = ""
	case tea.KeyCtrlC:
		m.state = StateMenu
	}
}

func (m *Model) handleInputPort(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyRunes:
		m.portInput += string(msg.Runes)
	case tea.KeyBackspace:
		if len(m.portInput) > 0 {
			m.portInput = m.portInput[:len(m.portInput)-1]
		}
	case tea.KeyEnter:
		if m.portInput == "" {
			break
		}
		m.state = StateRunning
		m.AppendLog(fmt.Sprintf("启动 %s 于 %s:%s", m.modeString(), m.ipInput, m.portInput))
		addr := m.ipInput + ":" + m.portInput
		switch m.mode {
		case ModeTCPServer:
			server := tcp.NewTCPServer(addr, m.logChan)
			m.tcpServer = server
			go server.Start()
		case ModeTCPClient:
			client := tcp.NewTCPClient(addr, m.logChan)
			m.tcpClient = client
			go client.Start()
		}
	case tea.KeyCtrlC:
		m.state = StateMenu
	}
}

func (m *Model) modeString() string {
	switch m.mode {
	case ModeTCPServer:
		return "TCP Server"
	case ModeTCPClient:
		return "TCP Client"
	}
	return "未知模式"
}

func (m *Model) handleRunningKeys(msg tea.KeyMsg) {
	switch msg.String() {
	case "up", "k":
		if m.runningCursor > 0 {
			m.runningCursor--
		}
	case "down", "j":
		if m.runningCursor < 1 {
			m.runningCursor++
		}
	case "enter":
		switch m.runningCursor {
		case 0:
			m.state = StateSendDataModeSelection
		case 1:
			m.closeTCP()
			m.state = StateMenu
			m.mode = ModeNone
			m.cursor = 0
			m.logs = []string{}
			m.runningCursor = 0
		}
	}
}

func (m *Model) handleSendDataModeSelectionKeys(msg tea.KeyMsg) {
	switch msg.String() {
	case "1":
		m.state = StateSendDataManualInput
		m.sendInput = ""
	case "2":
		m.state = StateSendDataFileInput
		m.filePath = ""
	case "esc":
		m.state = StateRunning
	}
}

func (m *Model) handleSendDataManualInputKeys(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyRunes:
		m.sendInput += string(msg.Runes)
	case tea.KeyBackspace:
		if len(m.sendInput) > 0 {
			m.sendInput = m.sendInput[:len(m.sendInput)-1]
		}
	case tea.KeyEnter:
		if m.sendInput != "" {
			if m.mode == ModeTCPServer && m.tcpServer != nil {
				m.tcpServer.Send(m.sendInput)
				m.AppendLog("发送: " + m.sendInput)
			} else if m.mode == ModeTCPClient && m.tcpClient != nil {
				m.tcpClient.Send(m.sendInput)
				m.AppendLog("发送: " + m.sendInput)
			}
		}
		m.sendInput = ""
		m.state = StateRunning
	case tea.KeyEsc:
		m.state = StateRunning
	}
}

func (m *Model) handleSendDataFileInputKeys(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyRunes:
		m.filePath += string(msg.Runes)
	case tea.KeyBackspace:
		if len(m.filePath) > 0 {
			m.filePath = m.filePath[:len(m.filePath)-1]
		}
	case tea.KeyEnter:
		if m.filePath != "" {
			// 调用文件读取辅助函数，逐行发送
			lines, err := utils.ReadFileLines(m.filePath)
			if err != nil {
				m.AppendLog(fmt.Sprintf("读取文件失败: %v", err))
			} else {
				for _, line := range lines {
					if m.mode == ModeTCPServer && m.tcpServer != nil {
						_ = m.tcpServer.Send(line)
						m.AppendLog("发送: " + line)
					} else if m.mode == ModeTCPClient && m.tcpClient != nil {
						_ = m.tcpClient.Send(line)
						m.AppendLog("发送: " + line)
					}
				}
			}
		}
		m.filePath = ""
		m.state = StateRunning
	case tea.KeyEsc:
		m.state = StateRunning
	}
}
