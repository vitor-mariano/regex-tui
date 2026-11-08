package regexview

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
)

func TestSubstituteWhitespace(t *testing.T) {
	plain := substituteWhitespace("a b")
	if strings.Contains(plain, " ") {
		t.Fatalf("unstyled space should become a glyph, got %q", plain)
	}
	if !strings.Contains(ansi.Strip(plain), "a·b") {
		t.Fatalf("unstyled space glyph = %q, want a·b", ansi.Strip(plain))
	}

	styled := evenMatchStyle.Render("a b")
	got := substituteWhitespace(styled)
	stripped := ansi.Strip(got)
	if !strings.Contains(stripped, "a·b") {
		t.Fatalf("styled space glyph = %q, want a·b", stripped)
	}
	if strings.Contains(got, whitespaceGlyphStyle.Render("·")) {
		t.Fatal("space inside a match should inherit match style, not the dim glyph style")
	}
}

func TestWrapLinePreservingSpaces(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		width int
		want  []string
	}{
		{
			name:  "empty",
			line:  "",
			width: 10,
			want:  []string{""},
		},
		{
			name:  "non-positive width",
			line:  "hello",
			width: 0,
			want:  []string{"hello"},
		},
		{
			name:  "fits",
			line:  "hello world",
			width: 11,
			want:  []string{"hello world"},
		},
		{
			name:  "wrap keeps trailing space",
			line:  "hello world",
			width: 10,
			want:  []string{"hello ", "world"},
		},
		{
			name:  "hard break when no space",
			line:  "hello",
			width: 3,
			want:  []string{"hel", "lo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapLinePreservingSpaces(tt.line, tt.width)
			if len(got) != len(tt.want) {
				t.Fatalf("got %#v, want %#v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got %#v, want %#v", got, tt.want)
				}
			}
		})
	}
}

func TestWrapLinePreservingSpacesSkipsANSI(t *testing.T) {
	line := evenMatchStyle.Render("abcd")
	got := wrapLinePreservingSpaces(line, 2)
	if len(got) != 2 {
		t.Fatalf("got %d pieces, want 2: %#v", len(got), got)
	}
	if ansi.Strip(got[0]) != "ab" || ansi.Strip(got[1]) != "cd" {
		t.Fatalf("ANSI wrap split = %q / %q, want ab / cd", ansi.Strip(got[0]), ansi.Strip(got[1]))
	}
}

func TestViewEmptySize(t *testing.T) {
	m := New(0, 8)
	m.SetValue("hello")
	if m.View() != "" {
		t.Fatal("width 0 should render empty")
	}
}

func TestViewWithoutExpression(t *testing.T) {
	m := New(40, 8)
	m.SetValue("hello")
	if !strings.Contains(ansi.Strip(m.View()), "hello") {
		t.Fatalf("View = %q", ansi.Strip(m.View()))
	}
}

func TestViewNoMatchWithoutGlobal(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression("xyz"); err != nil {
		t.Fatal(err)
	}
	m.SetValue("hello")
	m.SetGlobal(false)

	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("View panicked: %v", rec)
		}
	}()

	got := ansi.Strip(m.View())
	if !strings.Contains(got, "hello") {
		t.Fatalf("View = %q, want hello", got)
	}
}

func TestViewGlobalAlternatingStyles(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression(`\w+`); err != nil {
		t.Fatal(err)
	}
	m.SetValue("one two")
	m.SetGlobal(true)

	view := m.View()
	if !strings.Contains(view, evenMatchStyle.Render("one")) {
		t.Fatal("first match should use evenMatchStyle")
	}
	if !strings.Contains(view, oddMatchStyle.Render("two")) {
		t.Fatal("second match should use oddMatchStyle")
	}
}

func TestViewFirstMatchOnly(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression(`\w+`); err != nil {
		t.Fatal(err)
	}
	m.SetValue("one two")
	m.SetGlobal(false)

	view := m.View()
	if !strings.Contains(view, evenMatchStyle.Render("one")) {
		t.Fatal("first match should be highlighted")
	}
	if strings.Contains(view, evenMatchStyle.Render("two")) || strings.Contains(view, oddMatchStyle.Render("two")) {
		t.Fatal("second match should not be highlighted when global is off")
	}
}

