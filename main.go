package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/vitor-mariano/regex-tui/internal/screen"
	"github.com/vitor-mariano/regex-tui/internal/tty"
)

const (
	defaultRegex = "[A-Z]\\w+"
	defaultText  = "Hello World!"
)

func main() {
	hasStdin := hasStdin()

	config := getInitialConfig()

	var p *tea.Program
	if hasStdin {
		tty, err := tty.OpenInputTTY()
		if err != nil {
			log.Fatalf("failed to open TTY: %v\n", err)
		}
		defer tty.Close()

		p = tea.NewProgram(
			screen.New(config),
			tea.WithInput(tty),
		)
	} else {
		p = tea.NewProgram(screen.New(config))
	}

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start program: %v\n", err)
		os.Exit(1)
	}
}

func hasStdin() bool {
	stdinStat, _ := os.Stdin.Stat()
	return (stdinStat.Mode() & os.ModeCharDevice) == 0
}

func getInitialConfig() screen.Config {
	empty := flag.Bool("empty", false, "Start with empty expression and text")
	flag.BoolVar(empty, "e", false, "Start with empty expression and text (shorthand)")

	regex := flag.String("regex", "", "Initial regex pattern")
	flag.StringVar(regex, "r", "", "Initial regex pattern (shorthand)")

	text := flag.String("text", "", "Initial text subject")
	flag.StringVar(text, "t", "", "Initial text subject (shorthand)")

	noGlobal := flag.Bool("no-global", false, "Disable global flag (match only first occurrence)")

	insensitive := flag.Bool("insensitive", false, "Enable case-insensitive flag")

	regexp2 := flag.Bool("regexp2", false, "Use regexp2 engine (partial PCRE compatibility)")

	noWhitespaces := flag.Bool("no-whitespaces", false, "Disable whitespace visualization")

	flag.Parse()

	if hasStdin() && *text != "" {
		log.Fatal("error: cannot use --text/-t flag when reading from stdin")
	}

	var regexExpression string
	if *regex != "" {
		regexExpression = *regex
	} else if !*empty {
		regexExpression = defaultRegex
	}

	var textSubject string
	if hasStdin() {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("failed to read from stdin: %v\n", err)
		}
		textSubject = string(data)
	} else if *text != "" {
		textSubject = *text
	} else if !*empty {
		textSubject = defaultText
	}

	global := !*noGlobal
	whitespaces := !*noWhitespaces

	return screen.Config{
		InitialExpression: regexExpression,
		InitialSubject:    textSubject,
		Global:            global,
		Insensitive:       *insensitive,
		Regexp2:           *regexp2,
		Whitespaces:       whitespaces,
	}
}
