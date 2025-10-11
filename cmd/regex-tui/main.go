package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vitor-mariano/regex-tui/pkg/regexview"
)

type inputType int

const (
	inputTypeExpression inputType = iota
	inputTypeSubject
)

const (
	initialExpression = "([A-Z])\\w+"
	initialSubject    = "Hello World!"
)

type model struct {
	expressionInput textinput.Model
	subjectInput    textarea.Model
	subjectView     regexview.Model

	focusedInputType inputType
	expression       string
	subject          string
}

func initialModel() model {
	m := model{
		expressionInput: textinput.New(),
		subjectInput:    textarea.New(),
		subjectView:     regexview.New(),
	}

	m.expressionInput.SetValue(initialExpression)
	m.expressionInput.Prompt = ""
	m.expressionInput.Focus()

	m.subjectInput.SetValue(initialSubject)
	m.subjectInput.Prompt = ""
	m.subjectInput.ShowLineNumbers = false

	m.subjectView.SetExpressionString(initialExpression)
	m.subjectView.SetValue(initialSubject)

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.focusedInputType == inputTypeSubject {
		sm, cmd := m.subjectInput.Update(msg)
		m.subjectInput = sm
		m.subjectView.SetValue(m.subjectInput.Value())

		return m, cmd
	}

	cmds := make([]tea.Cmd, 2)

	m.expressionInput, cmds[0] = m.expressionInput.Update(msg)

	err := m.subjectView.SetExpressionString(m.expressionInput.Value())
	if err == nil {
		m.subjectView, cmds[1] = m.subjectView.Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab:
			switch m.focusedInputType {
			case inputTypeExpression:
				m.focusedInputType = inputTypeSubject
				m.expressionInput.Blur()
				m.subjectInput.Focus()
			case inputTypeSubject:
				m.focusedInputType = inputTypeExpression
				m.subjectInput.Blur()
				m.expressionInput.Focus()
			}
			return m, nil
		}
	}

	return m.updateInputs(msg)
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(m.expressionInput.View())
	b.WriteRune('\n')

	if m.focusedInputType == inputTypeSubject {
		b.WriteString(m.subjectInput.View())
	} else {
		b.WriteString(m.subjectView.View())
	}

	return b.String()
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
