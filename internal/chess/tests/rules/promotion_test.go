package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestPromotionChoicesAndValidation(t *testing.T) {
	for _, tc := range []struct {
		fen           string
		move          chess.Move
		capturedRight chess.CastlingRights
	}{
		{"4k3/P7/8/8/8/8/8/4K3 w - - 7 20", chess.Move{From: chess.A7, To: chess.A8}, 0},
		{"r3k3/1P6/8/8/8/8/8/4K3 w q - 7 20", chess.Move{From: chess.B7, To: chess.A8}, chess.BlackQueenSide},
		{"4k3/8/8/8/8/8/p7/4K3 b - - 7 20", chess.Move{From: chess.A2, To: chess.A1}, 0},
		{"4k3/8/8/8/8/8/1p6/R3K3 b Q - 7 20", chess.Move{From: chess.B2, To: chess.A1}, chess.WhiteQueenSide},
	} {
		p := perftPositionFromFEN(t, tc.fen)
		for _, kind := range []chess.PieceType{chess.Queen, chess.Rook, chess.Bishop, chess.Knight} {
			move := tc.move
			move.Promotion = kind
			next := after(t, p, move)
			assertPiece(t, next, move.From, chess.Piece{})
			assertPiece(t, next, move.To, chess.Piece{Type: kind, Color: p.SideToMove})
			if next.SideToMove == p.SideToMove || next.HalfmoveClock != 0 || next.CastlingRights != p.CastlingRights.Without(tc.capturedRight) {
				t.Errorf("Incorrect promotion metadata for %v", move)
			}
		}
		for _, kind := range []chess.PieceType{chess.NoPiece, chess.Pawn, chess.King} {
			move := tc.move
			move.Promotion = kind
			if _, err := rules.PositionAfterMove(p, move); !errors.Is(err, chess.ErrIllegalMove) {
				t.Errorf("Accepted invalid promotion %v", move)
			}
		}
	}
}
