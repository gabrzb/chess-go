package chess

import "testing"

func TestQueenAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   Square
		want bool
	}{
		{"right", Square(31), true},
		{"left", Square(24), true},
		{"up", Square(59), true},
		{"down", Square(3), true},
		{"up right", Square(63), true},
		{"up left", Square(48), true},
		{"down right", Square(6), true},
		{"down left", Square(0), true},
		{"knight move", Square(37), false},
		{"same square", Square(27), false},
		{"invalid destination", NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board Board
			queen := Piece{Type: Queen, Color: White}

			if err := board.SetPiece(Square(27), queen); err != nil {
				t.Fatalf("Failed to set piece: %v", err)
			}

			got := queenAttacksSquare(board, Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
