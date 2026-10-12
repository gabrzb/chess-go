package tui_test

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui"
)

func TestPlayMovesThroughInterface(t *testing.T) {
	cases := []struct {
		name   string
		moves  []chess.Move
		status chess.Status
	}{
		{
			name: "both colors and capture", status: chess.StatusOngoing,
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.D7, To: chess.D5},
				{From: chess.E4, To: chess.D5},
			},
		},
		{
			name: "black delivers checkmate", status: chess.StatusCheckmate,
			moves: []chess.Move{
				{From: chess.F2, To: chess.F3},
				{From: chess.E7, To: chess.E5},
				{From: chess.G2, To: chess.G4},
				{From: chess.D8, To: chess.H4},
			},
		},
		{
			name: "white delivers checkmate", status: chess.StatusCheckmate,
			moves: []chess.Move{
				{From: chess.E2, To: chess.E4},
				{From: chess.F7, To: chess.F6},
				{From: chess.D2, To: chess.D4},
				{From: chess.G7, To: chess.G5},
				{From: chess.D1, To: chess.H5},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			m := tui.NewModel(g)
			for index, move := range tc.moves {
				before, err := g.Snapshot()
				if err != nil {
					t.Fatal(err)
				}
				piece, _ := before.Board.PieceAt(move.From)
				m = moveCursorTo(t, m, move.From)
				symbol := ansi.Strip(boardCells(t, m)[move.From])
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				assertGameUnchanged(t, g, before, index+1)
				m = moveCursorTo(t, m, move.To)
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				after, err := g.Snapshot()
				if err != nil {
					t.Fatal(err)
				}
				if after.LastMove == nil || *after.LastMove != move {
					t.Fatalf("Expected interface to play %v, got snapshot %+v and view %q", move, after, m.View().Content)
				}
				wantTurn := chess.White
				if before.SideToMove == chess.White {
					wantTurn = chess.Black
				}
				if after.SideToMove != wantTurn || g.HistoryLen() != index+2 {
					t.Errorf("Expected turn %v and history length %d, got %v and %d", wantTurn, index+2, after.SideToMove, g.HistoryLen())
				}
				if source, _ := after.Board.PieceAt(move.From); source.Type != chess.NoPiece {
					t.Errorf("Expected empty source %v, got %v", move.From, source)
				}
				if target, _ := after.Board.PieceAt(move.To); target != piece {
					t.Errorf("Expected %v at %v, got %v", piece, move.To, target)
				}
				cells := boardCells(t, m)
				if ansi.Strip(cells[move.From]) != "." || ansi.Strip(cells[move.To]) != symbol {
					t.Errorf("Expected updated board in view after %v", move)
				}
				assertSelectionHighlights(t, m, chess.NoSquare, nil, move.To)
			}
			final, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			if final.Status != tc.status || final.Check != (tc.status == chess.StatusCheckmate) {
				t.Errorf("Expected status %v with matching check state, got %+v", tc.status, final)
			}
			if tc.status == chess.StatusCheckmate {
				move := chess.Move{From: chess.A2, To: chess.A3}
				if final.SideToMove == chess.Black {
					move = chess.Move{From: chess.A7, To: chess.A6}
				}
				length := g.HistoryLen()
				m = moveCursorTo(t, m, move.From)
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				m = moveCursorTo(t, m, move.To)
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				if !strings.Contains(m.View().Content, game.ErrGameOver.Error()) {
					t.Errorf("Expected game over message in view, got %q", m.View().Content)
				}
				assertGameUnchanged(t, g, final, length)
			}
		})
	}
}

func TestRejectedMovePreservesSelection(t *testing.T) {
	for _, side := range []struct {
		name        string
		setup       []chess.Move
		origin, own chess.Square
		legal       []chess.Square
	}{
		{name: "white", origin: chess.A2, own: chess.B2, legal: []chess.Square{chess.A3, chess.A4}},
		{name: "black", setup: []chess.Move{{From: chess.E2, To: chess.E4}}, origin: chess.A7, own: chess.B7, legal: []chess.Square{chess.A6, chess.A5}},
	} {
		for _, destination := range []chess.Square{chess.D4, side.own, side.origin} {
			t.Run(side.name+"/"+destination.String(), func(t *testing.T) {
				g := game.NewGame()
				for _, move := range side.setup {
					if _, err := g.Play(move); err != nil {
						t.Fatal(err)
					}
				}
				m := moveCursorTo(t, tui.NewModel(g), side.origin)
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				m = moveCursorTo(t, m, destination)
				beforeBoard := boardText(t, m)
				before, err := g.Snapshot()
				if err != nil {
					t.Fatal(err)
				}
				length := g.HistoryLen()
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				if !strings.Contains(m.View().Content, chess.ErrIllegalMove.Error()) {
					t.Errorf("Expected illegal move message in view, got %q", m.View().Content)
				}
				if boardText(t, m) != beforeBoard {
					t.Error("Expected rejected move to preserve visible board, selection and cursor")
				}
				assertSelectionHighlights(t, m, side.origin, side.legal, destination)
				assertGameUnchanged(t, g, before, length)
			})
		}
	}
}

func TestSuccessfulMoveClearsError(t *testing.T) {
	m, g := modelAt(t, chess.A2)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	m = moveCursorTo(t, m, chess.D4)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if !strings.Contains(m.View().Content, chess.ErrIllegalMove.Error()) {
		t.Fatal("Expected setup error message")
	}
	m = moveCursorTo(t, m, chess.A4)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	after, err := g.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	wantMove := chess.Move{From: chess.A2, To: chess.A4}
	if after.LastMove == nil || *after.LastMove != wantMove || g.HistoryLen() != 2 {
		t.Fatalf("Expected recovery move %v to be played once, got %+v", wantMove, after)
	}
	expected := moveCursorTo(t, tui.NewModel(g), chess.A4)
	if !reflect.DeepEqual(m.View(), expected.View()) {
		t.Error("Expected successful move to clear message, selection and destination highlights")
	}
}
