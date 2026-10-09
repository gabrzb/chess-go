package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestEvaluateStatusCheckmate(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(56), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(chess.Square(42), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(49), chess.Piece{Type: chess.Queen, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	got, err := rules.EvaluateStatus(board, chess.Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if got != chess.StatusCheckmate {
		t.Errorf("Expected StatusCheckmate, got %v", got)
	}
}

func TestEvaluateStatusStalemate(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(56), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(chess.Square(42), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(41), chess.Piece{Type: chess.Queen, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	got, err := rules.EvaluateStatus(board, chess.Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if got != chess.StatusStalemate {
		t.Errorf("Expected StatusStalemate, got %v", got)
	}
}

func TestEvaluateStatusOngoing(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(56), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(chess.Square(42), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(40), chess.Piece{Type: chess.Queen, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	got, err := rules.EvaluateStatus(board, chess.Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if got != chess.StatusOngoing {
		t.Errorf("Expected StatusOngoing, got %v", got)
	}
}

func TestEvaluateStatusErrors(t *testing.T) {
	var board chess.Board

	cases := []struct {
		name    string
		color   chess.Color
		wantErr error
	}{
		{"invalid color", chess.NoColor, chess.ErrInvalidColor},
		{"missing king", chess.White, chess.ErrKingNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := rules.EvaluateStatus(board, tc.color)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
