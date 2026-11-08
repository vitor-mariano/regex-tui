package regexview

import (
	"strings"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
	"github.com/muesli/reflow/wordwrap"
	. "github.com/vitor-mariano/regex-tui/pkg/regex"
	"github.com/vitor-mariano/regex-tui/pkg/regex/re2"
	"github.com/vitor-mariano/regex-tui/pkg/regex/regexp2"
)

const (
	spaceGlyph   = '·'
	newlineGlyph = "↵"

	ansiESC      = '\x1b' // ESC byte that starts an ANSI escape sequence
	ansiCSI      = '['    // CSI introducer following ESC
	ansiEnd      = 'm'    // SGR terminator
	ansiReset    = "\x1b[0m"
	ansiResetAlt = "\x1b[m"
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
	whitespaceGlyphStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))
)

type Model struct {
	expression      Regex
	baseExpStr      string
	global          bool
	insensitive     bool
	regexp2         bool
	showWhitespaces bool
	value           string
	width, height   int
	scrollYOffset   int
}

func New(width, height int) *Model {
	return &Model{
		width:  width,
		height: height,
	}
}

// substituteWhitespace replaces space characters with visible glyphs in an
// ANSI-styled string. It tracks whether each character is inside an ANSI-styled
// region (i.e. a regex match) so that:
//   - Spaces inside matches: glyph character replaces the original, inheriting
//     the surrounding match style.
//   - Spaces outside matches: glyph is wrapped in the dimmed style.
func substituteWhitespace(s string) string {
	var b strings.Builder
	styled := false
	i := 0
	for i < len(s) {
		// Detect ANSI CSI escape sequence: ESC [
		if s[i] == byte(ansiESC) && i+1 < len(s) && s[i+1] == byte(ansiCSI) {
			j := i + 2
			for j < len(s) && s[j] != byte(ansiEnd) {
				j++
			}
			if j < len(s) {
				seq := s[i : j+1]
				b.WriteString(seq)
				styled = seq != ansiReset && seq != ansiResetAlt
				i = j + 1
				continue
			}
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		switch r {
		case ' ':
			if styled {
				b.WriteRune(spaceGlyph)
			} else {
				b.WriteString(whitespaceGlyphStyle.Render(string(spaceGlyph)))
			}
		default:
			b.WriteString(s[i : i+size])
		}
		i += size
	}

	// Close any unclosed ANSI style to prevent it from leaking into
	// subsequent padding or adjacent content.
	if styled {
		b.WriteString(ansiReset)
	}

	return b.String()
}

// wrapLinePreservingSpaces wraps a single line (no newlines) at word boundaries
// while keeping space characters at the end of each wrapped line instead of
// consuming them. This is critical for whitespace visualization — standard
// wordwrap eats spaces at break points, making them invisible.
// The function is ANSI escape sequence aware.
func wrapLinePreservingSpaces(line string, width int) []string {
	if width <= 0 || line == "" {
		return []string{line}
	}

	type bp struct {
		afterByte int // byte index right after the space character
		visWidth  int // visible width up to and including the space
	}

	var result []string
	startByte := 0
	visWidth := 0
	var lastBreak *bp

	i := 0
	for i < len(line) {
		// Skip ANSI CSI sequences (ESC [ ... m)
		if line[i] == byte(ansiESC) && i+1 < len(line) && line[i+1] == byte(ansiCSI) {
			j := i + 2
			for j < len(line) && line[j] != byte(ansiEnd) {
				j++
			}
			if j < len(line) {
				i = j + 1
				continue
			}
		}

		_, size := utf8.DecodeRuneInString(line[i:])
		visWidth++

		if visWidth > width {
			if lastBreak != nil && lastBreak.afterByte > startByte {
				// Wrap at last space — space stays on current line
				result = append(result, line[startByte:lastBreak.afterByte])
				visWidth -= lastBreak.visWidth
				startByte = lastBreak.afterByte
				lastBreak = nil
			} else {
				// No break point; hard-break at current character
				result = append(result, line[startByte:i])
				startByte = i
				visWidth = 1
				lastBreak = nil
			}
		}

		// Record break point *after* overflow handling so the space that
		// triggers overflow goes to the new line instead of bloating the
		// current line beyond width.
		if line[i] == ' ' {
			lastBreak = &bp{afterByte: i + size, visWidth: visWidth}
		}

		i += size
	}

	result = append(result, line[startByte:])
	return result
}

// renderContainer wraps, substitutes whitespace glyphs, and pads the text to
// fill the view. nlStyles carries the match style for each hard newline in s
// (nil = unmatched, non-nil = use that style for the ↵ glyph). It may be nil
// when there is no active expression.
func (m *Model) renderContainer(s string, nlStyles []*lipgloss.Style) string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	var lines []string

	if m.showWhitespaces {
		// Split by hard newlines first, wrap each independently using a
		// space-preserving wrapper so that spaces at break points remain
		// visible as · glyphs. After wrapping, substitute glyphs and add ↵.
		hardLines := strings.Split(s, "\n")
		for i, hardLine := range hardLines {
			softLines := wrapLinePreservingSpaces(hardLine, m.width)
			for j := range softLines {
				softLines[j] = substituteWhitespace(softLines[j])
			}
			if i < len(hardLines)-1 {
				lastIdx := len(softLines) - 1
				if i < len(nlStyles) && nlStyles[i] != nil {
					softLines[lastIdx] += nlStyles[i].Render(newlineGlyph)
				} else {
					softLines[lastIdx] += whitespaceGlyphStyle.Render(newlineGlyph)
				}
			}
			lines = append(lines, softLines...)
		}
	} else {
		wrapped := wordwrap.String(s, m.width)
		lines = strings.Split(wrapped, "\n")
	}

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
		return m.renderContainer(m.value, nil)
	}

	var b strings.Builder
	lastIndex := 0

	// nlStyles tracks the match style for each hard newline in the value.
	// nil entries mean the newline is unmatched (dimmed ↵).
	var nlStyles []*lipgloss.Style

	var matches [][]int
	if m.global {
		matches = m.expression.FindAllStringIndex(m.value, -1)
	} else if idx := m.expression.FindStringIndex(m.value); idx != nil {
		matches = [][]int{idx}
	}
	for i, match := range matches {
		s := &evenMatchStyle
		if i%2 == 1 {
			s = &oddMatchStyle
		}

		// Unmatched text before this match — newlines here are unmatched.
		unmatched := m.value[lastIndex:match[0]]
		b.WriteString(unmatched)
		for range strings.Count(unmatched, "\n") {
			nlStyles = append(nlStyles, nil)
		}

		// Split match text at newlines so each hard-line gets its own
		// balanced ANSI open/close.  Without this, a match spanning \n
		// leaves an unclosed style that leaks into lipgloss padding.
		matchText := m.value[match[0]:match[1]]
		matchLines := strings.Split(matchText, "\n")
		for k, ml := range matchLines {
			if k > 0 {
				b.WriteRune('\n')
				nlStyles = append(nlStyles, s) // newline inside match
			}
			b.WriteString(s.Render(ml))
		}
		lastIndex = match[1]
	}

	// Trailing unmatched text — newlines here are unmatched.
	trailing := m.value[lastIndex:]
	b.WriteString(trailing)
	for range strings.Count(trailing, "\n") {
		nlStyles = append(nlStyles, nil)
	}

	return m.renderContainer(b.String(), nlStyles)
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

func (m *Model) SetShowWhitespaces(show bool) {
	m.showWhitespaces = show
}

func (m *Model) SetScrollYOffset(offset int) {
	m.scrollYOffset = offset
}

func (m *Model) Validate(expression string) error {
	_, err := m.newRegexp(expression)
	return err
}
