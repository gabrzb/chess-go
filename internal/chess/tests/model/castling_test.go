package model_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestCastlingRightHas(t *testing.T) {
	whiteRights := chess.WhiteKingSide | chess.WhiteQueenSide
	blackRights := chess.BlackKingSide | chess.BlackQueenSide

	cases := []struct {
		name   string
		rights chess.CastlingRights
		right  chess.CastlingRights
		want   bool
	}{
		{"white king side", whiteRights, chess.WhiteKingSide, true},
		{"white queen side", whiteRights, chess.WhiteQueenSide, true},
		{"black king side absent", whiteRights, chess.BlackKingSide, false},
		{"both white rights", whiteRights, whiteRights, true},
		{"partially present", whiteRights, chess.WhiteKingSide | chess.BlackKingSide, false},
		{"black king side", blackRights, chess.BlackKingSide, true},
		{"black queen side", blackRights, chess.BlackQueenSide, true},
		{"empty rights", 0, chess.WhiteKingSide, false},
		{"empty query", whiteRights, 0, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.rights.Has(tc.right)

			if got != tc.want {
				t.Errorf("Expected %v, got %v", tc.want, got)
			}
		})
	}
}
