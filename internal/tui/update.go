package tui

import tea "charm.land/bubbletea/v2"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		if key.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if len(m.promotionMoves) > 0 {
			m.handlePromotionKey(key.String())
			return m, nil
		}

		if key.String() == "h" {
			m.toggleReview()

			return m, nil
		}

		if m.reviewing {
			switch key.String() {
			case "left", "a":
				m.moveHistory(-1)
			case "right", "d":
				m.moveHistory(1)
			case "q":
				return m, tea.Quit
			case "esc":
				m.toggleReview()
			}

			return m, nil
		}

		switch key.String() {
		case "q":
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
		case "esc":
			m.clearSelection()
		}
	}

	return m, nil
}
