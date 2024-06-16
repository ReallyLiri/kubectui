package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

const (
	blueTint    = lipgloss.Color("#6f9ff8")
	greenTint   = lipgloss.Color("#6ff889")
	purpleTint  = lipgloss.Color("#ac6ff8")
	pinkTint    = lipgloss.Color("#EE6FF8")
	redTint     = lipgloss.Color("#d31919")
	blurredTint = lipgloss.Color("240")
)

var (
	subTitleTint  = help.New().Styles.ShortKey.GetForeground()
	ContextTint   = purpleTint
	NamespaceTint = pinkTint

	TextStyle                   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#0c0c0c", Dark: "#f0f0f0"})
	TitleStyle                  = lipgloss.NewStyle().Foreground(blueTint)
	SubTitleStyle               = lipgloss.NewStyle().Foreground(subTitleTint)
	BreadcrumbsSectionStyle     = SubTitleStyle
	BreadcrumbsTitleStyle       = lipgloss.NewStyle().Foreground(greenTint)
	BorderFocusedContextStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(ContextTint).Padding(0, 1)
	BorderFocusedNamespaceStyle = BorderFocusedContextStyle.BorderForeground(NamespaceTint)
	BorderBlurredStyle          = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(blurredTint).Padding(0, 1)
	NoDataStyle                 = lipgloss.NewStyle().Foreground(subTitleTint).AlignHorizontal(lipgloss.Left).Padding(2)
	ErrorStyle                  = lipgloss.NewStyle().Foreground(redTint)
	DangerStyle                 = lipgloss.NewStyle().Foreground(redTint)
)
