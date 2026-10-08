package chess

import "testing"

func TestKingAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   Square
		want bool
	}{
		{"e4", Square(28), true},
		{"c4", Square(26), true},
		{"d5", Square(35), true},
		{"d3", Square(19), true},
		{"e5", Square(36), true},
		{"c5", Square(34), true},
		{"e3", Square(20), true},
		{"c3", Square(18), true},
		{"two columns away", Square(29), false},
		{"two rows away", Square(43), false},
		{"same square", Square(27), false},
		{"invalid destination", NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := kingAttacksSquare(Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
