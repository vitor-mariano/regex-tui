package subject

import (
	"github.com/charmbracelet/bubbles/v2/textarea"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/styles"
	"github.com/vitor-mariano/regex-tui/pkg/components/regexview"
)

type Model struct {
	input         *textarea.Model
	view          *regexview.Model
	width, height int
}

func New(initialValue, initialExpression string) *Model {
	m := textarea.New()
	m.SetValue(initialValue)
	m.SetVirtualCursor(true)
	m.SetStyles(textarea.Styles{
		Cursor: textarea.CursorStyle{
			Color: styles.PrimaryColor,
			Blink: true,
		},
		Focused: textarea.StyleState{
			CursorLine: lipgloss.NewStyle().UnsetBackground(),
		},
	})
	m.Prompt = ""
	m.ShowLineNumbers = false

	sv := regexview.New(0, 0)
	sv.SetExpression(initialExpression)
	sv.SetValue(initialValue)

	return &Model{input: m, view: sv}
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.view.SetValue(m.input.Value())

	return cmd
}

func (m *Model) View() string {
	s := &styles.InputContainerStyle
	if m.input.Err != nil {
		s = &styles.ErrorInputContainerStyle
	} else if m.input.Focused() {
		s = &styles.FocusedInputContainerStyle
	}

	v := m.input.View()
	if !m.input.Focused() {
		v = m.view.View()
	}

	return s.Width(m.width).Render(v)
}

func (m *Model) SetSize(width, height int) {
	const subjectHSpacing = 4

	m.width = width
	m.height = height
	m.input.SetWidth(width - subjectHSpacing)
	m.input.SetHeight(height)

	m.view.SetSize(width-subjectHSpacing-1, height)
}

func (m *Model) GetInput() *textarea.Model {
	return m.input
}

func (m *Model) GetView() *regexview.Model {
	return m.view
}

func (m *Model) SetExpression(expression string) error {
	return m.view.SetExpression(expression)
}
