package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestBoardAfterMove(t *testing.T) {
	for _, tc := range []struct {
		name   string
		target chess.Piece
	}{
		{"empty destination", chess.Piece{}},
		{"capture", chess.Piece{Type: chess.Pawn, Color: chess.Black}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			rook := chess.Piece{Type: chess.Rook, Color: chess.White}
			mustSetPiece(t, &board, chess.D4, rook)
			mustSetPiece(t, &board, chess.D5, tc.target)
			before := board
			next, err := rules.BoardAfterMove(board, chess.Move{From: chess.D4, To: chess.D5})
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			assertPositionPiece(t, chess.Position{Board: next}, chess.D4, chess.Piece{})
			assertPositionPiece(t, chess.Position{Board: next}, chess.D5, rook)
			if board != before {
				t.Error("Expected original board to remain unchanged")
			}
		})
	}
}

func TestBoardAfterMoveErrors(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.D4, chess.Piece{Type: chess.Rook, Color: chess.White})

	cases := []struct {
		name    string
		move    chess.Move
		wantErr error
	}{
		{"invalid source", chess.Move{From: chess.NoSquare, To: chess.D5}, chess.ErrInvalidSquare},
		{"invalid destination", chess.Move{From: chess.D4, To: chess.NoSquare}, chess.ErrInvalidSquare},
		{"empty source", chess.Move{From: chess.E4, To: chess.D5}, chess.ErrEmptySource},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := rules.BoardAfterMove(board, tc.move)

			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}
