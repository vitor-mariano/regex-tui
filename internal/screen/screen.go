package screen

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/vitor-mariano/regex-tui/internal/clipboard"
	"github.com/vitor-mariano/regex-tui/internal/components/expression"
	"github.com/vitor-mariano/regex-tui/internal/components/options"
	"github.com/vitor-mariano/regex-tui/internal/components/subject"
	"github.com/vitor-mariano/regex-tui/pkg/components/multiselect"
)

type inputType int

const (
	inputTypeExpression inputType = iota
	inputTypeSubject
)

type editorFinishedMsg struct {
	tempFile string
	err      error
}

type Config struct {
	InitialExpression string
	InitialSubject    string
	Global            bool
	Insensitive       bool
	Regexp2           bool
	Whitespaces       bool
}

type model struct {
	expressionInput *expression.Model
	subjectInput    *subject.Model
	options         *options.Model
	help            help.Model

	focusedInputType inputType
	width, height    int
}

func New(config Config) model {
	si := subject.New(config.InitialSubject, config.InitialExpression)

	ei := expression.New(config.InitialExpression, si.GetView())
	ei.GetInput().Focus()

	d := options.New()
	d.OnToggle(func(item string, selected bool) {
		switch item {
		case options.GlobalOption:
			si.GetView().SetGlobal(selected)
		case options.InsensitiveOption:
			si.GetView().SetInsensitive(selected)
		case options.Regexp2Option:
			si.GetView().SetRegexp2(selected)
			ei.GetInput().Err = si.GetView().SetRegexp2(selected)
			// Force re-evaluation with the new engine.
			si.SetExpression(ei.GetInput().Value())
		case options.WhitespacesOption:
			si.GetView().SetShowWhitespaces(selected)
		}
	})

	var selectedOptions []string
	if config.Global {
		selectedOptions = append(selectedOptions, options.GlobalOption)
	}
	if config.Insensitive {
		selectedOptions = append(selectedOptions, options.InsensitiveOption)
	}
	if config.Regexp2 {
		selectedOptions = append(selectedOptions, options.Regexp2Option)
	}
	if config.Whitespaces {
		selectedOptions = append(selectedOptions, options.WhitespacesOption)
	}

	if len(selectedOptions) > 0 {
		d.SetSelected(selectedOptions...)
	}

	return model{
		expressionInput: ei,
		subjectInput:    si,
		options:         d,
		help:            help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.expressionInput.Init(),
		m.subjectInput.Init(),
	)
}

func (m *model) setSize(width, height int) {
	const subjectVSpacing = 8

	m.width = width
	m.height = height
	m.expressionInput.SetWidth(width)
	m.subjectInput.SetSize(width, height-subjectVSpacing)
}

func findEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	fallbacks := []string{"vim", "vi", "emacs"}
	for _, editor := range fallbacks {
		if path, err := exec.LookPath(editor); err == nil {
			return path
		}
	}

	return "nano"
}

func (m *model) openEditor() tea.Cmd {
	editor := findEditor()

	tmpFile, err := os.CreateTemp("", "regex-tui-*.txt")
	if err != nil {
		return nil
	}
	defer tmpFile.Close()

	content := m.subjectInput.GetInput().Value()
	if _, err := tmpFile.WriteString(content); err != nil {
		os.Remove(tmpFile.Name())
		return nil
	}

	return tea.ExecProcess(
		exec.Command(editor, tmpFile.Name()),
		func(err error) tea.Msg {
			return editorFinishedMsg{tempFile: tmpFile.Name(), err: err}
		},
	)
}

func (m *model) updateScreen(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0, 2)

	var shouldCopyRegexp bool
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.SwitchInput):
			var cmd tea.Cmd
			switch m.focusedInputType {
			case inputTypeExpression:
				m.focusedInputType = inputTypeSubject
				m.expressionInput.GetInput().Blur()
				cmd = m.subjectInput.Focus()

			case inputTypeSubject:
				m.focusedInputType = inputTypeExpression
				m.subjectInput.Blur()
				cmd = m.expressionInput.GetInput().Focus()
			}

			cmds = append(cmds, cmd)

		case key.Matches(msg, keys.ToggleOptions):
			if !m.options.IsOpen() {
				m.options.Open()
			}

		case key.Matches(msg, keys.CopyRegexp):
			shouldCopyRegexp = true

		case key.Matches(msg, keys.OpenEditor):
			return m.openEditor()
		}
	}

	if m.focusedInputType == inputTypeSubject {
		cmds = append(cmds, m.subjectInput.Update(msg))
	} else {
		cmds = append(cmds, m.expressionInput.Update(msg))

		// Check whether should copy the regexp
		if shouldCopyRegexp {
			if err := clipboard.WriteAll(m.expressionInput.GetInput().Value()); err != nil {
				fmt.Fprintf(os.Stderr, "failed to write to clipboard: %v\n", err)
			}
		}

		m.subjectInput.SetExpression(m.expressionInput.GetInput().Value())
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 2)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setSize(msg.Width, msg.Height)

	case editorFinishedMsg:
		if msg.err == nil {
			content, err := os.ReadFile(msg.tempFile)
			if err == nil {
				m.subjectInput.SetValue(string(content))
			}
		}
		os.Remove(msg.tempFile)
		return m, nil

	case tea.KeyPressMsg:
		if key.Matches(msg, keys.Exit) {
			if m.options.IsOpen() {
				break
			}

			return m, tea.Quit
		}
	}

	if m.options.IsOpen() {
		cmds = append(cmds, m.options.Update(msg))
	} else {
		cmds = append(cmds, m.updateScreen(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var helpKeyMap help.KeyMap = keys
	if m.options.IsOpen() {
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
	if m.options.IsOpen() {
		optionsLayer := lipgloss.NewLayer(m.options.View())
		optionsLayer.X((m.width - optionsLayer.Width()) / 2)
		optionsLayer.Y((m.height - optionsLayer.Height()) / 2)

		layers = append(layers, optionsLayer)
	}

	return tea.NewView(lipgloss.NewCompositor(layers...).Render())
}
