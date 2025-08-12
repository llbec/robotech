package bridge

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m BridgeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if _, ok := subModels[m.choice]; !ok {
			switch msg.String() {
			case "ctrl+c", "q":
				AppendLog("Bridge quit" + m.choice)
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
				cmds = append(cmds, subModels[m.choices[m.cursor]].Init())
			}
		}
		if _, ok := subModels[m.choice]; ok {
			subModels[m.choice], cmd = subModels[m.choice].Update(msg)
			cmds = append(cmds, cmd)
			if cmd != nil {
				m.choice = ""
				cmd = nil
				AppendLog("Back to Bridge")
			}
		}
	case string:
		AppendLog(msg)
	}
	return m, tea.Batch(cmds...)
}
