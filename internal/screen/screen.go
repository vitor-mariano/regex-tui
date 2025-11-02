package screen

import (
	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/components/expression"
	"github.com/vitor-mariano/regex-tui/internal/components/subject"
)

type inputType int

const (
	inputTypeExpression inputType = iota
	inputTypeSubject
)

const (
	initialExpression = "[A-Z]\\w+"
	initialSubject    = "Hello World!"
)

type model struct {
	expressionInput expression.Model
	subjectInput    subject.Model
	help            help.Model

	focusedInputType inputType
}

func New() model {
	ei := expression.New(initialExpression)
	ei.GetInput().Focus()

	return model{
		expressionInput: ei,
		subjectInput:    subject.New(initialSubject, initialExpression),
		help:            help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return m.expressionInput.Init()
}

func (m *model) setSize(width, height int) {
	const subjectVSpacing = 8

	m.expressionInput.SetWidth(width)
	m.subjectInput.SetSize(width, height-subjectVSpacing)
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if m.focusedInputType == inputTypeSubject {
		m.subjectInput, cmd = m.subjectInput.Update(msg)

		return cmd
	}

	m.expressionInput, cmd = m.expressionInput.Update(msg)
	m.subjectInput.SetExpressionString(m.expressionInput.GetInput().Value())

	return cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 2)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Exit):
			return m, tea.Quit

		case key.Matches(msg, keys.SwitchInput):
			var cmd tea.Cmd
			switch m.focusedInputType {
			case inputTypeExpression:
				m.focusedInputType = inputTypeSubject
				m.expressionInput.GetInput().Blur()
				cmd = m.subjectInput.GetInput().Focus()

			case inputTypeSubject:
				m.focusedInputType = inputTypeExpression
				m.subjectInput.GetInput().Blur()
				cmd = m.expressionInput.GetInput().Focus()
			}

			cmds = append(cmds, cmd)
		}
	}

	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		m.expressionInput.View(),
		m.subjectInput.View(),
		m.help.View(keys),
	))
}
