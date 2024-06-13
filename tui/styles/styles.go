package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

const (
	WhiteTint         = lipgloss.Color("#f0f0f0")
	BlueTint          = lipgloss.Color("#388de9")
	GreenTint         = lipgloss.Color("#46b17b")
	BorderFocusedTint = lipgloss.Color("63")
	BorderBluredTint  = lipgloss.Color("240")
)

var (
	SubTitleTint = help.New().Styles.ShortKey.GetForeground()

	TitleStyle              = lipgloss.NewStyle().Foreground(BlueTint)
	SubTitleStyle           = lipgloss.NewStyle().Foreground(SubTitleTint)
	BreadcrumbsSectionStyle = SubTitleStyle
	BreadcrumbsTitleStyle   = lipgloss.NewStyle().Foreground(GreenTint)
	BorderFocusedStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(BorderFocusedTint).Padding(0, 1)
	BorderBluredStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(BorderBluredTint).Padding(0, 1)
	NoDataStyle             = lipgloss.NewStyle().Foreground(SubTitleTint).AlignHorizontal(lipgloss.Left).Padding(2)
)
