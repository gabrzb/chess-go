package application

import (
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/delivery/tui/render"
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
		(m.screen == playingScreen && m.match.State().Thinking))
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
		length = width * m.frame / barFrames
	}

	return render.Band(width, length, double)
}

func (m application) logoText() string {
	if m.noAnimations || m.introDone || m.frame < barFrames {
		return render.Paint(render.Logo, render.Accent)
	}

	return render.LogoHighlight(m.frame-barFrames, glowFrames)
}
