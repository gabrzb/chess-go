package tui

func (m *model) moveCursor(colStep, rowStep int) {
	next, ok := m.cursor.Offset(colStep, rowStep)
	if ok {
		m.cursor = next
	}
}
