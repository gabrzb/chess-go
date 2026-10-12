package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

type model struct {
	game   *game.Game
	cursor chess.Square
}

func NewModel(g *game.Game) tea.Model {
	return model{
		game:   g,
		cursor: chess.A2,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
