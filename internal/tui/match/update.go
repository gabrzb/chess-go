package match

import tea "charm.land/bubbletea/v2"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startAIMessage:
		cmd := m.StartAI()
		return m, cmd
	case aiResultMessage:
		return m.receiveAIResult(msg)
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		if key.String() == "ctrl+c" {
			m.StopAI()
			return m, tea.Quit
		}

		if len(m.promotionMoves) > 0 {
			m.PromotionKey(key.String())
			cmd := m.StartAI()
			return m, cmd
		}

		if key.String() == "h" {
			m.toggleReview()
			if m.reviewing {
				m.StopAI()
				return m, nil
			}
			cmd := m.StartAI()
			return m, cmd
		}

		if m.reviewing {
			switch key.String() {
			case "left", "a":
				m.moveHistory(-1)
			case "right", "d":
				m.moveHistory(1)
			case "q":
				m.StopAI()
				return m, tea.Quit
			case "esc":
				m.toggleReview()
				cmd := m.StartAI()
				return m, cmd
			}

			return m, nil
		}

		switch key.String() {
		case "q":
			m.StopAI()
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
				cmd := m.StartAI()
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
			cmd := m.StartAI()
			return m, cmd
		case "esc":
			m.ClearSelection()
		}
	}

	return m, nil
}
