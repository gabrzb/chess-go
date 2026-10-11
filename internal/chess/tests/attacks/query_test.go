package attacks_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

func TestIsSquareAttacked(t *testing.T) {
	rook := placedPiece{chess.D4, chess.Piece{Type: chess.Rook, Color: chess.Black}}
	cases := []struct {
		name   string
		pieces []placedPiece
		to     chess.Square
		by     chess.Color
		want   bool
	}{
		{"empty board", nil, chess.D4, chess.Black, false},
		{"attacking color", []placedPiece{rook}, chess.H4, chess.Black, true},
		{"other color", []placedPiece{rook}, chess.H4, chess.White, false},
		{"pawn", []placedPiece{{chess.C5, chess.Piece{Type: chess.Pawn, Color: chess.Black}}}, chess.D4, chess.Black, true},
		{"knight", []placedPiece{{chess.C2, chess.Piece{Type: chess.Knight, Color: chess.Black}}}, chess.D4, chess.Black, true},
		{"bishop", []placedPiece{{chess.A1, chess.Piece{Type: chess.Bishop, Color: chess.Black}}}, chess.D4, chess.Black, true},
		{"rook", []placedPiece{{chess.D1, chess.Piece{Type: chess.Rook, Color: chess.Black}}}, chess.D4, chess.Black, true},
		{"queen", []placedPiece{{chess.H8, chess.Piece{Type: chess.Queen, Color: chess.Black}}}, chess.D4, chess.Black, true},
		{"king", []placedPiece{{chess.C4, chess.Piece{Type: chess.King, Color: chess.Black}}}, chess.D4, chess.Black, true},
		{"continue searching", []placedPiece{{chess.A1, chess.Piece{Type: chess.Knight, Color: chess.Black}}, rook}, chess.H4, chess.Black, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			board := boardWithPieces(t, tc.pieces...)
			if got := attacks.IsSquareAttacked(board, tc.to, tc.by); got != tc.want {
				t.Errorf("Expected attack on %v by %v = %v, got %v", tc.to, tc.by, tc.want, got)
			}
		})
	}
}
