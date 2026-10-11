package attacks_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

func TestIsSquareAttacked(t *testing.T) {
	rook := placedPiece{27, chess.Piece{Type: chess.Rook, Color: chess.Black}}
	cases := []struct {
		name   string
		pieces []placedPiece
		to     chess.Square
		by     chess.Color
		want   bool
	}{
		{"empty board", nil, 27, chess.Black, false},
		{"attacking color", []placedPiece{rook}, 31, chess.Black, true},
		{"other color", []placedPiece{rook}, 31, chess.White, false},
		{"pawn", []placedPiece{{34, chess.Piece{Type: chess.Pawn, Color: chess.Black}}}, 27, chess.Black, true},
		{"knight", []placedPiece{{10, chess.Piece{Type: chess.Knight, Color: chess.Black}}}, 27, chess.Black, true},
		{"bishop", []placedPiece{{0, chess.Piece{Type: chess.Bishop, Color: chess.Black}}}, 27, chess.Black, true},
		{"rook", []placedPiece{{3, chess.Piece{Type: chess.Rook, Color: chess.Black}}}, 27, chess.Black, true},
		{"queen", []placedPiece{{63, chess.Piece{Type: chess.Queen, Color: chess.Black}}}, 27, chess.Black, true},
		{"king", []placedPiece{{26, chess.Piece{Type: chess.King, Color: chess.Black}}}, 27, chess.Black, true},
		{"continue searching", []placedPiece{{0, chess.Piece{Type: chess.Knight, Color: chess.Black}}, rook}, 31, chess.Black, true},
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
