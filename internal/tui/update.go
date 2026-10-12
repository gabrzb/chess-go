package tui

import tea "charm.land/bubbletea/v2"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "q", "ctrl+c":
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
