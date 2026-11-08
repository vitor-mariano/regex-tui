package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/vitor-mariano/regex-tui/internal/screen"
)

func TestParseConfigDefaults(t *testing.T) {
	cfg, err := parseConfig(nil, strings.NewReader(""), false)
	if err != nil {
		t.Fatal(err)
	}
	want := screen.Config{
		InitialExpression: defaultRegex,
		InitialSubject:    defaultText,
		Global:            true,
		Whitespaces:       true,
	}
	if cfg != want {
		t.Fatalf("got %#v, want %#v", cfg, want)
	}
}

func TestParseConfigEmpty(t *testing.T) {
	cfg, err := parseConfig([]string{"-e"}, strings.NewReader(""), false)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.InitialExpression != "" || cfg.InitialSubject != "" {
		t.Fatalf("empty flag should clear expression and text, got %#v", cfg)
	}
}

func TestParseConfigRegexAndText(t *testing.T) {
	cfg, err := parseConfig([]string{"-r", `\d+`, "-t", "abc 123"}, strings.NewReader(""), false)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.InitialExpression != `\d+` || cfg.InitialSubject != "abc 123" {
		t.Fatalf("got expr=%q text=%q", cfg.InitialExpression, cfg.InitialSubject)
	}
}

func TestParseConfigFlags(t *testing.T) {
	cfg, err := parseConfig([]string{
		"--no-global",
		"--insensitive",
		"--regexp2",
		"--no-whitespaces",
	}, strings.NewReader(""), false)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Global {
		t.Fatal("expected Global=false")
	}
	if !cfg.Insensitive {
		t.Fatal("expected Insensitive=true")
	}
	if !cfg.Regexp2 {
		t.Fatal("expected Regexp2=true")
	}
	if cfg.Whitespaces {
		t.Fatal("expected Whitespaces=false")
	}
}

func TestParseConfigStdin(t *testing.T) {
	cfg, err := parseConfig([]string{"-r", "err"}, strings.NewReader("piped text"), true)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.InitialSubject != "piped text" {
		t.Fatalf("subject = %q, want piped text", cfg.InitialSubject)
	}
	if cfg.InitialExpression != "err" {
		t.Fatalf("expression = %q, want err", cfg.InitialExpression)
	}
}

func TestParseConfigStdinRejectsTextFlag(t *testing.T) {
	_, err := parseConfig([]string{"-t", "nope"}, strings.NewReader("piped"), true)
	if !errors.Is(err, errTextWithStdin) {
		t.Fatalf("err = %v, want errTextWithStdin", err)
	}
}
