package ui

import (
	"robotech/bridge"
	"robotech/squadron/tcpclient"

	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	StateInit State = iota
	StateInputHost
	StateEditMsg
	StateRunning
	StateSendDataFileInput
)

type TcpClientModel struct {
	// Add fields as needed for the model
	state     State
	client    *tcpclient.TCPClient
	hostInput string
}

func init() {
	bridge.RegisterSubModel("TCP Client", TcpClientModel{
		state:     StateInit,
		hostInput: "",
		client:    nil,
	})
}

func (m TcpClientModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *TcpClientModel) Clear() {
	m.hostInput = ""
	m.client.Close()
	m.state = StateInit
}
