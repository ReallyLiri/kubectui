package keymap

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

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
		Help, Quit, Select,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Tab},
		{Up, Down},
		{Left, Right},
		{Select},
		{Help, Quit},
	}
}
