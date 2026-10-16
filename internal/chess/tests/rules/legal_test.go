package rules_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestKingSafety(t *testing.T) {
	p := perftPositionFromFEN(t, "k3r3/8/8/8/8/8/4R3/4K3 w - - 0 1")
	if _, err := rules.PositionAfterMove(p, chess.Move{From: chess.E2, To: chess.D2}); !errors.Is(err, chess.ErrIllegalMove) {
		t.Error("Pinned rook exposed its king")
	}
	after(t, p, chess.Move{From: chess.E2, To: chess.E8})
	mustSetPiece(t, &p.Board, chess.E2, chess.Piece{})
	moves, err := rules.LegalMovesFrom(p, chess.E1)
	if err != nil || !slices.Contains(moves, chess.Move{From: chess.E1, To: chess.D1}) || slices.Contains(moves, chess.Move{From: chess.E1, To: chess.E2}) {
		t.Fatalf("Incorrect king escapes: %v, %v", moves, err)
	}
	p = perftPositionFromFEN(t, "4k3/4Q3/8/8/8/8/8/4K3 w - - 0 1")
	if _, err := rules.PositionAfterMove(p, chess.Move{From: chess.E7, To: chess.E8}); !errors.Is(err, chess.ErrIllegalMove) {
		t.Error("A king was captured")
	}
}

func TestPositionStatus(t *testing.T) {
	for _, tc := range []struct {
		fen   string
		want  chess.Status
		check bool
	}{
		{"7k/6Q1/5K2/8/8/8/8/8 b - - 0 1", chess.StatusCheckmate, true},
		{"7k/5K2/6Q1/8/8/8/8/8 b - - 0 1", chess.StatusStalemate, false},
		{"k3r3/8/8/8/8/8/8/4K3 w - - 0 1", chess.StatusOngoing, true},
	} {
		p := perftPositionFromFEN(t, tc.fen)
		status, err := rules.EvaluateStatus(p)
		check, checkErr := rules.IsInCheck(p.Board, p.SideToMove)
		if err != nil || checkErr != nil || status != tc.want || check != tc.check {
			t.Errorf("Incorrect status for %s: %v/%v, %v/%v", tc.fen, status, check, err, checkErr)
		}
	}
	if _, err := rules.EvaluateStatus(chess.Position{}); !errors.Is(err, chess.ErrInvalidColor) {
		t.Error("Expected invalid-color error")
	}
	if _, err := rules.EvaluateStatus(chess.Position{SideToMove: chess.White}); !errors.Is(err, chess.ErrKingNotFound) {
		t.Error("Expected missing-king error")
	}
}
