package rules_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func perft(t *testing.T, position chess.Position, depth int) uint64 {
	t.Helper()
	if depth < 0 {
		t.Fatalf("Invalid perft depth: %d", depth)
	}
	if depth == 0 {
		return 1
	}

	var total uint64
	for from := chess.A1; from <= chess.H8; from++ {
		piece, _ := position.Board.PieceAt(from)
		if piece.Type == chess.NoPiece || piece.Color != position.SideToMove {
			continue
		}

		moves, err := rules.LegalMovesFrom(position, from)
		if err != nil {
			t.Fatalf("Failed to generate moves from %v at depth %d: %v", from, depth, err)
		}

		if depth == 1 {
			total += uint64(len(moves))
			continue
		}

		for _, move := range moves {
			next, err := rules.PositionAfterMove(position, move)
			if err != nil {
				t.Fatalf("Failed to apply %v at depth %d: %v", move, depth, err)
			}

			total += perft(t, next, depth-1)
		}
	}
	return total
}

func assertPerft(t *testing.T, position chess.Position, depth int, want uint64) {
	t.Helper()

	before := position

	if got := perft(t, position, depth); got != want {
		t.Errorf("Perft depth %d: expected %d, got %d", depth, want, got)
	}

	if position != before {
		t.Error("Expected perft to preserve original position")
	}
}

func TestPerftInitialPosition(t *testing.T) {
	assertPerft(t, chess.NewInitialPosition(), 4, 197281)
}
