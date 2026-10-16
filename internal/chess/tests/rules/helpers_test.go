package rules_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func mustSetPiece(t *testing.T, board *chess.Board, square chess.Square, piece chess.Piece) {
	t.Helper()
	if err := board.SetPiece(square, piece); err != nil {
		t.Fatal(err)
	}
}

func after(t *testing.T, position chess.Position, move chess.Move) chess.Position {
	t.Helper()
	next, err := rules.PositionAfterMove(position, move)
	if err != nil {
		t.Fatalf("Applying %v: %v", move, err)
	}
	return next
}

func assertPiece(t *testing.T, position chess.Position, square chess.Square, want chess.Piece) {
	t.Helper()
	if got, _ := position.Board.PieceAt(square); got != want {
		t.Errorf("At %v: got %v, want %v", square, got, want)
	}
}
