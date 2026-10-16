package match

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui/render"
)

type State struct {
	Cursor, Selected                      chess.Square
	Message                               string
	PromotionPending, Reviewing, Thinking bool
	HistoryIndex                          int
	AIColor                               chess.Color
}

func (m Model) State() State {
	return State{Cursor: m.cursor, Selected: m.selected, Message: m.message,
		PromotionPending: len(m.promotionMoves) > 0, Reviewing: m.reviewing, Thinking: m.thinking,
		HistoryIndex: m.historyIndex, AIColor: m.aiColor}
}

func (m Model) Game() *game.Game { return m.game }

func (m *Model) ClearMessage() { m.message = "" }

func (m *Model) EndReview() { m.reviewing = false }

func (m *Model) Flip() { m.flipped = !m.flipped }

func (m *Model) Promote(kind chess.PieceType) {
	if err := m.choosePromotion(kind); err != nil {
		m.message = err.Error()
	}
}

func (m Model) RenderBoard(board chess.Board) string {
	return render.Board(board, render.BoardState{Cursor: m.cursor, Selected: m.selected, LegalMoves: m.legalMoves, Flipped: m.flipped})
}
