package ai_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
)

func setPiece(t *testing.T, board *chess.Board, square chess.Square, kind chess.PieceType, color chess.Color) {
	t.Helper()
	if err := board.SetPiece(square, chess.Piece{Type: kind, Color: color}); err != nil {
		t.Fatal(err)
	}
}

func TestEvaluationBalanceAndMaterial(t *testing.T) {
	initial := chess.NewInitialPosition()
	if got := ai.Evaluate(initial.Board); got != 0 {
		t.Errorf("Expected balanced initial position, got %d", got)
	}
	for _, color := range []chess.Color{chess.White, chess.Black} {
		t.Run(map[chess.Color]string{chess.White: "White", chess.Black: "Black"}[color], func(t *testing.T) {
			var pawn, queen chess.Board
			setPiece(t, &pawn, chess.D4, chess.Pawn, color)
			setPiece(t, &queen, chess.D4, chess.Queen, color)
			sign := 1
			if color == chess.Black {
				sign = -1
			}
			if ai.Evaluate(pawn)*sign <= 0 || ai.Evaluate(queen)*sign <= ai.Evaluate(pawn)*sign {
				t.Error("Expected material to favor its owner and a queen to outweigh a pawn")
			}
		})
	}
}

func TestEvaluationPositionAndColorSymmetry(t *testing.T) {
	for _, tc := range []struct {
		name         string
		kind         chess.PieceType
		weak, strong chess.Square
	}{
		{"knight centralization", chess.Knight, chess.A1, chess.D4},
		{"bishop centralization", chess.Bishop, chess.A1, chess.D4},
		{"pawn advance", chess.Pawn, chess.D2, chess.D5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var weak, strong, mirrored chess.Board
			setPiece(t, &weak, tc.weak, tc.kind, chess.White)
			setPiece(t, &strong, tc.strong, tc.kind, chess.White)
			mirror, ok := chess.A8.Offset(tc.strong.Column(), -tc.strong.Row())
			if !ok {
				t.Fatal("Invalid mirrored test square")
			}
			setPiece(t, &mirrored, mirror, tc.kind, chess.Black)
			before := strong
			if ai.Evaluate(strong) <= ai.Evaluate(weak) {
				t.Error("Expected improved piece placement to increase its score")
			}
			if ai.Evaluate(mirrored) != -ai.Evaluate(strong) {
				t.Error("Expected mirrored colors to receive opposite scores")
			}
			if strong != before {
				t.Error("Expected evaluation to preserve the board")
			}
		})
	}
}
