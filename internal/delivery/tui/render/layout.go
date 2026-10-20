package render

import "charm.land/lipgloss/v2"

const (
	DefaultWidth      = 80
	DefaultHeight     = 24
	MinPanelWidth     = 18
	WindowInset       = 4
	MenuWidth         = 52
	PromotionWidth    = 48
	MatchWidth        = 76
	WideMatchMinWidth = 72
	MaxHistoryRows    = 4
	ColumnGap         = 2
	LogoMinWidth      = 48
	LogoMinHeight     = 23
)

func ContentWidth(windowWidth, preferred int) int {
	return min(preferred, max(MinPanelWidth, windowWidth-WindowInset))
}

func Fits(content string, width, height int) bool {
	return lipgloss.Width(content) <= width && lipgloss.Height(content) <= height
}
