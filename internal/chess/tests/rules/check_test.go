package rules_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func TestFindKing(t *testing.T) {
	var board chess.Board

	mustSetPiece(t, &board, chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White})

	mustSetPiece(t, &board, chess.Square(60), chess.Piece{Type: chess.King, Color: chess.Black})

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

	mustSetPiece(t, &board, chess.Square(4), rook)

	got := rules.FindKing(board, chess.White)

	if got != chess.NoSquare {
		t.Errorf("FindKing(White) = %v; want NoSquare", got)
	}
}

func TestIsInCheck(t *testing.T) {
	cases := []struct {
		name                   string
		color, enemy           chess.Color
		kingFrom, attackerFrom chess.Square
		blocked, want          bool
		wantErr                error
	}{
		{"white", chess.White, chess.Black, 4, 60, false, true, nil},
		{"black", chess.Black, chess.White, 60, 4, false, true, nil},
		{"blocked", chess.White, chess.Black, 4, 60, true, false, nil},
		{"invalid color", chess.NoColor, chess.NoColor, chess.NoSquare, chess.NoSquare, false, false, chess.ErrInvalidColor},
		{"missing king", chess.White, chess.NoColor, chess.NoSquare, chess.NoSquare, false, false, chess.ErrKingNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			if tc.kingFrom.Valid() {
				mustSetPiece(t, &board, tc.kingFrom, chess.Piece{Type: chess.King, Color: tc.color})
			}
			if tc.attackerFrom.Valid() {
				mustSetPiece(t, &board, tc.attackerFrom, chess.Piece{Type: chess.Rook, Color: tc.enemy})
			}
			if tc.blocked {
				mustSetPiece(t, &board, 28, chess.Piece{Type: chess.Pawn, Color: chess.White})
			}
			got, err := rules.IsInCheck(board, tc.color)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("Expected error %v, got %v", tc.wantErr, err)
			}
			if err == nil && got != tc.want {
				t.Errorf("Expected check=%v, got %v", tc.want, got)
			}
		})
	}
}
