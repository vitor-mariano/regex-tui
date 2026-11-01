package screen

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/vitor-mariano/regex-tui/internal/testkeys"
)

func apply(m model, msg tea.Msg) (model, tea.Cmd) {
	next, cmd := m.Update(msg)
	return next.(model), cmd
}

func sized(t *testing.T, cfg Config) model {
	t.Helper()
	m, _ := apply(New(cfg), tea.WindowSizeMsg{Width: 80, Height: 24})
	return m
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func TestTabSwitchesFocus(t *testing.T) {
	m := sized(t, Config{InitialExpression: "a", InitialSubject: "abc"})
	if m.focusedInputType != inputTypeExpression {
		t.Fatal("focus should start on the expression")
	}

	m, _ = apply(m, testkeys.Press(tea.KeyTab))
	if m.focusedInputType != inputTypeSubject {
		t.Fatal("tab should focus the subject")
	}

	m, _ = apply(m, testkeys.Shift(tea.KeyTab))
	if m.focusedInputType != inputTypeExpression {
		t.Fatal("shift+tab should focus the expression")
	}
}

func TestOptionsOpenAndEscDoesNotQuit(t *testing.T) {
	m := sized(t, Config{})
	m, cmd := apply(m, testkeys.Ctrl('p'))
	if !m.options.IsOpen() {
		t.Fatal("ctrl+p should open options")
	}
	if isQuit(cmd) {
		t.Fatal("opening options should not quit")
	}

	m, cmd = apply(m, testkeys.Press(tea.KeyEsc))
	if m.options.IsOpen() {
		t.Fatal("esc should close options")
	}
	if isQuit(cmd) {
		t.Fatal("esc while options are open should not quit")
	}
}

func TestExitQuitsWhenOptionsClosed(t *testing.T) {
	m := sized(t, Config{})
	_, cmd := apply(m, testkeys.Press(tea.KeyEsc))
	if !isQuit(cmd) {
		t.Fatal("esc should quit when options are closed")
	}

	_, cmd = apply(m, testkeys.Ctrl('c'))
	if !isQuit(cmd) {
		t.Fatal("ctrl+c should quit when options are closed")
	}
}

func TestTypingUpdatesExpression(t *testing.T) {
	m := sized(t, Config{InitialExpression: "", InitialSubject: "foo"})
	before := m.subjectInput.GetView().View()
	for _, r := range "foo" {
		m, _ = apply(m, testkeys.Rune(r))
	}
	if got := m.expressionInput.GetInput().Value(); got != "foo" {
		t.Fatalf("expression = %q, want foo", got)
	}
	if m.subjectInput.GetView().View() == before {
		t.Fatal("typing should push the expression into the subject view")
	}
}

func TestConfigGlobal(t *testing.T) {
	cfg := Config{
		InitialExpression: `\w+`,
		InitialSubject:    "one two",
	}
	off := sized(t, cfg)
	cfg.Global = true
	on := sized(t, cfg)
	if off.subjectInput.GetView().View() == on.subjectInput.GetView().View() {
		t.Fatal("Global=false should not highlight every match")
	}
}

func TestConfigInsensitive(t *testing.T) {
	cfg := Config{
		InitialExpression: "hello",
		InitialSubject:    "HELLO",
	}
	off := sized(t, cfg)
	cfg.Insensitive = true
	on := sized(t, cfg)
	if off.subjectInput.GetView().View() == on.subjectInput.GetView().View() {
		t.Fatal("Insensitive should change matching")
	}
}

func TestConfigWhitespaces(t *testing.T) {
	cfg := Config{
		InitialExpression: "Hello",
		InitialSubject:    "Hello World",
	}
	off := sized(t, cfg)
	cfg.Whitespaces = true
	on := sized(t, cfg)
	if strings.Contains(ansi.Strip(off.subjectInput.GetView().View()), "·") {
		t.Fatal("Whitespaces=false should not show space glyphs")
	}
	if !strings.Contains(ansi.Strip(on.subjectInput.GetView().View()), "·") {
		t.Fatal("Whitespaces=true should show space glyphs")
	}
}

func TestConfigRegexp2(t *testing.T) {
	cfg := Config{
		InitialExpression: `foo(?=bar)`,
		InitialSubject:    "foobar foobaz",
		Global:            true,
	}
	off := sized(t, cfg)
	cfg.Regexp2 = true
	on := sized(t, cfg)
	if off.subjectInput.GetView().View() == on.subjectInput.GetView().View() {
		t.Fatal("Regexp2 should compile lookahead that RE2 rejects")
	}
}
