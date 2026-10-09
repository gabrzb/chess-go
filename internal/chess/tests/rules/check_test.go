package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestFindKing(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(60), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	cases := []struct {
		name  string
		color chess.Color
		want  chess.Square
	}{
		{"white", chess.White, chess.Square(4)},
		{"black", chess.Black, chess.Square(60)},
		{"invalid color", chess.NoColor, chess.NoSquare},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := rules.FindKing(board, tc.color)

			if got != tc.want {
				t.Errorf("FindKing(%v) = %v; want %v", tc.color, got, tc.want)
			}
		})
	}
}

func TestFindKingMissing(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.White}

	if err := board.SetPiece(chess.Square(4), rook); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	got := rules.FindKing(board, chess.White)

	if got != chess.NoSquare {
		t.Errorf("FindKing(White) = %v; want NoSquare", got)
	}
}

func TestIsInCheckWhite(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(60), chess.Piece{Type: chess.Rook, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	got, err := rules.IsInCheck(board, chess.White)
	if err != nil {
		t.Fatalf("IsInCheck returned an error: %v", err)
	}

	if !got {
		t.Errorf("IsInCheck(White) = %v; want true", got)
	}
}

func TestIsInCheckBlack(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.Rook, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white rook: %v", err)
	}

	if err := board.SetPiece(chess.Square(60), chess.Piece{Type: chess.King, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black king: %v", err)
	}

	got, err := rules.IsInCheck(board, chess.Black)
	if err != nil {
		t.Fatalf("IsInCheck returned an error: %v", err)
	}

	if !got {
		t.Errorf("IsInCheck(White) = %v; want true", got)
	}
}

func TestIsInCheckBlocked(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white king: %v", err)
	}

	if err := board.SetPiece(chess.Square(60), chess.Piece{Type: chess.Rook, Color: chess.Black}); err != nil {
		t.Fatalf("Failed to set black rook: %v", err)
	}

	if err := board.SetPiece(chess.Square(28), chess.Piece{Type: chess.Pawn, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set white pawn: %v", err)
	}

	got, err := rules.IsInCheck(board, chess.White)
	if err != nil {
		t.Fatalf("IsInCheck returned an error: %v", err)
	}

	if got {
		t.Errorf("IsInCheck(White) = %v; want false", got)
	}
}

func TestIsInCheckInvalidColor(t *testing.T) {
	var board chess.Board

	_, err := rules.IsInCheck(board, chess.NoColor)

	if !errors.Is(err, chess.ErrInvalidColor) {
		t.Errorf("IsInCheck returned error %v; want ErrInvalidColor", err)
	}
}

func TestIsInCheckMissingKing(t *testing.T) {
	var board chess.Board

	_, err := rules.IsInCheck(board, chess.White)

	if !errors.Is(err, chess.ErrKingNotFound) {
		t.Errorf("IsInCheck returned error %v; want ErrKingNotFound", err)
	}
}
