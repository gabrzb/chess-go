package chess

import "testing"

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
