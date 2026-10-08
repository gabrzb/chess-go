package chess

import "testing"

func TestSlidingAttacksSquare(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	got := slidingAttacksSquare(
		board, Square(27),
		Square(31),
		direction{colStep: 1, rowStep: 0},
	)

	if !got {
		t.Errorf("Expected an attack from 27 to 31, but got false")
	}
}

func TestSlidingAttacksSquareBlocked(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}
	blocker := Piece{Type: Pawn, Color: Black}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(29), blocker); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	got := slidingAttacksSquare(
		board, Square(27),
		Square(31),
		direction{colStep: 1, rowStep: 0},
	)

	if got {
		t.Errorf("Expected the piece at 29 to block the attack from 27 to 31, but got true")
	}
}

func TestSlidingAttacksSquareOwnPieceAtTarget(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}
	blocker := Piece{Type: Pawn, Color: White}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(31), blocker); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	got := slidingAttacksSquare(
		board, Square(27),
		Square(31),
		direction{colStep: 1, rowStep: 0},
	)

	if !got {
		t.Errorf("Expected an attack from 27 to 31 even though there's a piece of the same color at 31, but got false")
	}
}
