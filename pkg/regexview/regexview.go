package regexview

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	evenMatchStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("2")).
			Bold(true)
	oddMatchStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("3")).
			Bold(true)
)

type Model struct {
	expression *regexp.Regexp
	value      string
	width      int
	height     int
}

func New(width, height int) Model {
	return Model{
		width:  width,
		height: height,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m *Model) renderContainer(s string) string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Left,
		wordwrap.String(s, m.width),
	)
}

func (m Model) View() string {
	if m.expression == nil {
		return m.renderContainer(m.value)
	}

	var b strings.Builder
	lastIndex := 0

	matches := m.expression.FindAllStringIndex(m.value, -1)
	for i, match := range matches {
		s := &evenMatchStyle
		if i%2 == 1 {
			s = &oddMatchStyle
		}

		b.WriteString(m.value[lastIndex:match[0]])
		b.WriteString(s.Render(m.value[match[0]:match[1]]))
		lastIndex = match[1]
	}

	b.WriteString(m.value[lastIndex:])

	return m.renderContainer(b.String())
}

func (m *Model) SetExpression(expression *regexp.Regexp) {
	m.expression = expression
}

func (m *Model) SetExpressionString(expression string) error {
	expr, err := regexp.Compile(expression)

	m.SetExpression(expr)

	return err
}

func (m *Model) SetValue(value string) {
	m.value = value
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func (m *Model) SetSize(width, height int) {
	m.SetWidth(width)
	m.SetHeight(height)
}
