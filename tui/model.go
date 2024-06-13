package tui

import (
	"context"
	"facette.io/natsort"
	"fmt"
	"github.com/ahmetb/kubectx/core/kubeclient"
	"github.com/ahmetb/kubectx/core/kubeconfig"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/kubectui/tui/keymap"
	"github.com/reallyliri/kubectui/tui/list"
	"github.com/reallyliri/kubectui/tui/styles"
	"log"
)

var _ tea.Model = &model{}

func (m *model) Init() tea.Cmd {
	return nil
}

func Run(ctx context.Context, title string, kubeconf *kubeconfig.Kubeconfig) error {
	m, err := newModel(title, kubeconf)
	if err != nil {
		return err
	}
	program := tea.NewProgram(m)
	m.sender = program
	go func() {
		<-ctx.Done()
		program.Quit()
	}()
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}

func newModel(title string, kubeconf *kubeconfig.Kubeconfig) (*model, error) {
	contexts := kubeconf.ContextNames()
	natsort.Sort(contexts)

	helpVM := help.New()
	m := &model{
		kubeconf:            kubeconf,
		contexts:            contexts,
		namespacesByContext: make(map[string][]string, len(contexts)),
		state: modelState{
			currentContext: kubeconf.GetCurrentContext(),
			focused:        ContextList,
		},
		config: modelConfig{
			keymap: keymap.GetKeyMap(),
			title:  title,
		},
		vms: viewModels{
			help: helpVM,
		},
	}

	m.recreateContextList()
	m.onContextSelected(m.state.currentContext)

	return m, nil
}

func (m *model) onContextSelected(context string) {
	if m.state.selectedContext == context {
		return
	}

	m.state.selectedContext = context
	if context != "" {
		m.loadNamespaces(context)
	}
}

func (m *model) loadNamespaces(context string) {
	if _, ok := m.namespacesByContext[context]; ok {
		m.recreateNamespaceList(context)
		return
	}

	m.state.namespacesLoading = true
	go func() {
		m.doWithContext(context, func() {
			namespaces, err := kubeclient.QueryNamespaces(m.kubeconf)
			if err != nil {
				m.onError(fmt.Errorf("failed to query namespaces of %s: %w", context, err))
			}
			if namespaces != nil {
				natsort.Sort(namespaces)
			}
			m.namespacesByContext[context] = namespaces
		})
		m.state.namespacesLoading = false
		m.recreateNamespaceList(context)
		m.sender.Send(actionDoneMsg{})
	}()
}

func (m *model) recreateContextList() {
	m.vms.contextList = list.NewItemsList(m.contexts, "ctx", m.state.currentContext)
}

func (m *model) recreateNamespaceList(context string) {
	var err error
	m.state.currentNamespace, err = m.kubeconf.NamespaceOfContext(context)
	if err != nil {
		m.onError(fmt.Errorf("failed to get current namespace of %s: %w", context, err))
	}
	m.state.selectedNamespace = m.state.currentNamespace
	m.vms.namespaceList = list.NewItemsList(m.namespacesByContext[context], "ns", m.state.currentNamespace)
}

func (m *model) doWithContext(context string, action func()) {
	prevContext := m.kubeconf.GetCurrentContext()
	if err := m.kubeconf.ModifyCurrentContext(context); err != nil {
		m.onError(fmt.Errorf("failed to switch to context %s: %w", context, err))
	}
	defer func() {
		if err := m.kubeconf.ModifyCurrentContext(prevContext); err != nil {
			m.onError(fmt.Errorf("failed to switch back to context %s: %w", prevContext, err))
		}
	}()
	action()
}

func (m *model) onError(err error) {
	log.Fatalf(styles.ErrorStyle.Render("error: %v"), err)
}
