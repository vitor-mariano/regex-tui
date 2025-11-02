package screen

import (
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/styles"
)

var title = lipgloss.NewStyle().
	Background(styles.PrimaryColor).
	Bold(true).
	Foreground(styles.LightColor).
	Padding(0, 1).
	MarginLeft(1).
	MarginTop(1).
	Render("Regex TUI")
