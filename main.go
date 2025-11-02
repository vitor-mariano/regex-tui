package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vitor-mariano/regex-tui/internal/screen"
)

func main() {
	if _, err := tea.NewProgram(screen.New()).Run(); err != nil {
		log.Fatalf("failed to start program: %v\n", err)
	}
}
