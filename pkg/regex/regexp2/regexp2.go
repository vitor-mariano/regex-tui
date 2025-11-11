package regexp2

import (
	"github.com/dlclark/regexp2"
)

type Regexp2Regex struct {
	re *regexp2.Regexp
}

func New(expr string) (*Regexp2Regex, error) {
	re, err := regexp2.Compile(expr, 0)
	if err != nil {
		return nil, err
	}

	return &Regexp2Regex{re}, nil
}

func (regex *Regexp2Regex) FindAllStringIndex(s string, n int) [][]int {
	var matches [][]int
	match, err := regex.re.FindStringMatch(s)
	if err != nil {
		return matches
	}

	count := 0
	for match != nil && (n < 0 || count < n) {
		matches = append(matches, []int{match.Index, match.Index + match.Length})
		match, err = regex.re.FindNextMatch(match)
		if err != nil {
			break
		}
		count++
	}

	return matches
}

func (regex *Regexp2Regex) FindStringIndex(s string) []int {
	match, err := regex.re.FindStringMatch(s)
	if err != nil || match == nil {
		return nil
	}

	return []int{match.Index, match.Index + match.Length}
}
