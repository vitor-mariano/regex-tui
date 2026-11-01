package re2_test

import (
	"reflect"
	"testing"

	"github.com/vitor-mariano/regex-tui/pkg/regex/re2"
)

func TestNewInvalid(t *testing.T) {
	if _, err := re2.New("["); err == nil {
		t.Fatal("expected compile error for invalid pattern")
	}
}

func TestNewLookaheadUnsupported(t *testing.T) {
	if _, err := re2.New(`foo(?=bar)`); err == nil {
		t.Fatal("expected compile error for lookahead")
	}
}

func TestFindIndexes(t *testing.T) {
	tests := []struct {
		name  string
		expr  string
		input string
		n     int
		all   [][]int
		first []int
	}{
		{
			name:  "no match",
			expr:  "xyz",
			input: "abc",
			n:     -1,
			all:   nil,
			first: nil,
		},
		{
			name:  "single match",
			expr:  `World`,
			input: "Hello World!",
			n:     -1,
			all:   [][]int{{6, 11}},
			first: []int{6, 11},
		},
		{
			name:  "global words",
			expr:  `\w+`,
			input: "Hello World",
			n:     -1,
			all:   [][]int{{0, 5}, {6, 11}},
			first: []int{0, 5},
		},
		{
			name:  "limit n",
			expr:  `\w+`,
			input: "Hello World",
			n:     1,
			all:   [][]int{{0, 5}},
			first: []int{0, 5},
		},
		{
			name:  "n is zero",
			expr:  `\w+`,
			input: "Hello World",
			n:     0,
			all:   nil,
			first: []int{0, 5},
		},
		{
			name:  "unicode",
			expr:  `世界`,
			input: "Hello 世界",
			n:     -1,
			all:   [][]int{{6, 12}},
			first: []int{6, 12},
		},
		{
			name:  "empty input",
			expr:  `a+`,
			input: "",
			n:     -1,
			all:   nil,
			first: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re, err := re2.New(tt.expr)
			if err != nil {
				t.Fatalf("New(%q): %v", tt.expr, err)
			}

			all := re.FindAllStringIndex(tt.input, tt.n)
			if !reflect.DeepEqual(all, tt.all) {
				t.Errorf("FindAllStringIndex = %#v, want %#v", all, tt.all)
			}

			first := re.FindStringIndex(tt.input)
			if !reflect.DeepEqual(first, tt.first) {
				t.Errorf("FindStringIndex = %#v, want %#v", first, tt.first)
			}
		})
	}
}
