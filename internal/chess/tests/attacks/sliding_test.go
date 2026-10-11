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
		{"blocked path", []placedPiece{{29, chess.Piece{Type: chess.Pawn, Color: chess.Black}}}, false},
		{"own piece at target", []placedPiece{{31, chess.Piece{Type: chess.Pawn, Color: chess.White}}}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pieces := append([]placedPiece{{27, rook}}, tc.extra...)
			board := boardWithPieces(t, pieces...)
			got := attacks.SlidingAttacksSquare(board, 27, 31, chess.Direction{ColStep: 1})
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
			{"right", 31, true}, {"left", 24, true}, {"up", 59, true}, {"down", 3, true},
			{"diagonal", 36, false}, {"knight move", 37, false}, {"same square", 27, false}, {"invalid destination", chess.NoSquare, false},
		}},
		{"bishop", chess.Bishop, attacks.BishopAttacksSquare, []attackCase{
			{"up right", 63, true}, {"up left", 48, true}, {"down right", 6, true}, {"down left", 0, true},
			{"horizontal", 31, false}, {"vertical", 59, false}, {"same square", 27, false}, {"invalid destination", chess.NoSquare, false},
		}},
		{"queen", chess.Queen, attacks.QueenAttacksSquare, []attackCase{
			{"right", 31, true}, {"left", 24, true}, {"up", 59, true}, {"down", 3, true},
			{"up right", 63, true}, {"up left", 48, true}, {"down right", 6, true}, {"down left", 0, true},
			{"knight move", 37, false}, {"same square", 27, false}, {"invalid destination", chess.NoSquare, false},
		}},
	}
	for _, piece := range pieces {
		t.Run(piece.name, func(t *testing.T) {
			board := boardWithPieces(t, placedPiece{27, chess.Piece{Type: piece.pieceType, Color: chess.White}})
			for _, tc := range piece.cases {
				t.Run(tc.name, func(t *testing.T) {
					if got := piece.attack(board, 27, tc.to); got != tc.want {
						t.Errorf("Expected attack to %v = %v, got %v", tc.to, tc.want, got)
					}
				})
			}
		})
	}
}
