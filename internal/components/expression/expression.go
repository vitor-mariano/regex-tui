package expression

import (
	"github.com/charmbracelet/bubbles/v2/textinput"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/vitor-mariano/regex-tui/internal/styles"
	"github.com/vitor-mariano/regex-tui/pkg/components/regexview"
)

type Model struct {
	input textinput.Model
	view  *regexview.Model
	width int
}

func New(initialValue string, view *regexview.Model) *Model {
	m := textinput.New()
	m.SetValue(initialValue)
	m.SetVirtualCursor(true)
	m.SetStyles(textinput.Styles{
		Cursor: textinput.CursorStyle{
			Color: styles.PrimaryColor,
			Blink: true,
		},
	})
	m.Prompt = ""
	m.Placeholder = "Expression"

	model := &Model{input: m, view: view}
	model.input.Validate = func(expr string) error {
		return model.view.Validate(expr)
	}

	return model
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return cmd
}

func (m *Model) View() string {
	s := &styles.InputContainerStyle
	if m.input.Err != nil {
		s = &styles.ErrorInputContainerStyle
	} else if m.input.Focused() {
		s = &styles.FocusedInputContainerStyle
	}

	return s.Width(m.width).Render(m.input.View())
}

func (m *Model) SetWidth(width int) {
	m.width = width
	m.input.SetWidth(width)
}

func (m *Model) GetInput() *textinput.Model {
	return &m.input
}
