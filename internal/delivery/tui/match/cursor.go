package match

func (m *Model) moveCursor(colStep, rowStep int) {
	if m.flipped {
		colStep, rowStep = -colStep, -rowStep
	}

	next, ok := m.cursor.Offset(colStep, rowStep)

	if ok {
		m.cursor = next
	}
}
