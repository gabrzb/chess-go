package attacks_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

type placedPiece struct {
	square chess.Square
	piece  chess.Piece
}
type attackCase struct {
	name string
	to   chess.Square
	want bool
}

func boardWithPieces(t *testing.T, pieces ...placedPiece) chess.Board {
	t.Helper()
	var board chess.Board
	for _, item := range pieces {
		if err := board.SetPiece(item.square, item.piece); err != nil {
			t.Fatalf("Failed to set piece at %v: %v", item.square, err)
		}
	}
	return board
}
