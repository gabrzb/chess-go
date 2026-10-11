package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

func TestSlidingMoves(t *testing.T) {
	ownPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemyPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	cases := []struct {
		name         string
		direction    chess.Direction
		extra        []placedPiece
		destinations []chess.Square
	}{
		{"right", chess.Direction{ColStep: 1}, nil, []chess.Square{chess.E4, chess.F4, chess.G4, chess.H4}},
		{"left", chess.Direction{ColStep: -1}, nil, []chess.Square{chess.C4, chess.B4, chess.A4}},
		{"up", chess.Direction{RowStep: 1}, nil, []chess.Square{chess.D5, chess.D6, chess.D7, chess.D8}},
		{"down", chess.Direction{RowStep: -1}, nil, []chess.Square{chess.D3, chess.D2, chess.D1}},
		{"own piece blocks", chess.Direction{ColStep: 1}, []placedPiece{{chess.F4, ownPawn}}, []chess.Square{chess.E4}},
		{"capture ends ray", chess.Direction{ColStep: 1}, []placedPiece{{chess.F4, enemyPawn}}, []chess.Square{chess.E4, chess.F4}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, chess.D4, chess.Piece{Type: chess.Rook, Color: chess.White})
			for _, item := range tc.extra {
				mustSetPiece(t, &board, item.square, item.piece)
			}
			assertMoves(t, movegen.SlidingMoves(board, chess.D4, tc.direction), chess.D4, tc.destinations)
		})
	}
}

func TestRookMoves(t *testing.T) {
	runMoveCases(t, chess.Rook, movegen.RookMoves, []moveCase{
		{"center", chess.D4, chess.White, nil, rookDestinations},
	})
}

func TestRookMovesFromEmptySquare(t *testing.T) {
	assertMoves(t, movegen.RookMoves(chess.Board{}, chess.D4), chess.D4, nil)
}

func TestBishopMoves(t *testing.T) {
	ownPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemyPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	runMoveCases(t, chess.Bishop, movegen.BishopMoves, []moveCase{
		{"center", chess.D4, chess.White, nil, bishopDestinations},
		{"own piece blocks", chess.D4, chess.White, []placedPiece{{chess.F6, ownPawn}}, []chess.Square{chess.E5, chess.C5, chess.B6, chess.A7, chess.E3, chess.F2, chess.G1, chess.C3, chess.B2, chess.A1}},
		{"capture ends ray", chess.D4, chess.White, []placedPiece{{chess.F6, enemyPawn}}, []chess.Square{chess.E5, chess.F6, chess.C5, chess.B6, chess.A7, chess.E3, chess.F2, chess.G1, chess.C3, chess.B2, chess.A1}},
	})
}

func TestQueenMoves(t *testing.T) {
	ownPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemyPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	runMoveCases(t, chess.Queen, movegen.QueenMoves, []moveCase{
		{"center", chess.D4, chess.White, nil, queenDestinations},
		{"own piece blocks", chess.D4, chess.White, []placedPiece{{chess.F4, ownPawn}}, []chess.Square{chess.E4, chess.C4, chess.B4, chess.A4, chess.D5, chess.D6, chess.D7, chess.D8, chess.D3, chess.D2, chess.D1, chess.E5, chess.F6, chess.G7, chess.H8, chess.C5, chess.B6, chess.A7, chess.E3, chess.F2, chess.G1, chess.C3, chess.B2, chess.A1}},
		{"capture ends ray", chess.D4, chess.White, []placedPiece{{chess.F4, enemyPawn}}, []chess.Square{chess.E4, chess.F4, chess.C4, chess.B4, chess.A4, chess.D5, chess.D6, chess.D7, chess.D8, chess.D3, chess.D2, chess.D1, chess.E5, chess.F6, chess.G7, chess.H8, chess.C5, chess.B6, chess.A7, chess.E3, chess.F2, chess.G1, chess.C3, chess.B2, chess.A1}},
	})
}
