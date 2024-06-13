package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/samber/lo"
)

type tablesListItem string

var _ list.DefaultItem = tablesListItem("")

func (t tablesListItem) FilterValue() string {
	return string(t)
}

func (t tablesListItem) Title() string {
	return string(t)
}

func (t tablesListItem) Description() string {
	return ""
}

func NewItemsList(items []string, componentName string) list.Model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	l := list.New(
		lo.Map(items, func(name string, _ int) list.Item {
			return tablesListItem(name)
		}),
		delegate,
		0,
		0,
	)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetStatusBarItemName(componentName, componentName)
	l.SetShowPagination(false)
	return l
}
