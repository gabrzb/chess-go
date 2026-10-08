package chess

import "testing"

func TestPawnAttackSquare(t *testing.T) {
	cases := []struct {
		name  string
		color Color
		to    Square
		want  bool
	}{
		{"white left diagonal", White, Square(34), true},
		{"white right diagonal", White, Square(36), true},
		{"white forward", White, Square(35), false},
		{"black left diagonal", Black, Square(18), true},
		{"black right diagonal", Black, Square(20), true},
		{"black forward", Black, Square(19), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := pawnAttacksSquare(Square(27), tc.to, tc.color)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
