package engine_test

import (
	"testing"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/engine"
)

func setPiece(t *testing.T, board *chess.Board, square chess.Square, kind chess.PieceType, color chess.Color) {
	t.Helper()
	if err := board.SetPiece(square, chess.Piece{Type: kind, Color: color}); err != nil {
		t.Fatal(err)
	}
}

func TestEvaluation(t *testing.T) {
	if score := engine.Evaluate(chess.NewInitialPosition().Board); score != 0 {
		t.Fatal("Initial position is not balanced")
	}
	var white, black, edge chess.Board
	setPiece(t, &white, chess.D4, chess.Knight, chess.White)
	setPiece(t, &black, chess.D5, chess.Knight, chess.Black)
	setPiece(t, &edge, chess.A1, chess.Knight, chess.White)
	if engine.Evaluate(white) <= 0 || engine.Evaluate(black) != -engine.Evaluate(white) || engine.Evaluate(white) <= engine.Evaluate(edge) {
		t.Error("Evaluation must respect material, color symmetry and centralization")
	}
	setPiece(t, &white, chess.D4, chess.Queen, chess.White)
	if engine.Evaluate(white) <= -engine.Evaluate(black) {
		t.Error("Queen must outweigh a knight")
	}
}
