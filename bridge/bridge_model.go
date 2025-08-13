package bridge

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type BridgeModel struct {
	// Add fields as needed for the bridge model
	choices  []string // items on the to-do list
	cursor   int      // which to-do list item our cursor is pointing at
	choice   string   // which to-do list item is selected
	viewport viewport.Model
	ready    bool // 用于窗口大小就绪判断
	logs     []string
}

func (m BridgeModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func NewModel() BridgeModel {
	subs := []string{}
	for k := range subModels {
		subs = append(subs, k)
	}
	return BridgeModel{
		choices: subs,
		logs:    []string{},
	}
}

func RegisterSubModel(choice string, m tea.Model) {
	if _, ok := subModels[choice]; !ok {
		subModels[choice] = m
	}
}
