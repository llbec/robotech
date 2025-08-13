package ui

import (
	"os"
	"robotech/bridge"
	"robotech/squadron/tcpclient"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	StateInit State = iota
	StateInputHost
	StateRunning
	StageSelectFile
	StageSendFile
)

type TcpClientModel struct {
	// Add fields as needed for the model
	state        State
	client       *tcpclient.TCPClient
	hostInput    string
	filepicker   filepicker.Model
	selectedFile string
}

func init() {
	bridge.RegisterSubModel("TCP Client", NewTCPClientModel())
}

func NewTCPClientModel() TcpClientModel {
	fp := filepicker.New()
	fp.AllowedTypes = []string{".msg"} // 自定义过滤
	fp.CurrentDirectory, _ = os.Getwd()
	fp.ShowHidden = true // 显示隐藏文件
	fp.SetHeight(5)      // 设置文件选择器的高度
	return TcpClientModel{
		state:      StateInit,
		filepicker: fp,
		hostInput:  "",
		client:     nil,
	}
}

func (m TcpClientModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.filepicker.Init()
}

func (m *TcpClientModel) Clear() {
	m.hostInput = ""
	if m.client != nil {
		m.client.Close()
	}
	m.state = StateInit
}
