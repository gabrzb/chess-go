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

func newModel() model {
	return model{
		game:   game.NewGame(),
		cursor: chess.A2,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
