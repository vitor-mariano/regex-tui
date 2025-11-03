package screen

import (
	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/components/expression"
	"github.com/vitor-mariano/regex-tui/internal/components/subject"
	"github.com/vitor-mariano/regex-tui/internal/styles"
	"github.com/vitor-mariano/regex-tui/pkg/components/multiselect"
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
	expressionInput     *expression.Model
	subjectInput        *subject.Model
	options             *multiselect.Model
	isOptionsDialogOpen bool
	help                help.Model

	focusedInputType inputType
	width, height    int
}

const (
	globalOption      = "Global"
	insensitiveOption = "Insensitive"
)

func New() model {
	ei := expression.New(initialExpression)
	ei.GetInput().Focus()

	si := subject.New(initialSubject, initialExpression)

	mi := multiselect.New([]string{globalOption, insensitiveOption})
	mi.OnToggle(func(item string, selected bool) {
		switch item {
		case globalOption:
			si.GetView().SetGlobal(selected)
		case insensitiveOption:
			si.GetView().SetInsensitive(selected)
		}
	})
	mi.SetSelected(globalOption)

	return model{
		expressionInput: ei,
		subjectInput:    si,
		options:         mi,
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
		cmd = m.subjectInput.Update(msg)

		return cmd
	}

	cmd = m.expressionInput.Update(msg)
	m.subjectInput.SetExpression(m.expressionInput.GetInput().Value())

	return cmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 2)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.Exit):
			return m, tea.Quit

		case key.Matches(msg, keys.SwitchInput):
			if m.isOptionsDialogOpen {
				break
			}

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

		case key.Matches(msg, keys.ToggleOptions):
			m.isOptionsDialogOpen = !m.isOptionsDialogOpen
		}
	}

	if m.isOptionsDialogOpen {
		cmds = append(cmds, m.options.Update(msg))
	} else {
		cmds = append(cmds, m.updateInputs(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var helpKeyMap help.KeyMap = keys
	if m.isOptionsDialogOpen {
		helpKeyMap = multiselect.Keys
	}

	baseLayer := lipgloss.NewLayer(lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		m.expressionInput.View(),
		m.subjectInput.View(),
		m.help.View(helpKeyMap),
	))

	layers := []*lipgloss.Layer{baseLayer}
	if m.isOptionsDialogOpen {
		optionsLayer := lipgloss.NewLayer(
			lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(styles.MutedColor).
				Padding(1, 4, 1, 2).
				Render(m.options.View()),
		)
		optionsLayer.X((m.width - optionsLayer.GetWidth()) / 2)
		optionsLayer.Y((m.height - optionsLayer.GetHeight()) / 2)

		layers = append(layers, optionsLayer)
	}

	return tea.NewView(lipgloss.NewCanvas(layers...).Render())
}