func TestViewWhitespaceGlyphs(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression("Hello"); err != nil {
		t.Fatal(err)
	}
	m.SetValue("Hello World\nNext")
	m.SetGlobal(true)
	m.SetShowWhitespaces(true)

	got := ansi.Strip(m.View())
	if !strings.Contains(got, "Hello·World") {
		t.Fatalf("missing space glyph: %q", got)
	}
	if !strings.Contains(got, "↵") {
		t.Fatalf("missing newline glyph: %q", got)
	}
}

func TestViewScrollOffset(t *testing.T) {
	m := New(10, 2)
	m.SetValue("a\nb\nc\nd")
	m.SetScrollYOffset(2)

	got := ansi.Strip(m.View())
	if strings.Contains(got, "a") {
		t.Fatalf("scrolled view should not show first line: %q", got)
	}
	if !strings.Contains(got, "c") {
		t.Fatalf("scrolled view should show offset line: %q", got)
	}
}

func TestSetExpressionInvalidKeepsPrevious(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression("foo"); err != nil {
		t.Fatal(err)
	}
	m.SetValue("foo bar")
	m.SetGlobal(true)
	before := m.View()

	if err := m.SetExpression("["); err == nil {
		t.Fatal("expected compile error")
	}
	if m.View() != before {
		t.Fatal("invalid pattern should keep the previous compiled regex")
	}
	if m.baseExpStr != "foo" {
		t.Fatalf("baseExpStr = %q, want foo", m.baseExpStr)
	}
}

func TestValidate(t *testing.T) {
	m := New(40, 8)
	if err := m.Validate(`\w+`); err != nil {
		t.Fatal(err)
	}
	if err := m.Validate("["); err == nil {
		t.Fatal("expected invalid pattern")
	}
}

func TestSetInsensitive(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression("hello"); err != nil {
		t.Fatal(err)
	}
	m.SetValue("HELLO")
	m.SetGlobal(true)

	sensitive := m.View()
	m.SetInsensitive(true)
	insensitive := m.View()
	if sensitive == insensitive {
		t.Fatal("SetInsensitive should change matching")
	}
	if !strings.Contains(insensitive, evenMatchStyle.Render("HELLO")) {
		t.Fatal("case-insensitive match should highlight HELLO")
	}
}

func TestSetRegexp2Lookahead(t *testing.T) {
	m := New(40, 8)
	if err := m.SetRegexp2(true); err != nil {
		t.Fatal(err)
	}
	if err := m.SetExpression(`foo(?=bar)`); err != nil {
		t.Fatal(err)
	}
	m.SetValue("foobar foobaz")
	m.SetGlobal(true)

	view := m.View()
	if !strings.Contains(view, evenMatchStyle.Render("foo")) {
		t.Fatal("lookahead should match foo in foobar")
	}
	stripped := ansi.Strip(view)
	if !strings.Contains(stripped, "foobar") || !strings.Contains(stripped, "foobaz") {
		t.Fatalf("View = %q", stripped)
	}
}

func TestViewGoldenGlobalMatches(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression(`[A-Z]\w+`); err != nil {
		t.Fatal(err)
	}
	m.SetValue("Hello World!")
	m.SetGlobal(true)
	golden.RequireEqual(t, m.View())
}

func TestViewGoldenWhitespaces(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression(`[A-Z]\w+`); err != nil {
		t.Fatal(err)
	}
	m.SetValue("Hello World\nNext line")
	m.SetGlobal(true)
	m.SetShowWhitespaces(true)
	golden.RequireEqual(t, m.View())
}

func TestViewGoldenFirstMatchOnly(t *testing.T) {
	m := New(40, 8)
	if err := m.SetExpression(`\w+`); err != nil {
		t.Fatal(err)
	}
	m.SetValue("one two three")
	m.SetGlobal(false)
	golden.RequireEqual(t, m.View())
}
