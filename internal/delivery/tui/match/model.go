package match

import (
	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/application"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

type Model struct {
	game           *game.Game
	cursor         chess.Square
	selected       chess.Square
	legalMoves     []chess.Move
	message        string
	aiError        error
	promotionMoves []chess.Move
	reviewing      bool
	historyIndex   int
	opponent       application.Opponent
	flipped        bool
}

func New(g *game.Game) Model {
	return Model{
		game:     g,
		cursor:   chess.A2,
		selected: chess.NoSquare,
	}
}

func (m Model) Init() tea.Cmd {
	if m.opponent.Enabled() {
		return func() tea.Msg { return startAIMessage{} }
	}

	return nil
}
