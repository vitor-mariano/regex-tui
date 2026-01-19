package regexview

import (
	"strings"

	"charm.land/lipgloss/v2"
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
	whitespaceStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("242"))
	evenMatchWhitespaceStyle = evenMatchStyle.
					Foreground(lipgloss.Color("172"))
	oddMatchWhitespaceStyle = oddMatchStyle.
				Foreground(lipgloss.Color("25"))
)

func whitespaceStyleFor(matchStyle *lipgloss.Style) lipgloss.Style {
	if matchStyle == &evenMatchStyle {
		return evenMatchWhitespaceStyle
	}
	if matchStyle == &oddMatchStyle {
		return oddMatchWhitespaceStyle
	}
	return whitespaceStyle
}

func renderChar(r rune, matchStyle *lipgloss.Style) string {
	switch r {
	case ' ':
		return whitespaceStyleFor(matchStyle).Render("·")
	case '\n':
		return whitespaceStyleFor(matchStyle).Render("↵") + "\n"
	default:
		if matchStyle != nil {
			return matchStyle.Render(string(r))
		}
		return string(r)
	}
}

type Model struct {
	expression    Regex
	baseExpStr    string
	value         string
	width, height int
	scrollYOffset int
	global        bool
	insensitive   bool
	regexp2       bool
}

func New(width, height int) *Model {
	return &Model{
		width:  width,
		height: height,
	}
}

func (m *Model) renderContainer(s string) string {
	wrapped := wordwrap.String(s, m.width)
	lines := strings.Split(wrapped, "\n")

	start := m.scrollYOffset
	if start > len(lines) {
		start = len(lines)
	}

	end := start + m.height
	if end > len(lines) {
		end = len(lines)
	}

	visibleLines := lines[start:end]

	for len(visibleLines) < m.height {
		visibleLines = append(visibleLines, "")
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Left,
		strings.Join(visibleLines, "\n"),
	)
}

func (m *Model) View() string {
	if m.expression == nil {
		var b strings.Builder
		for _, r := range m.value {
			b.WriteString(renderChar(r, nil))
		}
		return m.renderContainer(b.String())
	}

	var matches [][]int
	if m.global {
		matches = m.expression.FindAllStringIndex(m.value, -1)
	} else if match := m.expression.FindStringIndex(m.value); match != nil {
		matches = [][]int{match}
	}

	matchStyles := make(map[int]*lipgloss.Style)
	for i, match := range matches {
		style := &evenMatchStyle
		if i%2 == 1 {
			style = &oddMatchStyle
		}
		for pos := match[0]; pos < match[1]; pos++ {
			matchStyles[pos] = style
		}
	}

	var b strings.Builder
	for i, r := range m.value {
		style := matchStyles[i]
		b.WriteString(renderChar(r, style))
	}

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

func (m *Model) SetScrollYOffset(offset int) {
	m.scrollYOffset = offset
}

func (m *Model) Validate(expression string) error {
	_, err := m.newRegexp(expression)
	return err
}
