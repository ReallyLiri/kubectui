package tui

import (
	"context"
	"fmt"
	"log"

	"facette.io/natsort"
	"github.com/ahmetb/kubectx/core/kubeclient"
	"github.com/ahmetb/kubectx/core/kubeconfig"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/reallyliri/kubectui/tui/keymap"
	"github.com/reallyliri/kubectui/tui/list"
	"github.com/reallyliri/kubectui/tui/styles"
	"github.com/reallyliri/syncmap"
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
		namespacesByContext: syncmap.SyncMap[string, []string]{},
		state: modelState{
			currentContext:    kubeconf.GetCurrentContext(),
			focused:           ContextList,
			namespacesLoading: syncmap.SyncMap[string, bool]{},
		},
		config: modelConfig{
			keymap: keymap.GetKeyMap(),
			title:  title,
		},
		vms: viewModels{
			help:  helpVM,
			input: textinput.New(),
		},
	}

	m.recreateContextList()
	m.onContextSelected(m.state.selectedContext)
	m.state.currentNamespace, _ = kubeconf.NamespaceOfContext(m.state.currentContext)

	return m, nil
}

func (m *model) onContextSelected(context string) {
	m.state.selectedContext = context
	if context != "" {
		m.loadNamespaces(context)
	}
}

func (m *model) onNamespaceSelected(namespace string) {
	m.state.selectedNamespace = namespace
}

func (m *model) loadNamespaces(context string) {
	if _, ok := m.namespacesByContext.Load(context); ok {
		m.recreateNamespaceList(context)
		return
	}

	m.state.namespacesLoading.Store(context, true)
	go func() {
		m.doWithContext(context, func() {
			namespaces, err := kubeclient.QueryNamespaces(m.kubeconf)
			if err != nil {
				m.onError(fmt.Errorf("failed to query namespaces of %s: %w", context, err))
			}
			if namespaces != nil {
				natsort.Sort(namespaces)
			}
			m.namespacesByContext.Store(context, namespaces)
		})
		m.recreateNamespaceList(context)
		m.state.namespacesLoading.Store(context, false)
		m.sender.Send(actionDoneMsg{})
	}()
}

func (m *model) recreateContextList() {
	m.vms.contextList = list.NewItemsList(m.contexts, "ctx", m.state.currentContext, styles.ContextTint)
	if m.state.selectedContext == "" {
		m.state.selectedContext = m.vms.contextList.SelectedItem().FilterValue()
	}
}

func (m *model) recreateNamespaceList(context string) {
	m.mut.Lock()
	defer m.mut.Unlock()
	if m.state.selectedContext != context {
		return
	}
	namespaceOfContext, err := m.kubeconf.NamespaceOfContext(context)
	if err != nil {
		m.onError(fmt.Errorf("failed to get current namespace of %s: %w", context, err))
	}
	m.state.selectedNamespace = namespaceOfContext
	m.state.currentNamespace = namespaceOfContext
	namespaces, _ := m.namespacesByContext.Load(context)
	m.vms.namespaceList = list.NewItemsList(namespaces, "ns", m.state.currentNamespace, styles.NamespaceTint)
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

func (m *model) setCurrentContextFromSelected() {
	if m.state.selectedContext == m.state.currentContext {
		return
	}
	if err := m.kubeconf.ModifyCurrentContext(m.state.selectedContext); err != nil {
		m.onError(fmt.Errorf("failed to switch to context %s: %w", m.state.selectedContext, err))
	}
	m.saveKubeconfig()
	m.state.currentContext = m.state.selectedContext
	var err error
	m.state.currentNamespace, err = m.kubeconf.NamespaceOfContext(m.state.currentContext)
	if err != nil {
		m.onError(fmt.Errorf("failed to get current namespace of %s: %w", m.state.currentContext, err))
	}
	m.recreateContextList()
}

func (m *model) setCurrentNamespaceFromSelected() {
	m.setCurrentContextFromSelected()
	if m.state.selectedNamespace == m.state.currentNamespace {
		return
	}
	if err := m.kubeconf.SetNamespace(m.state.currentContext, m.state.selectedNamespace); err != nil {
		m.onError(fmt.Errorf("failed to switch to namespace %s: %w", m.state.selectedNamespace, err))
	}
	m.saveKubeconfig()
	m.recreateNamespaceList(m.state.currentContext)
}

func (m *model) renameSelectedContext(newName string) {
	prevName := m.state.selectedContext
	if newName == prevName {
		return
	}
	if err := m.kubeconf.ModifyContextName(m.state.selectedContext, newName); err != nil {
		m.onError(fmt.Errorf("failed to rename context '%s' to '%s': %w", m.state.selectedContext, newName, err))
	}
	if m.state.currentContext == prevName {
		m.state.currentContext = newName
		if err := m.kubeconf.ModifyCurrentContext(newName); err != nil {
			m.onError(fmt.Errorf("failed to switch to renamed context '%s': %w", newName, err))
		}
	}
	m.saveKubeconfig()
	m.contexts = m.kubeconf.ContextNames()
	natsort.Sort(m.contexts)
	namespaces, _ := m.namespacesByContext.LoadAndDelete(prevName)
	m.namespacesByContext.Store(newName, namespaces)
	m.state.selectedContext = newName
	m.recreateContextList()
}

func (m *model) deleteSelectedContext() {
	selectedContext := m.state.selectedContext
	if err := m.kubeconf.DeleteContextEntry(selectedContext); err != nil {
		m.onError(fmt.Errorf("failed to delete context '%s': %w", selectedContext, err))
	}
	if m.state.currentContext == selectedContext {
		m.state.currentContext = ""
		if err := m.kubeconf.UnsetCurrentContext(); err != nil {
			m.onError(fmt.Errorf("failed to unset current context: %w", err))
		}
	}
	m.saveKubeconfig()
	m.contexts = m.kubeconf.ContextNames()
	m.namespacesByContext.Delete(selectedContext)
	natsort.Sort(m.contexts)
	m.state.selectedContext = ""
	m.recreateContextList()
}

func (m *model) saveKubeconfig() {
	if err := m.kubeconf.Save(); err != nil {
		m.onError(fmt.Errorf("failed to save kubeconfig: %w", err))
	}
}

func (m *model) refresh() {
	if err := m.kubeconf.Parse(); err != nil {
		m.onError(fmt.Errorf("failed to parse kubeconfig: %w", err))
	}
	contexts := m.kubeconf.ContextNames()
	natsort.Sort(contexts)
	m.contexts = contexts
	m.state.selectedContext = ""
	m.state.selectedNamespace = ""
	m.namespacesByContext = syncmap.SyncMap[string, []string]{}
	m.recreateContextList()
	m.onContextSelected(m.state.selectedContext)
}

func (m *model) onError(err error) {
	log.Fatalf(styles.ErrorStyle.Render("error: %v"), err)
}
