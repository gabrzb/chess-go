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
		{"white left diagonal", chess.White, chess.C5, true},
		{"white right diagonal", chess.White, chess.E5, true},
		{"white forward", chess.White, chess.D5, false},
		{"black left diagonal", chess.Black, chess.C3, true},
		{"black right diagonal", chess.Black, chess.E3, true},
		{"black forward", chess.Black, chess.D3, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := attacks.PawnAttacksSquare(chess.D4, tc.to, tc.color)

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
		{"knight", attacks.KnightAttacksSquare, []chess.Square{chess.F5, chess.F3, chess.B5, chess.B3, chess.E6, chess.C6, chess.E2, chess.C2}, []chess.Square{chess.E5, chess.D4, chess.NoSquare}},
		{"king", attacks.KingAttacksSquare, []chess.Square{chess.E4, chess.C4, chess.D5, chess.D3, chess.E5, chess.C5, chess.E3, chess.C3}, []chess.Square{chess.F4, chess.D6, chess.D4, chess.NoSquare}},
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
							if got := piece.attack(chess.D4, to); got != group.want {
								t.Errorf("Expected attack to %v = %v, got %v", to, group.want, got)
							}
						})
					}
				})
			}
		})
	}
}
