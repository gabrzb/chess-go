package tui_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/tui"
)

func TestSelectionMoveAndErrorRecovery(t *testing.T) {
	g := gameWith(t)
	m := tui.NewModel(g)
	for _, move := range []chess.Move{{From: chess.A2, To: chess.A4}, {From: chess.A7, To: chess.A5}} {
		before, length := g.Position(), g.HistoryLen()
		m = moveCursorTo(t, m, move.From)
		m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
		m = moveCursorTo(t, m, chess.D4)
		m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
		middle, _ := move.From.Offset(0, (move.To.Row()-move.From.Row())/2)
		assertHighlights(t, m, move.From, middle, move.To)
		if !strings.Contains(m.View().Content, chess.ErrIllegalMove.Error()) {
			t.Fatal("Expected invalid-destination message")
		}
		assertUnchanged(t, g, before, length)
		m = moveCursorTo(t, m, move.To)
		m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
		state, err := g.Snapshot()
		if err != nil || state.LastMove == nil || *state.LastMove != move || g.HistoryLen() != length+1 {
			t.Fatalf("Move %v was not recorded once: %v", move, err)
		}
		if strings.Contains(m.View().Content, chess.ErrIllegalMove.Error()) {
			t.Error("Successful move kept the error")
		}
		assertHighlights(t, m, chess.NoSquare)
	}
	m = moveCursorTo(t, m, chess.B2)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEscape})
	assertHighlights(t, m, chess.NoSquare)
}

func TestStatusPresentation(t *testing.T) {
	for _, tc := range []struct {
		moves []chess.Move
		text  string
	}{
		{[]chess.Move{{From: chess.E2, To: chess.E4}, {From: chess.D7, To: chess.D5}, {From: chess.F1, To: chess.B5}}, "Turn: Black (check)"},
		{[]chess.Move{{From: chess.F2, To: chess.F3}, {From: chess.E7, To: chess.E5}, {From: chess.G2, To: chess.G4}, {From: chess.D8, To: chess.H4}}, "Checkmate! Black wins."},
		{[]chess.Move{{From: chess.E2, To: chess.E4}, {From: chess.F7, To: chess.F6}, {From: chess.D2, To: chess.D4}, {From: chess.G7, To: chess.G5}, {From: chess.D1, To: chess.H5}}, "Checkmate! White wins."},
		{[]chess.Move{
			{From: chess.E2, To: chess.E3}, {From: chess.A7, To: chess.A5},
			{From: chess.D1, To: chess.H5}, {From: chess.A8, To: chess.A6},
			{From: chess.H5, To: chess.A5}, {From: chess.H7, To: chess.H5},
			{From: chess.A5, To: chess.C7}, {From: chess.A6, To: chess.H6},
			{From: chess.H2, To: chess.H4}, {From: chess.F7, To: chess.F6},
			{From: chess.C7, To: chess.D7}, {From: chess.E8, To: chess.F7},
			{From: chess.D7, To: chess.B7}, {From: chess.D8, To: chess.D3},
			{From: chess.B7, To: chess.B8}, {From: chess.D3, To: chess.H7},
			{From: chess.B8, To: chess.C8}, {From: chess.F7, To: chess.G6},
			{From: chess.C8, To: chess.E6},
		}, "Stalemate! Draw."},
	} {
		g := gameWith(t, tc.moves...)
		m := tui.NewModel(g)
		if !strings.Contains(m.View().Content, tc.text) {
			t.Errorf("Missing status %q", tc.text)
		}
	}
}
