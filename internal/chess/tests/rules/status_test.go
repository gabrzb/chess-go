package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestEvaluateStatus(t *testing.T) {
	cases := []struct {
		name  string
		queen chess.Square
		want  chess.Status
	}{
		{"checkmate", 49, chess.StatusCheckmate},
		{"stalemate", 41, chess.StatusStalemate},
		{"ongoing", 40, chess.StatusOngoing},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, 56, chess.Piece{Type: chess.King, Color: chess.Black})
			mustSetPiece(t, &board, 42, chess.Piece{Type: chess.King, Color: chess.White})
			mustSetPiece(t, &board, tc.queen, chess.Piece{Type: chess.Queen, Color: chess.White})
			got, err := rules.EvaluateStatus(board, chess.Black)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Expected status %v, got %v", tc.want, got)
			}
		})
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
