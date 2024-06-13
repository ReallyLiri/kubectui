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
		key.WithHelp("←/h", "move left"),
	)
	Right = key.NewBinding(
		key.WithKeys(tea.KeyRight.String(), "l"),
		key.WithHelp("→/l", "move right"),
	)
	Tab = key.NewBinding(
		key.WithKeys(tea.KeyTab.String(), "n"),
		key.WithHelp("Tab/n", "change focus"),
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
	return []key.Binding{Help, Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{Tab},
		{Up, Down},
		{Left, Right},
		{Help, Quit},
	}
}
