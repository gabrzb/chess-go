package tui

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	barFrames     = 5
	glowFrames    = 8
	introFrames   = barFrames + glowFrames
	frameDuration = 60 * time.Millisecond
)

type animationMessage struct{ id uint64 }

func combineCommands(first, second tea.Cmd) tea.Cmd {
	if first == nil {
		return second
	}
	if second == nil {
		return first
	}
	return tea.Batch(first, second)
}

func (m application) animating() bool {
	return !m.noAnimations && ((m.screen == mainScreen && !m.introDone) ||
		(m.screen == resultScreen && m.frame < barFrames) ||
		(m.screen == playingScreen && m.match.thinking))
}

func (m application) animationTick() tea.Cmd {
	if !m.animating() {
		return nil
	}
	id := m.animationID
	return tea.Tick(frameDuration, func(time.Time) tea.Msg { return animationMessage{id: id} })
}

func (m application) band(width int, double bool) string {
	length := width
	if !m.noAnimations && ((m.screen == mainScreen && !m.introDone) || m.screen == resultScreen) && m.frame < barFrames {
		length = max(2, width*m.frame/barFrames)
	}
	bar := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(paint(strings.Repeat("=", length), muted))
	if double {
		bar += "\n" + bar
	}
	return bar
}

func (m application) logoText() string {
	if m.noAnimations || m.introDone || m.frame < barFrames {
		return paint(logo, accent)
	}
	column := (m.frame - barFrames) * lipgloss.Width(logo) / glowFrames
	lines := strings.Split(logo, "\n")
	for row, line := range lines {
		var text strings.Builder
		for index, char := range []rune(line) {
			color := accent
			if index >= column && index < column+3 {
				color = foreground
			}
			text.WriteString(paint(string(char), color))
		}
		lines[row] = text.String()
	}
	return strings.Join(lines, "\n")
}

func (m application) thinkingText() string {
	if m.noAnimations {
		return "Calculando jogada…"
	}
	frames := []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")
	return string(frames[m.frame%len(frames)]) + " Calculando jogada…"
}
