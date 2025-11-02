package screen

import "github.com/charmbracelet/bubbles/v2/key"

type keyMap struct {
	Exit        key.Binding
	SwitchInput key.Binding
}

var keys = keyMap{
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("esc/ctrl+c", "exit"),
	),
	SwitchInput: key.NewBinding(
		key.WithKeys("tab", "shift+tab"),
		key.WithHelp("tab", "switch input"),
	),
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Exit, k.SwitchInput},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Exit, k.SwitchInput}
}
