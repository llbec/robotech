package main

import (
	"log"
	"os"

	"robotech/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logChan := make(chan string, 100)

	model := ui.NewModel(logChan)

	p := tea.NewProgram(model)

	// 启动一个协程，将 logChan 中的日志消息转发为 tea.Msg 传给 p.Send()
	go func() {
		for msg := range logChan {
			p.Send(msg)
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
