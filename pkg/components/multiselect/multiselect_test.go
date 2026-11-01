package multiselect

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/vitor-mariano/regex-tui/internal/testkeys"
)

func TestNavigationWraps(t *testing.T) {
	m := New([]string{"a", "b", "c"})
	if m.current != 0 {
		t.Fatalf("current = %d, want 0", m.current)
	}

	m.Update(testkeys.Press(tea.KeyDown))
	if m.current != 1 {
		t.Fatalf("down current = %d, want 1", m.current)
	}

	m.Update(testkeys.Rune('j'))
	m.Update(testkeys.Rune('j'))
	if m.current != 0 {
		t.Fatalf("down wrap current = %d, want 0", m.current)
	}

	m.Update(testkeys.Press(tea.KeyUp))
	if m.current != 2 {
		t.Fatalf("up wrap current = %d, want 2", m.current)
	}

	m.Update(testkeys.Rune('k'))
	if m.current != 1 {
		t.Fatalf("up current = %d, want 1", m.current)
	}
}

func TestToggle(t *testing.T) {
	m := New([]string{"a", "b"})
	var events []struct {
		item     string
		selected bool
	}
	m.OnToggle(func(item string, selected bool) {
		events = append(events, struct {
			item     string
			selected bool
		}{item, selected})
	})

	m.Update(testkeys.Press(tea.KeySpace))
	if !m.selected.Contains("a") {
		t.Fatal("space should select the current item")
	}
	if len(events) != 1 || events[0].item != "a" || !events[0].selected {
		t.Fatalf("onToggle events = %#v", events)
	}

	m.Update(testkeys.Rune('x'))
	if m.selected.Contains("a") {
		t.Fatal("x should deselect the current item")
	}
	if len(events) != 2 || events[1].selected {
		t.Fatalf("onToggle events = %#v", events)
	}
}

func TestSetSelectedIsAdditive(t *testing.T) {
	m := New([]string{"a", "b", "c"})
	m.SetSelected("a")
	m.SetSelected("c")
	if !m.selected.Contains("a") || !m.selected.Contains("c") || m.selected.Contains("b") {
		t.Fatal("SetSelected should add without clearing previous items")
	}
}
