package tui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

type model struct {
	game           *game.Game
	cursor         chess.Square
	selected       chess.Square
	legalMoves     []chess.Move
	message        string
	promotionMoves []chess.Move
	reviewing      bool
	historyIndex   int
	opponent       ai.MoveChooser
	aiColor        chess.Color
	aiContext      context.Context
	cancelAI       context.CancelFunc
	thinking       bool
	requestID      uint64
}

func NewModel(g *game.Game) tea.Model {
	return model{
		game:     g,
		cursor:   chess.A2,
		selected: chess.NoSquare,
	}
}

func (m model) Init() tea.Cmd {
	if m.opponent != nil {
		return func() tea.Msg { return startAIMessage{} }
	}
	return nil
}
