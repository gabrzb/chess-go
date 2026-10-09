package attacks_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

func TestSlidingAttacksSquare(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	got := attacks.SlidingAttacksSquare(
		board, chess.Square(27),
		chess.Square(31),
		chess.Direction{ColStep: 1, RowStep: 0},
	)

	if !got {
		t.Errorf("Expected an attack from 27 to 31, but got false")
	}
}

func TestSlidingAttacksSquareBlocked(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}
	blocker := chess.Piece{Type: chess.Pawn, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(29), blocker); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	got := attacks.SlidingAttacksSquare(
		board, chess.Square(27),
		chess.Square(31),
		chess.Direction{ColStep: 1, RowStep: 0},
	)

	if got {
		t.Errorf("Expected the piece at 29 to block the attack from 27 to 31, but got true")
	}
}

func TestSlidingAttacksSquareOwnPieceAtTarget(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}
	blocker := chess.Piece{Type: chess.Pawn, Color: chess.White}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(chess.Square(31), blocker); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	got := attacks.SlidingAttacksSquare(
		board, chess.Square(27),
		chess.Square(31),
		chess.Direction{ColStep: 1, RowStep: 0},
	)

	if !got {
		t.Errorf("Expected an attack from 27 to 31 even though there's a piece of the same color at 31, but got false")
	}
}

func TestRookAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   chess.Square
		want bool
	}{
		{"right", chess.Square(31), true},
		{"left", chess.Square(24), true},
		{"up", chess.Square(59), true},
		{"down", chess.Square(3), true},
		{"diagonal", chess.Square(36), false},
		{"knight move", chess.Square(37), false},
		{"same square", chess.Square(27), false},
		{"invalid destination", chess.NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			rook := chess.Piece{Type: chess.Rook, Color: chess.White}

			if err := board.SetPiece(chess.Square(27), rook); err != nil {
				t.Fatalf("Failed to set piece: %v", err)
			}

			got := attacks.RookAttacksSquare(board, chess.Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestBishopAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   chess.Square
		want bool
	}{
		{"up right", chess.Square(63), true},
		{"up left", chess.Square(48), true},
		{"down right", chess.Square(6), true},
		{"down left", chess.Square(0), true},
		{"horizontal", chess.Square(31), false},
		{"vertical", chess.Square(59), false},
		{"same square", chess.Square(27), false},
		{"invalid destination", chess.NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			bishop := chess.Piece{Type: chess.Bishop, Color: chess.White}

			if err := board.SetPiece(chess.Square(27), bishop); err != nil {
				t.Fatalf("Failed to set piece: %v", err)
			}

			got := attacks.BishopAttacksSquare(board, chess.Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestQueenAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   chess.Square
		want bool
	}{
		{"right", chess.Square(31), true},
		{"left", chess.Square(24), true},
		{"up", chess.Square(59), true},
		{"down", chess.Square(3), true},
		{"up right", chess.Square(63), true},
		{"up left", chess.Square(48), true},
		{"down right", chess.Square(6), true},
		{"down left", chess.Square(0), true},
		{"knight move", chess.Square(37), false},
		{"same square", chess.Square(27), false},
		{"invalid destination", chess.NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			queen := chess.Piece{Type: chess.Queen, Color: chess.White}

			if err := board.SetPiece(chess.Square(27), queen); err != nil {
				t.Fatalf("Failed to set piece: %v", err)
			}

			got := attacks.QueenAttacksSquare(board, chess.Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
