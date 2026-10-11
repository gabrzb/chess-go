package attacks_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

func TestSlidingAttacksSquare(t *testing.T) {
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}
	cases := []struct {
		name  string
		extra []placedPiece
		want  bool
	}{
		{"clear path", nil, true},
		{"blocked path", []placedPiece{{chess.F4, chess.Piece{Type: chess.Pawn, Color: chess.Black}}}, false},
		{"own piece at target", []placedPiece{{chess.H4, chess.Piece{Type: chess.Pawn, Color: chess.White}}}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pieces := append([]placedPiece{{chess.D4, rook}}, tc.extra...)
			board := boardWithPieces(t, pieces...)
			got := attacks.SlidingAttacksSquare(board, chess.D4, chess.H4, chess.Direction{ColStep: 1})
			if got != tc.want {
				t.Errorf("Expected attack=%v, got %v", tc.want, got)
			}
		})
	}
}

func TestSlidingPieceAttacks(t *testing.T) {
	pieces := []struct {
		name      string
		pieceType chess.PieceType
		attack    func(chess.Board, chess.Square, chess.Square) bool
		cases     []attackCase
	}{
		{"rook", chess.Rook, attacks.RookAttacksSquare, []attackCase{
			{"right", chess.H4, true}, {"left", chess.A4, true}, {"up", chess.D8, true}, {"down", chess.D1, true},
			{"diagonal", chess.E5, false}, {"knight move", chess.F5, false}, {"same square", chess.D4, false}, {"invalid destination", chess.NoSquare, false},
		}},
		{"bishop", chess.Bishop, attacks.BishopAttacksSquare, []attackCase{
			{"up right", chess.H8, true}, {"up left", chess.A7, true}, {"down right", chess.G1, true}, {"down left", chess.A1, true},
			{"horizontal", chess.H4, false}, {"vertical", chess.D8, false}, {"same square", chess.D4, false}, {"invalid destination", chess.NoSquare, false},
		}},
		{"queen", chess.Queen, attacks.QueenAttacksSquare, []attackCase{
			{"right", chess.H4, true}, {"left", chess.A4, true}, {"up", chess.D8, true}, {"down", chess.D1, true},
			{"up right", chess.H8, true}, {"up left", chess.A7, true}, {"down right", chess.G1, true}, {"down left", chess.A1, true},
			{"knight move", chess.F5, false}, {"same square", chess.D4, false}, {"invalid destination", chess.NoSquare, false},
		}},
	}
	for _, piece := range pieces {
		t.Run(piece.name, func(t *testing.T) {
			board := boardWithPieces(t, placedPiece{chess.D4, chess.Piece{Type: piece.pieceType, Color: chess.White}})
			for _, tc := range piece.cases {
				t.Run(tc.name, func(t *testing.T) {
					if got := piece.attack(board, chess.D4, tc.to); got != tc.want {
						t.Errorf("Expected attack to %v = %v, got %v", tc.to, tc.want, got)
					}
				})
			}
		})
	}
}
