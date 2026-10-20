package match

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func TestStateDoesNotExposeMutableMoves(t *testing.T) {
	m := New(game.NewGame(), 1)
	move := chess.Move{From: chess.A2, To: chess.A4}
	m.legalMoves = []chess.Move{move}
	m.promotionMoves = []chess.Move{{From: chess.A7, To: chess.A8, Promotion: chess.Queen}}
	state := m.State()
	state.LegalMoves[0].To = chess.H8
	state.PromotionChoices[0] = chess.Knight
	if m.legalMoves[0] != move || m.promotionMoves[0].Promotion != chess.Queen {
		t.Fatal("changing presentation data mutated the match")
	}
}
