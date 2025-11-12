package regexview

import (
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/muesli/reflow/wordwrap"
	. "github.com/vitor-mariano/regex-tui/pkg/regex"
	"github.com/vitor-mariano/regex-tui/pkg/regex/re2"
	"github.com/vitor-mariano/regex-tui/pkg/regex/regexp2"
)

var (
	evenMatchStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("220")).
			Foreground(lipgloss.Color("232")).
			Bold(true)
	oddMatchStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("117")).
			Foreground(lipgloss.Color("232")).
			Bold(true)
)

type Model struct {
	expression    Regex
	baseExpStr    string
	global        bool
	insensitive   bool
	regexp2       bool
	value         string
	width, height int
}

func New(width, height int) *Model {
	return &Model{
		width:  width,
		height: height,
	}
}

func (m *Model) renderContainer(s string) string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Left,
		wordwrap.String(s, m.width),
	)
}

func (m *Model) View() string {
	if m.expression == nil {
		return m.renderContainer(m.value)
	}

	var b strings.Builder
	lastIndex := 0

	var matches [][]int
	if m.global {
		matches = m.expression.FindAllStringIndex(m.value, -1)
	} else {
		matches = [][]int{m.expression.FindStringIndex(m.value)}
	}
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

func (m *Model) newRegexp(expression string) (Regex, error) {
	if m.regexp2 {
		return regexp2.New(expression)
	}

	return re2.New(expression)
}

func (m *Model) setRegexp(expression string) error {
	prefix := ""
	if m.insensitive {
		prefix = "(?i)"
	}

	regex, err := m.newRegexp(prefix + expression)
	if err != nil {
		return err
	}

	m.expression = regex
	return nil
}

func (m *Model) SetExpression(expression string) error {
	err := m.setRegexp(expression)
	if err == nil {
		m.baseExpStr = expression
	}

	return err
}

func (m *Model) SetGlobal(global bool) {
	m.global = global
}

func (m *Model) SetInsensitive(insensitive bool) {
	m.insensitive = insensitive
	m.setRegexp(m.baseExpStr)
}

func (m *Model) SetRegexp2(regexp2 bool) error {
	m.regexp2 = regexp2
	return m.setRegexp(m.baseExpStr)
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

func (m *Model) Validate(expression string) error {
	_, err := m.newRegexp(expression)
	return err
}
