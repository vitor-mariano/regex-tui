package options

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/vitor-mariano/regex-tui/internal/testkeys"
)

func TestOpenClose(t *testing.T) {
	m := New()
	if m.IsOpen() {
		t.Fatal("dialog should start closed")
	}

	m.Open()
	if !m.IsOpen() {
		t.Fatal("Open should mark the dialog open")
	}

	m.Update(testkeys.Press(tea.KeyEsc))
	if m.IsOpen() {
		t.Fatal("esc should close the dialog")
	}

	m.Open()
	m.Update(testkeys.Ctrl('p'))
	if m.IsOpen() {
		t.Fatal("ctrl+p should close the dialog")
	}
}

func TestOnToggleDoesNotForceGlobal(t *testing.T) {
	m := New()
	var called bool
	m.OnToggle(func(string, bool) { called = true })
	if called {
		t.Fatal("registering OnToggle should not select items")
	}
	if isChecked(m) {
		t.Fatal("Global should not be selected until SetSelected")
	}

	m.SetSelected(GlobalOption)
	if !isChecked(m) {
		t.Fatal("SetSelected(Global) should show a check mark")
	}
}

func TestUpdateTogglesCurrentItem(t *testing.T) {
	m := New()
	var gotItem string
	var gotSelected bool
	m.OnToggle(func(item string, selected bool) {
		gotItem = item
		gotSelected = selected
	})

	m.Update(testkeys.Press(tea.KeySpace))
	if gotItem != GlobalOption || !gotSelected {
		t.Fatalf("space toggle = %q selected=%v, want Global true", gotItem, gotSelected)
	}

	m.Update(testkeys.Press(tea.KeyDown))
	m.Update(testkeys.Rune('x'))
	if gotItem != InsensitiveOption || !gotSelected {
		t.Fatalf("x toggle = %q selected=%v, want Insensitive true", gotItem, gotSelected)
	}
}

func isChecked(m *Model) bool {
	return strings.Contains(ansi.Strip(m.View()), "✓")
}
