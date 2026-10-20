package rules_test

import (
	"testing"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/rules"
)

func TestApplyingGeneratedMovesMatchesValidatedPath(t *testing.T) {
	for _, fen := range []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 20",
		"4k3/P7/8/8/8/8/8/4K3 w - - 0 20",
	} {
		position := perftPositionFromFEN(t, fen)
		before := position
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
				if position != before {
					t.Fatalf("%s, %v: move mutated the parent position", fen, move)
				}
				if err != nil || got != want {
					t.Fatalf("%s, %v: position differs: %v", fen, move, err)
				}
			}
		}
	}
}

func TestSearchBranchesHaveIndependentBoards(t *testing.T) {
	parent := chess.NewInitialPosition()
	before := parent
	left, err := rules.PositionAfterMove(parent, chess.Move{From: chess.E2, To: chess.E4})
	if err != nil {
		t.Fatal(err)
	}
	right, err := rules.PositionAfterMove(parent, chess.Move{From: chess.D2, To: chess.D4})
	if err != nil {
		t.Fatal(err)
	}
	rightBefore := right
	if _, err := left.Board.RemovePiece(chess.E1); err != nil {
		t.Fatal(err)
	}
	left.CastlingRights = 0
	left.EnPassant = chess.NoSquare
	left.HalfmoveClock++
	left.FullmoveNumber++
	if parent != before || right != rightBefore {
		t.Fatal("modifying one search branch changed the parent or sibling")
	}
	if left == right {
		t.Fatal("different branches produced the same position")
	}
}
