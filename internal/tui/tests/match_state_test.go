package tui_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui/match"
)

func TestStateDoesNotExposeMutableMoves(t *testing.T) {
	m := match.New(game.NewGame(), 1)
	next, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = next.(match.Model)
	state := m.State()
	if len(state.LegalMoves) == 0 {
		t.Fatal("selection did not provide moves")
	}
	want := state.LegalMoves[0]
	state.LegalMoves[0].To = chess.H8
	if m.State().LegalMoves[0] != want {
		t.Fatal("changing presentation data mutated legal moves")
	}

	pending, _, _ := pendingPromotion(t, false, true)
	promotion := pending.(match.Model)
	state = promotion.State()
	if len(state.PromotionChoices) == 0 {
		t.Fatal("promotion choices missing")
	}
	wantKind := state.PromotionChoices[0]
	state.PromotionChoices[0] = chess.NoPiece
	if promotion.State().PromotionChoices[0] != wantKind {
		t.Fatal("changing presentation data mutated promotion choices")
	}
}
