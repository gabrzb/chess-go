package chess

import (
	"errors"
	"testing"
)

func TestFindKing(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(4), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(60), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	cases := []struct {
		name  string
		color Color
		want  Square
	}{
		{"white", White, Square(4)},
		{"black", Black, Square(60)},
		{"invalid color", NoColor, NoSquare},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := findKing(board, tc.color)

			if got != tc.want {
				t.Errorf("FindKing(%v) = %v; want %v", tc.color, got, tc.want)
			}
		})
	}
}

func TestFindKingMissing(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}

	if err := board.SetPiece(Square(4), rook); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	got := findKing(board, White)

	if got != NoSquare {
		t.Errorf("FindKing(White) = %v; want NoSquare", got)
	}
}

func TestIsInCheckWhite(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(4), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(60), Piece{Type: Rook, Color: Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	got, err := IsInCheck(board, White)
	if err != nil {
		t.Fatalf("IsInCheck returned an error: %v", err)
	}

	if !got {
		t.Errorf("IsInCheck(White) = %v; want true", got)
	}
}

func TestIsInCheckBlack(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(4), Piece{Type: Rook, Color: White}); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	if err := board.SetPiece(Square(60), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	got, err := IsInCheck(board, Black)
	if err != nil {
		t.Fatalf("IsInCheck returned an error: %v", err)
	}

	if !got {
		t.Errorf("IsInCheck(White) = %v; want true", got)
	}
}

func TestIsInCheckBlocked(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(4), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(60), Piece{Type: Rook, Color: Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	if err := board.SetPiece(Square(28), Piece{Type: Pawn, Color: White}); err != nil {
		t.Fatalf("Failed to set white pawn: %v", err)
	}

	got, err := IsInCheck(board, White)
	if err != nil {
		t.Fatalf("IsInCheck returned an error: %v", err)
	}

	if got {
		t.Errorf("IsInCheck(White) = %v; want false", got)
	}
}

func TestIsInCheckInvalidColor(t *testing.T) {
	var board Board

	_, err := IsInCheck(board, NoColor)

	if !errors.Is(err, ErrInvalidColor) {
		t.Errorf("IsInCheck returned error %v; want ErrInvalidColor", err)
	}
}

func TestIsInCheckMissingKing(t *testing.T) {
	var board Board

	_, err := IsInCheck(board, White)

	if !errors.Is(err, ErrKingNotFound) {
		t.Errorf("IsInCheck returned error %v; want ErrKingNotFound", err)
	}
}
