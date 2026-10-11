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

func TestJumpAttacks(t *testing.T) {
	pieces := []struct {
		name         string
		attack       func(chess.Square, chess.Square) bool
		hits, misses []chess.Square
	}{
		{"knight", attacks.KnightAttacksSquare, []chess.Square{37, 21, 33, 17, 44, 42, 12, 10}, []chess.Square{36, 27, chess.NoSquare}},
		{"king", attacks.KingAttacksSquare, []chess.Square{28, 26, 35, 19, 36, 34, 20, 18}, []chess.Square{29, 43, 27, chess.NoSquare}},
	}
	for _, piece := range pieces {
		t.Run(piece.name, func(t *testing.T) {
			for _, group := range []struct {
				name    string
				squares []chess.Square
				want    bool
			}{
				{"attacked", piece.hits, true}, {"not attacked", piece.misses, false},
			} {
				t.Run(group.name, func(t *testing.T) {
					for _, to := range group.squares {
						t.Run(to.String(), func(t *testing.T) {
							if got := piece.attack(27, to); got != group.want {
								t.Errorf("Expected attack to %v = %v, got %v", to, group.want, got)
							}
						})
					}
				})
			}
		})
	}
}
