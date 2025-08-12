package bridge

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BridgeModel struct {
	// Add fields as needed for the bridge model
	choices []string // items on the to-do list
	cursor  int      // which to-do list item our cursor is pointing at
	choice  string   // which to-do list item is selected
}

func (m BridgeModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
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
