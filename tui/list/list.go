package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/reallyliri/kubectui/tui/styles"
	"github.com/samber/lo"
	"strings"
)

const selectedIndicator = "[*]"

type tablesListItem struct {
	name     string
	selected bool
}

var _ list.DefaultItem = tablesListItem{}

func (t tablesListItem) FilterValue() string {
	return t.name
}

func (t tablesListItem) Title() string {
	if t.selected {
		return fmt.Sprintf("%s %s", selectedIndicator, t.name)
	}
	return fmt.Sprintf("%s %s", strings.Repeat(" ", len(selectedIndicator)), t.name)
}

func (t tablesListItem) Description() string {
	return ""
}

func NewItemsList(items []string, componentName string, selectedItem string, selectedItemStyle lipgloss.TerminalColor) list.Model {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(selectedItemStyle).BorderForeground(selectedItemStyle)
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	var selectedIndex int
	l := list.New(
		lo.Map(items, func(name string, i int) list.Item {
			selected := selectedItem == name
			if selected {
				selectedIndex = i
			}
			return tablesListItem{
				name:     name,
				selected: selected,
			}
		}),
		delegate,
		0,
		0,
	)
	l.SetFilteringEnabled(true)
	l.Styles.ActivePaginationDot = styles.TitleStyle
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetStatusBarItemName(componentName, componentName)
	l.SetShowPagination(false)
	l.Select(selectedIndex)
	return l
}
