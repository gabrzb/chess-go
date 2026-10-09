package attacks_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

func TestPawnAttackSquare(t *testing.T) {
	cases := []struct {
		name  string
		color chess.Color
		to    chess.Square
		want  bool
	}{
		{"white left diagonal", chess.White, chess.Square(34), true},
		{"white right diagonal", chess.White, chess.Square(36), true},
		{"white forward", chess.White, chess.Square(35), false},
		{"black left diagonal", chess.Black, chess.Square(18), true},
		{"black right diagonal", chess.Black, chess.Square(20), true},
		{"black forward", chess.Black, chess.Square(19), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := attacks.PawnAttacksSquare(chess.Square(27), tc.to, tc.color)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestKnightAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   chess.Square
		want bool
	}{
		{"f5", chess.Square(37), true},
		{"f3", chess.Square(21), true},
		{"b5", chess.Square(33), true},
		{"b3", chess.Square(17), true},
		{"e6", chess.Square(44), true},
		{"c6", chess.Square(42), true},
		{"e2", chess.Square(12), true},
		{"c2", chess.Square(10), true},
		{"adjacent diagonal", chess.Square(36), false},
		{"same square", chess.Square(27), false},
		{"invalid destination", chess.NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := attacks.KnightAttacksSquare(chess.Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestKingAttacks(t *testing.T) {
	cases := []struct {
		name string
		to   chess.Square
		want bool
	}{
		{"e4", chess.Square(28), true},
		{"c4", chess.Square(26), true},
		{"d5", chess.Square(35), true},
		{"d3", chess.Square(19), true},
		{"e5", chess.Square(36), true},
		{"c5", chess.Square(34), true},
		{"e3", chess.Square(20), true},
		{"c3", chess.Square(18), true},
		{"two columns away", chess.Square(29), false},
		{"two rows away", chess.Square(43), false},
		{"same square", chess.Square(27), false},
		{"invalid destination", chess.NoSquare, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := attacks.KingAttacksSquare(chess.Square(27), tc.to)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
