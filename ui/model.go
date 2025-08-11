package ui

import (
	"robotech/tcp"

	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	StateMenu State = iota
	StateInputIP
	StateInputPort
	StateRunning
	StateSendDataModeSelection
	StateSendDataManualInput
	StateSendDataFileInput
)

type Mode int

const (
	ModeNone Mode = iota
	ModeTCPServer
	ModeTCPClient
)

type Model struct {
	state   State
	mode    Mode
	choices []string
	cursor  int

	ipInput   string
	portInput string

	logs []string

	logChan chan string

	tcpServer *tcp.TCPServer
	tcpClient *tcp.TCPClient

	// 发送数据相关
	sendInput string
	filePath  string

	runningCursor int // 运行时选项光标 0:发送数据 1:关闭连接
}

func NewModel(logChan chan string) Model {
	return Model{
		state:   StateMenu,
		mode:    ModeNone,
		choices: []string{"TCP Server", "TCP Client", "退出"},
		cursor:  0,
		logs:    []string{},
		logChan: logChan,
	}
}

func (m *Model) AppendLog(s string) {
	m.logs = append(m.logs, s)
}

func (m *Model) closeTCP() {
	if m.tcpServer != nil {
		m.tcpServer.Close()
		m.tcpServer = nil
	}
	if m.tcpClient != nil {
		m.tcpClient.Close()
		m.tcpClient = nil
	}
}

func (m Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
