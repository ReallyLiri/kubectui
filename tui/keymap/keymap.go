package keymap

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// hjkl is vim standard for movement

var (
	Up = key.NewBinding(
		key.WithKeys(tea.KeyUp.String(), "k"),
		key.WithHelp("↑/k", "move up"),
	)
	Down = key.NewBinding(
		key.WithKeys(tea.KeyDown.String(), "j"),
		key.WithHelp("↓/j", "move down"),
	)
	Left = key.NewBinding(
		key.WithKeys(tea.KeyLeft.String(), "h"),
		key.WithHelp("←/h", "switch focus left"),
	)
	Right = key.NewBinding(
		key.WithKeys(tea.KeyRight.String(), "l"),
		key.WithHelp("→/l", "switch focus right"),
	)
	Tab = key.NewBinding(
		key.WithKeys(tea.KeyTab.String(), "n"),
		key.WithHelp("Tab/n", "next focus"),
	)
	Select = key.NewBinding(
		key.WithKeys(tea.KeyEnter.String(), "s"),
		key.WithHelp("Enter/s", "select current ctx or ns"),
	)
	Rename = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename ctx"),
	)
	Refresh = key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "refresh"),
	)
	Delete = key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete ctx"),
	)
	Cancel = key.NewBinding(
		key.WithKeys(tea.KeyEscape.String()),
		key.WithHelp(tea.KeyEscape.String(), "cancel"),
	)
	Search = key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	)
	Help = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	)
	Quit = key.NewBinding(
		key.WithKeys("q", tea.KeyEscape.String(), tea.KeyCtrlC.String()),
		key.WithHelp("q", "quit"),
	)
)

type keyMap struct {
}

var _ help.KeyMap = keyMap{}

func GetKeyMap() help.KeyMap {
	return keyMap{}
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		Help, Quit, Select, Rename, Delete,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Tab},
		{Left, Up, Down, Right},
		{Rename, Delete, Select},
		{Refresh, Help, Quit},
	}
}
