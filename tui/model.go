package tui

import (
	"context"
	"facette.io/natsort"
	"fmt"
	"github.com/ahmetb/kubectx/core/kubeclient"
	"github.com/ahmetb/kubectx/core/kubeconfig"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/kubectui/tui/keymap"
	"github.com/samber/lo"
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
			contextList: newItemsList(contexts, ContextList),
			help:        help.New(),
		},
	}

	m.onContextSelected(m.state.currentContext)

	return m, nil
}

func (m *model) onContextSelected(context string) {
	if m.state.selectedContext == context {
		return
	}

	m.state.selectedContext = context
	if context != "" {
		if _, ok := m.namespacesByContext[context]; !ok {
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
		}
		var err error
		m.state.currentNamespace, err = m.kubeconf.NamespaceOfContext(context)
		if err != nil {
			m.onError(fmt.Errorf("failed to get current namespace of %s: %w", context, err))
		}
		m.state.selectedNamespace = m.state.currentNamespace
	}
	m.vms.namespaceList = newItemsList(m.namespacesByContext[context], NamespaceList)
}

func newItemsList(items []string, component Component) list.Model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)
	lst := list.New(
		lo.Map(items, func(name string, _ int) list.Item {
			return TablesListItem(name)
		}),
		delegate,
		0,
		0,
	)
	lst.SetFilteringEnabled(false)
	lst.SetShowHelp(false)
	lst.SetShowTitle(false)
	switch component {
	case ContextList:
		lst.SetStatusBarItemName("ctx", "ctx")
	case NamespaceList:
		lst.SetStatusBarItemName("ns", "ns")
	}
	lst.SetShowPagination(false)
	return lst
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
	log.Fatalf("error: %v", err)
}
