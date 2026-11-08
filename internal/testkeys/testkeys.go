// Package testkeys builds Bubble Tea v2 key messages for tests.
package testkeys

import tea "charm.land/bubbletea/v2"

func Press(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func Rune(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Text: string(r), Code: r})
}

func Ctrl(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Mod: tea.ModCtrl})
}

func Shift(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code, Mod: tea.ModShift})
}
