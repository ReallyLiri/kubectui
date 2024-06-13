package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/kubectui/tui/keymap"
)

const maxWidth = 250

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch tmsg := msg.(type) {
	case tea.WindowSizeMsg:
		if tmsg.Width > 0 && tmsg.Height > 0 {
			m.state.termWidth = tmsg.Width
			if m.state.termWidth > maxWidth {
				m.state.termWidth = maxWidth
			}
			m.state.termHeight = tmsg.Height
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(tmsg, keymap.Tab), key.Matches(tmsg, keymap.Left), key.Matches(tmsg, keymap.Right):
			m.state.focused = (m.state.focused + 1) % ComponentCount
		case key.Matches(tmsg, keymap.Up), key.Matches(tmsg, keymap.Down):
			switch m.state.focused {
			case ContextList:
				m.vms.contextList, cmd = m.vms.contextList.Update(msg)
				nextContext := m.vms.contextList.SelectedItem().FilterValue()
				m.onContextSelected(nextContext)
			case NamespaceList:
				m.vms.namespaceList, cmd = m.vms.namespaceList.Update(msg)
			}
		case key.Matches(tmsg, keymap.Help):
			m.vms.help.ShowAll = !m.vms.help.ShowAll
		case key.Matches(tmsg, keymap.Quit):
			m.state.quitting = true
			return m, tea.Quit
		}
	}
	return m, cmd
}
