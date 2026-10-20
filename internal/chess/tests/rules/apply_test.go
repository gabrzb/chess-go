package rules_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestApplyingGeneratedMovesMatchesValidatedPath(t *testing.T) {
	for _, fen := range []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 20",
		"4k3/P7/8/8/8/8/8/4K3 w - - 0 20",
	} {
		position := perftPositionFromFEN(t, fen)
		for square := chess.A1; square <= chess.H8; square++ {
			piece, _ := position.Board.PieceAt(square)
			if piece.Type == chess.NoPiece || piece.Color != position.SideToMove {
				continue
			}
			moves, err := rules.LegalMovesFrom(position, square)
			if err != nil {
				t.Fatal(err)
			}
			for _, move := range moves {
				want, err := rules.PositionAfterMove(position, move)
				if err != nil {
					t.Fatal(err)
				}
				got, err := rules.PositionAfterLegalMove(position, move)
				if err != nil || got != want {
					t.Fatalf("%s, %v: position differs: %v", fen, move, err)
				}
			}
		}
	}
}
