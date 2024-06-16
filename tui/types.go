package tui

import (
	"github.com/ahmetb/kubectx/core/kubeconfig"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"sync"
)

type MessageSender interface {
	Send(msg tea.Msg)
}

type modelState struct {
	currentContext    string
	currentNamespace  string
	selectedContext   string
	selectedNamespace string
	namespacesLoading map[string]bool
	focused           Component
	renaming          bool
	deleting          bool
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
	input         textinput.Model
}

type model struct {
	mut                 sync.Mutex
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

type actionDoneMsg struct {
}

var _ tea.Msg = actionDoneMsg{}
