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
	rookDestinations   = []chess.Square{chess.E4, chess.F4, chess.G4, chess.H4, chess.C4, chess.B4, chess.A4, chess.D5, chess.D6, chess.D7, chess.D8, chess.D3, chess.D2, chess.D1}
	bishopDestinations = []chess.Square{chess.E5, chess.F6, chess.G7, chess.H8, chess.C5, chess.B6, chess.A7, chess.E3, chess.F2, chess.G1, chess.C3, chess.B2, chess.A1}
	queenDestinations  = []chess.Square{chess.E4, chess.F4, chess.G4, chess.H4, chess.C4, chess.B4, chess.A4, chess.D5, chess.D6, chess.D7, chess.D8, chess.D3, chess.D2, chess.D1, chess.E5, chess.F6, chess.G7, chess.H8, chess.C5, chess.B6, chess.A7, chess.E3, chess.F2, chess.G1, chess.C3, chess.B2, chess.A1}
	knightDestinations = []chess.Square{chess.F5, chess.F3, chess.B5, chess.B3, chess.E6, chess.C6, chess.E2, chess.C2}
	kingDestinations   = []chess.Square{chess.E4, chess.C4, chess.D5, chess.D3, chess.E5, chess.C5, chess.E3, chess.C3}
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
