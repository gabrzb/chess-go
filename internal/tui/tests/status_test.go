package tui_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui"
)

func TestViewGameStatus(t *testing.T) {
	cases := []struct {
		name   string
		moves  []chess.Move
		status chess.Status
		check  bool
		want   string
	}{
		{name: "initial White turn", want: "Turn: White"},
		{
			name: "Black turn", want: "Turn: Black",
			moves: []chess.Move{{From: chess.E2, To: chess.E4}},
		},
		{
			name: "Black in check", check: true, want: "Turn: Black (check)",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.D7, To: chess.D5},
				{From: chess.F1, To: chess.B5},
			},
		},
		{
			name: "White in check", check: true, want: "Turn: White (check)",
			moves: []chess.Move{
				{From: chess.D2, To: chess.D4},
				{From: chess.E7, To: chess.E5},
				{From: chess.G1, To: chess.F3},
				{From: chess.F8, To: chess.B4},
			},
		},
		{
			name: "check answered", want: "Turn: White",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.D7, To: chess.D5},
				{From: chess.F1, To: chess.B5},
				{From: chess.C7, To: chess.C6},
			},
		},
		{
			name: "Black wins", status: chess.StatusCheckmate, check: true, want: "Checkmate! Black wins.",
			moves: []chess.Move{
				{From: chess.F2, To: chess.F3},
				{From: chess.E7, To: chess.E5},
				{From: chess.G2, To: chess.G4},
				{From: chess.D8, To: chess.H4},
			},
		},
		{
			name: "White wins", status: chess.StatusCheckmate, check: true, want: "Checkmate! White wins.",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.F7, To: chess.F6},
				{From: chess.D2, To: chess.D4},
				{From: chess.G7, To: chess.G5},
				{From: chess.D1, To: chess.H5},
			},
		},
		{
			name: "stalemate", status: chess.StatusStalemate, want: "Stalemate! Draw.",
			moves: []chess.Move{
				{From: chess.E2, To: chess.E3},
				{From: chess.A7, To: chess.A5},
				{From: chess.D1, To: chess.H5},
				{From: chess.A8, To: chess.A6},
				{From: chess.H5, To: chess.A5},
				{From: chess.H7, To: chess.H5},
				{From: chess.A5, To: chess.C7},
				{From: chess.A6, To: chess.H6},
				{From: chess.H2, To: chess.H4},
				{From: chess.F7, To: chess.F6},
				{From: chess.C7, To: chess.D7},
				{From: chess.E8, To: chess.F7},
				{From: chess.D7, To: chess.B7},
				{From: chess.D8, To: chess.D3},
				{From: chess.B7, To: chess.B8},
				{From: chess.D3, To: chess.H7},
				{From: chess.B8, To: chess.C8},
				{From: chess.F7, To: chess.G6},
				{From: chess.C8, To: chess.E6},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			m := tui.NewModel(g)
			for _, move := range tc.moves {
				m = moveCursorTo(t, m, move.From)
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				m = moveCursorTo(t, m, move.To)
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				snapshot, err := g.Snapshot()
				if err != nil {
					t.Fatal(err)
				}
				if snapshot.LastMove == nil || *snapshot.LastMove != move {
					t.Fatalf("Expected setup move %v to be played, got %+v", move, snapshot)
				}
			}
			before, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			if before.Status != tc.status || before.Check != tc.check {
				t.Fatalf("Expected status/check %v/%v, got %v/%v", tc.status, tc.check, before.Status, before.Check)
			}
			length := g.HistoryLen()
			sections := strings.Split(ansi.Strip(m.View().Content), "\n\n")
			if len(sections) < 3 {
				t.Fatalf("Expected game status below board, got %q", m.View().Content)
			}
			if got := sections[2]; got != tc.want {
				t.Errorf("Expected status text %q, got %q", tc.want, got)
			}
			assertGameUnchanged(t, g, before, length)
		})
	}
}
