package bridge

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BridgeModel struct {
	// Add fields as needed for the bridge model
	choices []string // items on the to-do list
	cursor  int      // which to-do list item our cursor is pointing at
}

func (m BridgeModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *BridgeModel) AppendLog(s string) {
	if len(BridgeLogs) >= 200 {
		BridgeLogs = BridgeLogs[1:] // 保持日志长度不超过200
	}
	BridgeLogs = append(BridgeLogs, s)
}

func NewModel() BridgeModel {
	BridgeLogs = []string{} // 初始化日志
	subs := []string{}
	for k := range subModels {
		subs = append(subs, k)
	}
	return BridgeModel{
		choices: subs,
	}
}

func RegisterSubModel(choice string, m tea.Model) {
	if _, ok := subModels[choice]; !ok {
		subModels[choice] = m
	}
}
