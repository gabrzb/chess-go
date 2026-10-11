package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

func TestCanOccupySquare(t *testing.T) {
	cases := []struct {
		name   string
		square chess.Square
		piece  chess.Piece
		want   bool
	}{
		{"empty", chess.D4, chess.Piece{}, true},
		{"invalid", chess.NoSquare, chess.Piece{}, false},
		{"own piece", chess.D4, chess.Piece{Type: chess.Pawn, Color: chess.White}, false},
		{"enemy piece", chess.D4, chess.Piece{Type: chess.Pawn, Color: chess.Black}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			if tc.square.Valid() {
				mustSetPiece(t, &board, tc.square, tc.piece)
			}
			if got := movegen.CanOccupySquare(board, tc.square, chess.White); got != tc.want {
				t.Errorf("Expected occupancy permission %v, got %v", tc.want, got)
			}
		})
	}
}

func TestJumpMoves(t *testing.T) {
	var board chess.Board
	mustSetPiece(t, &board, chess.D4, chess.Piece{Type: chess.Knight, Color: chess.White})
	assertMoves(t, movegen.JumpMoves(board, chess.D4, chess.Direction{ColStep: 2, RowStep: 1}), chess.D4, []chess.Square{chess.F5})
}

func TestPseudoLegalMovesFromUnavailableSquare(t *testing.T) {
	for _, from := range []chess.Square{chess.D4, chess.NoSquare} {
		t.Run(from.String(), func(t *testing.T) {
			assertMoves(t, movegen.PseudoLegalMoves(chess.Board{}, from), from, nil)
		})
	}
}

func TestPseudoLegalMovesByPiece(t *testing.T) {
	cases := []struct {
		name         string
		pieceType    chess.PieceType
		destinations []chess.Square
	}{
		{"pawn", chess.Pawn, []chess.Square{chess.D5}},
		{"knight", chess.Knight, knightDestinations},
		{"bishop", chess.Bishop, bishopDestinations},
		{"rook", chess.Rook, rookDestinations},
		{"queen", chess.Queen, queenDestinations},
		{"king", chess.King, kingDestinations},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, chess.D4, chess.Piece{Type: tc.pieceType, Color: chess.White})
			assertMoves(t, movegen.PseudoLegalMoves(board, chess.D4), chess.D4, tc.destinations)
		})
	}
}
