package main

import (
	"errors"
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

var errTextWithStdin = errors.New("error: cannot use --text/-t flag when reading from stdin")

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
	cfg, err := parseConfig(os.Args[1:], os.Stdin, hasStdin())
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		log.Fatal(err)
	}
	return cfg
}

func parseConfig(args []string, stdin io.Reader, stdinIsPipe bool) (screen.Config, error) {
	fs := flag.NewFlagSet("regex-tui", flag.ContinueOnError)

	empty := fs.Bool("empty", false, "Start with empty expression and text")
	fs.BoolVar(empty, "e", false, "Start with empty expression and text (shorthand)")

	regex := fs.String("regex", "", "Initial regex pattern")
	fs.StringVar(regex, "r", "", "Initial regex pattern (shorthand)")

	text := fs.String("text", "", "Initial text subject")
	fs.StringVar(text, "t", "", "Initial text subject (shorthand)")

	noGlobal := fs.Bool("no-global", false, "Disable global flag (match only first occurrence)")
	insensitive := fs.Bool("insensitive", false, "Enable case-insensitive flag")
	regexp2 := fs.Bool("regexp2", false, "Use regexp2 engine (partial PCRE compatibility)")
	noWhitespaces := fs.Bool("no-whitespaces", false, "Disable whitespace visualization")

	if err := fs.Parse(args); err != nil {
		return screen.Config{}, err
	}

	if stdinIsPipe && *text != "" {
		return screen.Config{}, errTextWithStdin
	}

	var regexExpression string
	if *regex != "" {
		regexExpression = *regex
	} else if !*empty {
		regexExpression = defaultRegex
	}

	var textSubject string
	if stdinIsPipe {
		data, err := io.ReadAll(stdin)
		if err != nil {
			return screen.Config{}, fmt.Errorf("failed to read from stdin: %w", err)
		}
		textSubject = string(data)
	} else if *text != "" {
		textSubject = *text
	} else if !*empty {
		textSubject = defaultText
	}

	return screen.Config{
		InitialExpression: regexExpression,
		InitialSubject:    textSubject,
		Global:            !*noGlobal,
		Insensitive:       *insensitive,
		Regexp2:           *regexp2,
		Whitespaces:       !*noWhitespaces,
	}, nil
}
