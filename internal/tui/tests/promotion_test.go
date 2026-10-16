package tui_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui"
)

var whitePromotionSetup = []chess.Move{
	{From: chess.A2, To: chess.A4}, {From: chess.B8, To: chess.C6},
	{From: chess.A4, To: chess.A5}, {From: chess.H7, To: chess.H5},
	{From: chess.A5, To: chess.A6}, {From: chess.H5, To: chess.H4},
	{From: chess.A6, To: chess.B7}, {From: chess.H4, To: chess.H3},
}

var blackPromotionSetup = append(append([]chess.Move{}, whitePromotionSetup...),
	chess.Move{From: chess.B7, To: chess.B8, Promotion: chess.Knight},
	chess.Move{From: chess.H3, To: chess.G2}, chess.Move{From: chess.G1, To: chess.F3})

func pendingPromotion(t *testing.T, black, capture bool) (tea.Model, *game.Game, chess.Move) {
	t.Helper()
	setup, move := whitePromotionSetup, chess.Move{From: chess.B7, To: chess.B8}
	if black {
		setup, move = blackPromotionSetup, chess.Move{From: chess.G2, To: chess.G1}
	}
	if capture {
		move.To = chess.A8
		if black {
			move.To = chess.H1
		}
	}
	g := gameWith(t, setup...)
	before := g.Position()
	m, cmd := submit(t, tui.NewModel(g), move)
	if cmd != nil || !strings.Contains(m.View().Content, "Promotion:") {
		t.Fatal("Missing promotion choices")
	}
	assertUnchanged(t, g, before, len(setup)+1)
	return m, g, move
}

func TestPromotionKeyboardChoices(t *testing.T) {
	for _, tc := range []struct {
		black, capture bool
		key            rune
		kind           chess.PieceType
	}{
		{false, false, 'q', chess.Queen}, {false, true, 'r', chess.Rook},
		{true, false, 'b', chess.Bishop}, {true, true, 'n', chess.Knight},
	} {
		m, g, move := pendingPromotion(t, tc.black, tc.capture)
		color, length := g.Position().SideToMove, g.HistoryLen()
		m = updateModel(t, m, tea.KeyPressMsg{Code: tc.key})
		move.Promotion = tc.kind
		state, err := g.Snapshot()
		piece, _ := g.Position().Board.PieceAt(move.To)
		if err != nil || state.LastMove == nil || *state.LastMove != move || piece != (chess.Piece{Type: tc.kind, Color: color}) || g.HistoryLen() != length+1 {
			t.Fatalf("Incorrect UI promotion %v: %v", move, err)
		}
		if strings.Contains(m.View().Content, "Promotion:") {
			t.Error("Completed promotion still shown")
		}
		assertHighlights(t, m, chess.NoSquare)
	}
}

func TestPromotionCancellationAndInput(t *testing.T) {
	m, g, move := pendingPromotion(t, false, true)
	before, length, view := g.Position(), g.HistoryLen(), m
	for _, msg := range []tea.Msg{tea.KeyPressMsg{Code: 'h'}, tea.KeyPressMsg{Code: tea.KeyRight}, tea.KeyPressMsg{Code: tea.KeyEnter}, tea.KeyReleaseMsg{Code: 'q'}} {
		m = updateModel(t, m, msg)
	}
	assertSameView(t, m, view)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEscape})
	assertUnchanged(t, g, before, length)
	assertSameView(t, m, moveCursorTo(t, tui.NewModel(g), move.To))
}
