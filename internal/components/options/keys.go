package options

import "github.com/charmbracelet/bubbles/v2/key"

type keyMap struct {
	Exit key.Binding
}

var keys = keyMap{
	Exit: key.NewBinding(
		key.WithKeys("esc", "ctrl+p"),
	),
}
