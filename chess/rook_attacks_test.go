package chess

import "testing"

func TestRookAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   Square
		want bool
	}{
		{"right", Square(31), true},
		{"left", Square(24), true},
		{"up", Square(59), true},
		{"down", Square(3), true},
		{"diagonal", Square(36), false},
		{"knight move", Square(37), false},
		{"same square", Square(27), false},
		{"invalid destination", NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board Board
			rook := Piece{Type: Rook, Color: White}

			if err := board.SetPiece(Square(27), rook); err != nil {
				t.Fatalf("Failed to set piece: %v", err)
			}

			got := rookAttacksSquare(board, Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
