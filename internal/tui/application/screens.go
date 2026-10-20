package application

import (
	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui/match"
)

func (m application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if tick, ok := msg.(animationMessage); ok {
		if tick.id != m.animationID || !m.animating() {
			return m, nil
		}

		m.frame++

		if m.screen == mainScreen && m.frame >= introFrames {
			m.introDone = true
		}

		return m, m.animationTick()
	}

	next, cmd := m.update(msg)
	n := next.(application)

	if n.animating() && (n.animationID != m.animationID || !m.animating()) {
		n.frame = 0
		n.animationID++
		cmd = combineCommands(cmd, n.animationTick())
	}

	return n, cmd
}

func (m application) update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = max(1, size.Width), max(1, size.Height)
		return m, nil
	}

	if match.IsAIResult(msg, m.session) {
		if m.match.Game() == nil {
			return m, nil
		}

		next, cmd := m.match.Update(msg)
		m.match = next.(match.Model)
		m.afterMatch()

		return m, cmd
	}

	key, ok := msg.(tea.KeyPressMsg)

	if !ok {
		return m, nil
	}

	if key.String() == "ctrl+c" {
		m.match.StopAI()
		return m, tea.Quit
	}

	if m.screen == mainScreen {
		if key.String() == "a" {
			m.noAnimations = !m.noAnimations
			m.introDone = true
			m.animationID++

			return m, nil
		}
		if key.String() == "enter" && !m.noAnimations && !m.introDone {
			m.introDone = true
			m.animationID++

			return m, nil
		}
	}
	if key.String() == "?" && m.screen != helpScreen {
		m.helpBack, m.helpChoice = m.screen, m.choice
		m.show(helpScreen)

		return m, nil
	}

	if m.screen == playingScreen || m.screen == finalBoardScreen {
		return m.updateBoard(key)
	}

	if key.String() == "esc" || key.String() == "q" {
		return m.back()
	}

	options := m.menuOptions()
	if len(options) == 0 {
		return m, nil
	}

	switch key.String() {
	case "up", "w":
		m.choice = cycleIndex(m.choice, -1, len(options))
	case "down", "s":
		m.choice = cycleIndex(m.choice, 1, len(options))
	case "left", "a", "right", "d":
		if m.screen == configScreen {
			step := 1

			if key.String() == "left" || key.String() == "a" {
				step = -1
			}

			m.adjustConfig(step)
		}
	case "enter":
		cmd := m.activate()
		return m, cmd
	}
	return m, nil
}

func (m application) back() (tea.Model, tea.Cmd) {
	switch m.screen {
	case mainScreen:
		return m, tea.Quit
	case configScreen:
		m.show(mainScreen)
	case pauseScreen:
		cmd := m.resume()
		return m, cmd
	case resignScreen, abandonScreen:
		m.show(pauseScreen)
	case drawScreen:
		s, _ := m.match.Game().Snapshot()

		if err := m.match.Game().RespondDraw(opposite(s.DrawOffer), false); err != nil {
			m.logger.Print(err)
		}

		cmd := m.resume()
		return m, cmd
	case helpScreen:
		m.show(m.helpBack)
		m.choice = m.helpChoice

		if m.screen == playingScreen {
			cmd := m.match.StartAI()
			return m, cmd
		}
	case resultScreen:
		m.show(finalBoardScreen)
	case errorScreen:
		m.show(abandonScreen)
	}
	return m, nil
}

func (m application) updateBoard(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.match.State().PromotionPending && m.screen == playingScreen {
		switch key.String() {
		case "up", "w":
			m.choice = cycleIndex(m.choice, -1, len(m.match.State().PromotionChoices))
		case "down", "s":
			m.choice = cycleIndex(m.choice, 1, len(m.match.State().PromotionChoices))
		case "enter":
			kind := m.match.State().PromotionChoices[m.choice]
			m.match.Promote(kind)
		default:
			m.match.PromotionKey(key.String())
		}
		m.afterMatch()
		if m.screen == playingScreen {
			cmd := m.match.StartAI()
			return m, cmd
		}
		return m, nil
	}

	if key.String() == "f" {
		m.match.Flip()
		return m, nil
	}

	if key.String() == "esc" && !m.match.State().Reviewing {
		if m.screen == finalBoardScreen {
			m.show(resultScreen)
		} else if m.match.State().Selected.Valid() || m.match.State().Message != "" {
			m.match.ClearSelection()
		} else {
			m.show(pauseScreen)
		}
		return m, nil
	}

	if key.String() == "q" {
		if m.screen == finalBoardScreen {
			m.show(resultScreen)
		} else {
			m.show(pauseScreen)
		}

		return m, nil
	}

	if m.screen == finalBoardScreen && key.String() != "h" && !m.match.State().Reviewing {
		return m, nil
	}

	if !m.boardFits() && key.String() != "h" && key.String() != "esc" {
		return m, nil
	}

	wasPromotion := m.match.State().PromotionPending
	next, cmd := m.match.Update(key)
	m.match = next.(match.Model)

	if !wasPromotion && m.match.State().PromotionPending {
		m.choice = 0
	}

	m.afterMatch()
	return m, cmd
}

func (m *application) afterMatch() {
	if err := m.match.State().AIError; err != nil {
		m.logger.Print(err)
		m.match.ClearMessage()
		m.show(errorScreen)
		return
	}

	snapshot, err := m.match.Game().Snapshot()
	if err != nil {
		m.logger.Print(err)
		m.show(errorScreen)
		return
	}

	m.refreshLabels()
	if m.config.computer && snapshot.SideToMove == m.match.State().AIColor && m.match.Game().ClaimableDraw() != game.NoReason {
		if err := m.match.Game().ClaimDraw(); err != nil {
			m.logger.Print(err)
			return
		}
		snapshot, _ = m.match.Game().Snapshot()
	}

	if snapshot.Status != chess.StatusOngoing && m.screen == playingScreen {
		m.show(resultScreen)
	}
}
