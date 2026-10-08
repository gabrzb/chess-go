package chess

import "testing"

func TestBishopAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   Square
		want bool
	}{
		{"up right", Square(63), true},
		{"up left", Square(48), true},
		{"down right", Square(6), true},
		{"down left", Square(0), true},
		{"horizontal", Square(31), false},
		{"vertical", Square(59), false},
		{"same square", Square(27), false},
		{"invalid destination", NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board Board
			bishop := Piece{Type: Bishop, Color: White}

			if err := board.SetPiece(Square(27), bishop); err != nil {
				t.Fatalf("Failed to set piece: %v", err)
			}

			got := bishopAttacksSquare(board, Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
