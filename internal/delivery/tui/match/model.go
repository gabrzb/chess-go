package match

import (
	"context"

	tea "charm.land/bubbletea/v2"

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
	opponent       chess.MoveChooser
	aiColor        chess.Color
	aiContext      context.Context
	cancelAI       context.CancelFunc
	thinking       bool
	requestID      uint64
	sessionID      uint64
	flipped        bool
}

func New(g *game.Game, session uint64) Model {
	return Model{
		game:      g,
		sessionID: session,
		cursor:    chess.A2,
		selected:  chess.NoSquare,
	}
}

func (m Model) Init() tea.Cmd {
	if m.opponent != nil {
		return func() tea.Msg { return startAIMessage{} }
	}

	return nil
}
