package chess

import "testing"

func TestKnightAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   Square
		want bool
	}{
		{"f5", Square(37), true},
		{"f3", Square(21), true},
		{"b5", Square(33), true},
		{"b3", Square(17), true},
		{"e6", Square(44), true},
		{"c6", Square(42), true},
		{"e2", Square(12), true},
		{"c2", Square(10), true},
		{"adjacent diagonal", Square(36), false},
		{"same square", Square(27), false},
		{"invalid destination", NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := knightAttacksSquare(Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
