package screen

import "charm.land/bubbles/v2/key"

type keyMap struct {
	Exit          key.Binding
	SwitchInput   key.Binding
	ToggleOptions key.Binding
	OpenEditor    key.Binding
	CopyRegexp    key.Binding
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
	ToggleOptions: key.NewBinding(
		key.WithKeys("ctrl+p"),
		key.WithHelp("ctrl+p", "options"),
	),
	OpenEditor: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("ctrl+o", "edit text"),
	),
	CopyRegexp: key.NewBinding(
		key.WithKeys("ctrl+y"),
		key.WithHelp("ctrl+y", "copy regexp"),
	),
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Exit, k.SwitchInput},
		{k.ToggleOptions, k.OpenEditor},
	}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Exit, k.SwitchInput, k.ToggleOptions, k.OpenEditor, k.CopyRegexp}
}
