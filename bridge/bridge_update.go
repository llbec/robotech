package bridge

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m BridgeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-20) // 留出命令区空间
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 20
		}
	case tea.KeyMsg:
		if _, ok := subModels[m.choice]; !ok {
			switch msg.String() {
			case "ctrl+c", "q":
				//m.AppendLog("Bridge quit" + m.choice)
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				m.choice = m.choices[m.cursor]
				cmds = append(cmds, subModels[m.choice].Init())
			}
		} // 不用else，因为enter可能会触发子模型的更新
		if _, ok := subModels[m.choice]; ok {
			subModels[m.choice], cmd = subModels[m.choice].Update(msg)
			if cmd == nil {
				switch msg.String() {
				case "up":
					m.viewport.ScrollUp(1)
				case "down":
					m.viewport.ScrollDown(1)
				case "pgup":
					m.viewport.PageUp()
				case "pgdown":
					m.viewport.PageDown()
				}
			}
			cmds = append(cmds, cmd)
		}
	case string:
		m.AppendLog(msg)
	case BridgeMsg:
		switch msg.Cmd {
		case 1:
			m.choice = ""
			//m.AppendLog("Back to Bridge")
		}
	default:
		//m.AppendLog(fmt.Sprintf("Unhandled message type: %v, choice:%v", msg, m.choice))
		if _, ok := subModels[m.choice]; ok {
			subModels[m.choice], cmd = subModels[m.choice].Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *BridgeModel) AppendLog(s string) {
	if len(m.logs) >= 200 {
		m.logs = m.logs[1:] // 保持日志长度不超过200
	}
	m.logs = append(m.logs, s)
	m.viewport.SetContent(strings.Join(m.logs, "\n"))
}
