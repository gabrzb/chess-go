package match

import (
	"slices"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

type State struct {
	Cursor, Selected                      chess.Square
	Message                               string
	AIError                               error
	PromotionPending, Reviewing, Thinking bool
	HistoryIndex                          int
	AIColor                               chess.Color
	Flipped                               bool
	LegalMoves                            []chess.Move
	PromotionChoices                      []chess.PieceType
}

func (m Model) State() State {
	var promotions []chess.PieceType

	for _, move := range m.promotionMoves {
		promotions = append(promotions, move.Promotion)
	}

	return State{
		Cursor: m.cursor, Selected: m.selected, Message: m.message, AIError: m.aiError,
		PromotionPending: len(m.promotionMoves) > 0, Reviewing: m.reviewing, Thinking: m.opponent.Thinking(),
		HistoryIndex: m.historyIndex, AIColor: m.opponent.Color(), Flipped: m.flipped,
		LegalMoves: slices.Clone(m.legalMoves), PromotionChoices: promotions,
	}
}

func (m Model) Game() *game.Game { return m.game }

func (m *Model) ClearMessage() {
	m.message = ""
	m.aiError = nil
}

func (m *Model) EndReview() { m.reviewing = false }

func (m *Model) Flip() { m.flipped = !m.flipped }

func (m *Model) Promote(kind chess.PieceType) {
	if err := m.choosePromotion(kind); err != nil {
		m.message = err.Error()
	}
}
