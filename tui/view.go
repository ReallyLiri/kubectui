package tui

import (
	_ "embed"
	"github.com/charmbracelet/lipgloss"
	"github.com/reallyliri/kubectui/tui/format"
	"github.com/reallyliri/kubectui/tui/styles"
)

func (m *model) View() string {
	if m.state.quitting || m.state.termWidth == 0 || m.state.termHeight == 0 {
		return ""
	}

	borderWidth, borderHeight := styles.BorderFocusedStyle.GetFrameSize()

	title := titleView(m.config.title, m.state.currentContext, m.state.currentNamespace)

	if m.state.deleting {
		content := lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				"Are you sure you want to ",
				styles.DangerStyle.Render("delete"),
			),
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				"the context ",
				lipgloss.NewStyle().Foreground(styles.GreenTint).Render(m.state.selectedContext),
				styles.SubTitleStyle.Render(" (y/n or esc to cancel)"),
			),
			m.vms.input.View(),
		)
		return lipgloss.JoinVertical(
			lipgloss.Top,
			title,
			" ",
			content,
		)
	}

	if m.state.renaming {
		content := lipgloss.JoinVertical(
			lipgloss.Top,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				"Enter the new name for the context ",
				lipgloss.NewStyle().Foreground(styles.GreenTint).Render(m.state.selectedContext),
				styles.SubTitleStyle.Render(" (enter to apply or esc to cancel)"),
			),
			m.vms.input.View(),
		)
		return lipgloss.JoinVertical(
			lipgloss.Top,
			title,
			" ",
			content,
		)
	}

	footer := m.vms.help.View(m.config.keymap)
	centerHeight := m.state.termHeight - lipgloss.Height(title) - lipgloss.Height(footer) - 5

	var contextsList, namespacesList string

	mainWidth := m.state.termWidth/ComponentCount - borderWidth - 3
	mainHeight := centerHeight - borderHeight + 3

	if len(m.contexts) == 0 {
		contextsList = m.emptyMessage(ContextList, mainWidth, mainHeight)
	} else {
		m.vms.contextList.SetSize(mainWidth, mainHeight)
		contextsList = withBorder(m.vms.contextList.View(), m.state.focused == ContextList)
		if m.state.selectedContext != "" {
			if m.state.namespacesLoading || len(m.namespacesByContext[m.state.selectedContext]) == 0 {
				namespacesList = m.emptyMessage(NamespaceList, mainWidth, mainHeight)
			} else {
				m.vms.namespaceList.SetSize(mainWidth, mainHeight)
				namespacesList = withBorder(m.vms.namespaceList.View(), m.state.focused == NamespaceList)
			}
		}
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			contextsList,
			" ",
			namespacesList,
		),
		footer,
	)
}

func titleView(title, context, namespace string) string {
	parts := make([]string, 0, 8)
	parts = append(parts, styles.TitleStyle.Render(title))
	parts = append(parts, " ")
	if context != "" {
		parts = append(parts, styles.BreadcrumbsSectionStyle.Render(" ", format.BreadcrumbsSeparator, " ctx: "))
		parts = append(parts, styles.BreadcrumbsTitleStyle.Render(context))
		if namespace != "" {
			parts = append(parts, styles.BreadcrumbsSectionStyle.Render(" ", format.BreadcrumbsSeparator, " ns: "))
			parts = append(parts, styles.BreadcrumbsTitleStyle.Render(namespace))
		} else {
			parts = append(parts, styles.BreadcrumbsSectionStyle.Render(" ", format.BreadcrumbsSeparator, " No current ns"))
		}
	} else {
		parts = append(parts, styles.BreadcrumbsSectionStyle.Render(" ", format.BreadcrumbsSeparator, " No current ctx"))
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, parts...)
}

func withBorder(ui string, focused bool) string {
	if focused {
		return styles.BorderFocusedStyle.Render(ui)
	} else {
		return styles.BorderBlurredStyle.Render(ui)
	}
}

func (m *model) emptyMessage(component Component, width, height int) string {
	var message string
	switch component {
	case ContextList:
		message = "No contexts"
	case NamespaceList:
		if m.state.namespacesLoading {
			message = "Loading namespaces..."
		} else {
			message = "No namespaces"
		}
	}
	return withBorder(
		styles.NoDataStyle.
			Width(width).
			Height(height).
			Render(message),
		m.state.focused == component,
	)
}
