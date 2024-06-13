package styles

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

const (
	BlueTint          = lipgloss.Color("#388de9")
	GreenTint         = lipgloss.Color("#46b17b")
	RedTint           = lipgloss.Color("#ff0000")
	BorderFocusedTint = lipgloss.Color("63")
	BorderBlurredTint = lipgloss.Color("240")
)

var (
	SubTitleTint = help.New().Styles.ShortKey.GetForeground()

	TitleStyle              = lipgloss.NewStyle().Foreground(BlueTint)
	SubTitleStyle           = lipgloss.NewStyle().Foreground(SubTitleTint)
	BreadcrumbsSectionStyle = SubTitleStyle
	BreadcrumbsTitleStyle   = lipgloss.NewStyle().Foreground(GreenTint)
	BorderFocusedStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(BorderFocusedTint).Padding(0, 1)
	BorderBlurredStyle      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(BorderBlurredTint).Padding(0, 1)
	NoDataStyle             = lipgloss.NewStyle().Foreground(SubTitleTint).AlignHorizontal(lipgloss.Left).Padding(2)
	ErrorStyle              = lipgloss.NewStyle().Foreground(RedTint)
	DangerStyle             = lipgloss.NewStyle().Foreground(RedTint)
)
