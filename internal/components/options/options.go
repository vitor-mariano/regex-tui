package options

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/styles"
	"github.com/vitor-mariano/regex-tui/pkg/components/multiselect"
)

const (
	GlobalOption      = "Global"
	InsensitiveOption = "Insensitive"
	Regexp2Option     = "Regexp2"
)

type Model struct {
	options             *multiselect.Model
	isOptionsDialogOpen bool
}

func New() *Model {
	return &Model{
		options:             multiselect.New([]string{GlobalOption, InsensitiveOption, Regexp2Option}),
		isOptionsDialogOpen: false,
	}
}

func (m *Model) IsOpen() bool {
	return m.isOptionsDialogOpen
}

func (m *Model) Open() {
	m.isOptionsDialogOpen = true
}

func (m *Model) OnToggle(onToggle func(item string, selected bool)) {
	m.options.OnToggle(onToggle)
	m.options.SetSelected(GlobalOption)
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if key.Matches(msg, keys.Exit) {
			m.isOptionsDialogOpen = false
		}
	}

	return m.options.Update(msg)
}

func (m *Model) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.MutedColor).
		Padding(1, 4, 1, 2).
		Render(m.options.View())
}
