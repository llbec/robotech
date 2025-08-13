package main

import (
	"log"
	"os"

	"robotech/bridge"
	_ "robotech/squadron/tcpclient/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := bridge.NewModel()

	p := tea.NewProgram(model)

	bridge.LogChan = make(chan string, 1024*8)
	go func() {
		for logMsg := range bridge.LogChan {
			p.Send(logMsg)
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
