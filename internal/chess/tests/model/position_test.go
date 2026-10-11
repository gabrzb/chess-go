package model_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestNewPosition(t *testing.T) {
	var whiteBoard chess.Board
	if err := whiteBoard.SetPiece(4, chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set king: %v", err)
	}
	cases := []struct {
		name  string
		color chess.Color
		board chess.Board
	}{
		{"white", chess.White, whiteBoard},
		{"black", chess.Black, chess.Board{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := chess.NewPosition(tc.board, tc.color)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			want := chess.Position{Board: tc.board, SideToMove: tc.color, EnPassant: chess.NoSquare, FullmoveNumber: 1}
			if got != want {
				t.Error("Expected supplied board and turn, no special rights, halfmove=0 and fullmove=1")
			}
		})
	}
}

func TestNewPositionInvalidColor(t *testing.T) {
	var board chess.Board

	_, err := chess.NewPosition(board, chess.NoColor)

	if !errors.Is(err, chess.ErrInvalidColor) {
		t.Errorf("Expected ErrInvalidColor, got %v", err)
	}
}
