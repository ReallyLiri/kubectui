package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/kubectui/tui/keymap"
	"strings"
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
		case m.state.deleting || m.state.renaming:
			switch {
			case key.Matches(tmsg, keymap.Select):
				value := strings.ToLower(strings.TrimSpace(m.vms.input.Value()))
				if m.state.deleting && (value == "y" || value == "yes") {
					m.deleteSelectedContext()
				}
				if m.state.renaming && value != "" {
					m.renameSelectedContext(value)
				}
				m.finishAction()
			case key.Matches(tmsg, keymap.Cancel):
				m.finishAction()
			default:
				m.vms.input, cmd = m.vms.input.Update(msg)
			}
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
				nextNamespace := m.vms.namespaceList.SelectedItem().FilterValue()
				m.onNamespaceSelected(nextNamespace)
			}
		case key.Matches(tmsg, keymap.Select):
			switch m.state.focused {
			case ContextList:
				m.setCurrentContextFromSelected()
			case NamespaceList:
				m.setCurrentNamespaceFromSelected()
			}
		case key.Matches(tmsg, keymap.Rename):
			m.state.renaming = true
			m.vms.input.SetValue(m.state.selectedContext)
			m.vms.input.Focus()
		case key.Matches(tmsg, keymap.Delete):
			m.state.deleting = true
			m.vms.input.Focus()
		case key.Matches(tmsg, keymap.Help):
			m.vms.help.ShowAll = !m.vms.help.ShowAll
		case key.Matches(tmsg, keymap.Quit):
			m.state.quitting = true
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m *model) finishAction() {
	m.vms.input.SetValue("")
	m.vms.input.Blur()
	m.state.deleting = false
	m.state.renaming = false
}
