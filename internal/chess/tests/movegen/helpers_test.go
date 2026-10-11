package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

type placedPiece struct {
	square chess.Square
	piece  chess.Piece
}
type moveCase struct {
	name         string
	from         chess.Square
	color        chess.Color
	extra        []placedPiece
	destinations []chess.Square
}

var (
	rookDestinations   = []chess.Square{28, 29, 30, 31, 26, 25, 24, 35, 43, 51, 59, 19, 11, 3}
	bishopDestinations = []chess.Square{36, 45, 54, 63, 34, 41, 48, 20, 13, 6, 18, 9, 0}
	queenDestinations  = []chess.Square{28, 29, 30, 31, 26, 25, 24, 35, 43, 51, 59, 19, 11, 3, 36, 45, 54, 63, 34, 41, 48, 20, 13, 6, 18, 9, 0}
	knightDestinations = []chess.Square{37, 21, 33, 17, 44, 42, 12, 10}
	kingDestinations   = []chess.Square{28, 26, 35, 19, 36, 34, 20, 18}
)

func mustSetPiece(t *testing.T, board *chess.Board, square chess.Square, piece chess.Piece) {
	t.Helper()
	if err := board.SetPiece(square, piece); err != nil {
		t.Fatalf("Failed to set piece at %v: %v", square, err)
	}
}

func assertMoves(t *testing.T, moves []chess.Move, from chess.Square, destinations []chess.Square) {
	t.Helper()
	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d: %v", len(destinations), len(moves), moves)
	}
	for i, to := range destinations {
		want := chess.Move{From: from, To: to}
		if moves[i] != want {
			t.Errorf("Move %d: expected %v, got %v", i, want, moves[i])
		}
	}
}

func runMoveCases(t *testing.T, pieceType chess.PieceType, generate func(chess.Board, chess.Square) []chess.Move, cases []moveCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, tc.from, chess.Piece{Type: pieceType, Color: tc.color})
			for _, item := range tc.extra {
				mustSetPiece(t, &board, item.square, item.piece)
			}
			assertMoves(t, generate(board, tc.from), tc.from, tc.destinations)
		})
	}
}
