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
			mustSetPiece(t, &board, 27, rook)
			mustSetPiece(t, &board, 35, tc.target)
			before := board
			next, err := rules.BoardAfterMove(board, chess.Move{From: 27, To: 35})
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			assertPositionPiece(t, chess.Position{Board: next}, 27, chess.Piece{})
			assertPositionPiece(t, chess.Position{Board: next}, 35, rook)
			if board != before {
				t.Error("Expected original board to remain unchanged")
			}
		})
	}
}

func TestBoardAfterMoveErrors(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.Square(27), chess.Piece{Type: chess.Rook, Color: chess.White})

	cases := []struct {
		name    string
		move    chess.Move
		wantErr error
	}{
		{"invalid source", chess.Move{From: chess.NoSquare, To: chess.Square(35)}, chess.ErrInvalidSquare},
		{"invalid destination", chess.Move{From: chess.Square(27), To: chess.NoSquare}, chess.ErrInvalidSquare},
		{"empty source", chess.Move{From: chess.Square(28), To: chess.Square(35)}, chess.ErrEmptySource},
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
