package tui

import "github.com/gabrzb/chess-go/internal/chess"

func (m *model) playSelectedMove() error {
	move := chess.Move{
		From: m.selected,
		To:   m.cursor,
	}

	if _, err := m.game.Play(move); err != nil {
		return err
	}

	m.clearSelection()
	return nil
}
