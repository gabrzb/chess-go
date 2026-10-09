package chess

import (
	"errors"
	"testing"
)

func TestEvaluateStatusCheckmate(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(56), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(Square(42), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(49), Piece{Type: Queen, Color: White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	got, err := EvaluateStatus(board, Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if got != StatusCheckmate {
		t.Errorf("Expected StatusCheckmate, got %v", got)
	}
}

func TestEvaluateStatusStalemate(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(56), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(Square(42), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(41), Piece{Type: Queen, Color: White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	got, err := EvaluateStatus(board, Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if got != StatusStalemate {
		t.Errorf("Expected StatusStalemate, got %v", got)
	}
}

func TestEvaluateStatusOngoing(t *testing.T) {
	var board Board

	if err := board.SetPiece(Square(56), Piece{Type: King, Color: Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	if err := board.SetPiece(Square(42), Piece{Type: King, Color: White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(Square(40), Piece{Type: Queen, Color: White}); err != nil {
		t.Fatalf("Failed to set white queen: %v", err)
	}

	got, err := EvaluateStatus(board, Black)
	if err != nil {
		t.Fatalf("Unexpected error for color Black: %v", err)
	}

	if got != StatusOngoing {
		t.Errorf("Expected StatusOngoing, got %v", got)
	}
}

func TestEvaluateStatusErrors(t *testing.T) {
	var board Board

	cases := []struct {
		name    string
		color   Color
		wantErr error
	}{
		{"invalid color", NoColor, ErrInvalidColor},
		{"missing king", White, ErrKingNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := EvaluateStatus(board, tc.color)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
