package tui

import (
	"github.com/ahmetb/kubectx/core/kubeconfig"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type MessageSender interface {
	Send(msg tea.Msg)
}

type modelState struct {
	currentContext    string
	currentNamespace  string
	selectedContext   string
	selectedNamespace string
	namespacesLoading bool
	focused           Component
	quitting          bool
	termWidth         int
	termHeight        int
}

type modelConfig struct {
	title  string
	keymap help.KeyMap
}

type viewModels struct {
	help          help.Model
	contextList   list.Model
	namespaceList list.Model
}

type model struct {
	kubeconf            *kubeconfig.Kubeconfig
	contexts            []string
	namespacesByContext map[string][]string

	state  modelState
	config modelConfig
	vms    viewModels
	sender MessageSender
}

type Component int

const (
	ContextList Component = iota
	NamespaceList
)

const ComponentCount = 2

type TablesListItem string

var _ list.DefaultItem = TablesListItem("")

func (t TablesListItem) FilterValue() string {
	return string(t)
}

func (t TablesListItem) Title() string {
	return string(t)
}

func (t TablesListItem) Description() string {
	return ""
}

type actionDoneMsg struct {
}

var _ tea.Msg = actionDoneMsg{}
