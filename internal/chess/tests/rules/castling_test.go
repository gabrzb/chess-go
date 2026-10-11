package rules_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func whiteKingSidePosition(t *testing.T) chess.Position {
	t.Helper()

	var board chess.Board

	pieces := []struct {
		square chess.Square
		piece  chess.Piece
	}{
		{4, chess.Piece{Type: chess.King, Color: chess.White}},
		{7, chess.Piece{Type: chess.Rook, Color: chess.White}},
		{60, chess.Piece{Type: chess.King, Color: chess.Black}},
	}

	for _, item := range pieces {
		if err := board.SetPiece(item.square, item.piece); err != nil {
			t.Fatalf("Failed to set piece on board: %v", err)
		}
	}

	position, err := chess.NewPosition(board, chess.White)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	position.CastlingRights = chess.WhiteKingSide

	return position
}

func TestCanCastleKingSideWhiteAllowed(t *testing.T) {
	position := whiteKingSidePosition(t)

	if !rules.CanCastleKingSide(position) {
		t.Error("Expected white king-side castling to be allowed")
	}
}

func TestCanCastleKingSideWhiteWithoutRight(t *testing.T) {
	for _, rights := range []chess.CastlingRights{0, chess.WhiteQueenSide} {
		position := whiteKingSidePosition(t)
		position.CastlingRights = rights

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected white king-side castling to be disallowed with rights %v", rights)
		}
	}
}

func TestCanCastleKingSideWhiteBlocked(t *testing.T) {
	for _, square := range []chess.Square{5, 6} {
		position := whiteKingSidePosition(t)
		blocker := chess.Piece{Type: chess.Knight, Color: chess.White}

		if err := position.Board.SetPiece(square, blocker); err != nil {
			t.Fatalf("Failed to set piece on board: %v", err)
		}

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected white king-side castling to be disallowed when square %v is blocked", square)
		}
	}
}

func TestCanCastleKingSideWhiteAttackedSquares(t *testing.T) {
	cases := []struct {
		name string
		from chess.Square
	}{
		{"king in check", 52},
		{"crossing check", 53},
		{"destination check", 54},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			position := whiteKingSidePosition(t)
			attacker := chess.Piece{Type: chess.Queen, Color: chess.Black}

			if err := position.Board.SetPiece(tc.from, attacker); err != nil {
				t.Fatalf("Failed to set piece on board: %v", err)
			}

			if rules.CanCastleKingSide(position) {
				t.Errorf("Expected the attacker at %v to prevent castling", tc.from)
			}
		})
	}
}

func TestCanCastleKingSideWhiteMissingPiece(t *testing.T) {
	for _, square := range []chess.Square{4, 7} {
		position := whiteKingSidePosition(t)

		if _, err := position.Board.RemovePiece(square); err != nil {
			t.Fatalf("Failed to remove piece from board: %v", err)
		}

		if rules.CanCastleKingSide(position) {
			t.Errorf("Expected missing piece at %v to prevent castling", square)
		}
	}
}
