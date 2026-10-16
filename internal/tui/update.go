package tui

import tea "charm.land/bubbletea/v2"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startAIMessage:
		cmd := m.startAITurn()
		return m, cmd
	case aiResultMessage:
		return m.receiveAIResult(msg)
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		if key.String() == "ctrl+c" {
			m.stopAITurn()
			return m, tea.Quit
		}

		if len(m.promotionMoves) > 0 {
			m.handlePromotionKey(key.String())
			cmd := m.startAITurn()
			return m, cmd
		}

		if key.String() == "h" {
			m.toggleReview()
			if m.reviewing {
				m.stopAITurn()
				return m, nil
			}
			cmd := m.startAITurn()
			return m, cmd
		}

		if m.reviewing {
			switch key.String() {
			case "left", "a":
				m.moveHistory(-1)
			case "right", "d":
				m.moveHistory(1)
			case "q":
				m.stopAITurn()
				return m, tea.Quit
			case "esc":
				m.toggleReview()
				cmd := m.startAITurn()
				return m, cmd
			}

			return m, nil
		}

		switch key.String() {
		case "q":
			m.stopAITurn()
			return m, tea.Quit
		case "up", "w":
			m.moveCursor(0, 1)
		case "down", "s":
			m.moveCursor(0, -1)
		case "left", "a":
			m.moveCursor(-1, 0)
		case "right", "d":
			m.moveCursor(1, 0)
		case "enter":
			if m.isAITurn() {
				cmd := m.startAITurn()
				return m, cmd
			}
			var err error

			if m.selected.Valid() {
				err = m.playSelectedMove()
			} else {
				err = m.selectPiece()
			}

			if err != nil {
				m.message = err.Error()
			} else {
				m.message = ""
			}
			cmd := m.startAITurn()
			return m, cmd
		case "esc":
			m.clearSelection()
		}
	}

	return m, nil
}
