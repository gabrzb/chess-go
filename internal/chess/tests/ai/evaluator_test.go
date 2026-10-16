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

func TestEvaluation(t *testing.T) {
	if score := ai.Evaluate(chess.NewInitialPosition().Board); score != 0 {
		t.Fatal("Initial position is not balanced")
	}
	var white, black, edge chess.Board
	setPiece(t, &white, chess.D4, chess.Knight, chess.White)
	setPiece(t, &black, chess.D5, chess.Knight, chess.Black)
	setPiece(t, &edge, chess.A1, chess.Knight, chess.White)
	if ai.Evaluate(white) <= 0 || ai.Evaluate(black) != -ai.Evaluate(white) || ai.Evaluate(white) <= ai.Evaluate(edge) {
		t.Error("Evaluation must respect material, color symmetry and centralization")
	}
	setPiece(t, &white, chess.D4, chess.Queen, chess.White)
	if ai.Evaluate(white) <= -ai.Evaluate(black) {
		t.Error("Queen must outweigh a knight")
	}
}
