package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess/game"
)

type model struct {
	game *game.Game
}

func newModel() model {
	return model{
		game: game.NewGame(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
