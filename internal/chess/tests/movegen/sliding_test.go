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
		{"right", chess.Direction{ColStep: 1}, nil, []chess.Square{28, 29, 30, 31}},
		{"left", chess.Direction{ColStep: -1}, nil, []chess.Square{26, 25, 24}},
		{"up", chess.Direction{RowStep: 1}, nil, []chess.Square{35, 43, 51, 59}},
		{"down", chess.Direction{RowStep: -1}, nil, []chess.Square{19, 11, 3}},
		{"own piece blocks", chess.Direction{ColStep: 1}, []placedPiece{{29, ownPawn}}, []chess.Square{28}},
		{"capture ends ray", chess.Direction{ColStep: 1}, []placedPiece{{29, enemyPawn}}, []chess.Square{28, 29}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			mustSetPiece(t, &board, 27, chess.Piece{Type: chess.Rook, Color: chess.White})
			for _, item := range tc.extra {
				mustSetPiece(t, &board, item.square, item.piece)
			}
			assertMoves(t, movegen.SlidingMoves(board, 27, tc.direction), 27, tc.destinations)
		})
	}
}

func TestRookMoves(t *testing.T) {
	runMoveCases(t, chess.Rook, movegen.RookMoves, []moveCase{
		{"center", 27, chess.White, nil, rookDestinations},
	})
}

func TestRookMovesFromEmptySquare(t *testing.T) {
	assertMoves(t, movegen.RookMoves(chess.Board{}, 27), 27, nil)
}

func TestBishopMoves(t *testing.T) {
	ownPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemyPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	runMoveCases(t, chess.Bishop, movegen.BishopMoves, []moveCase{
		{"center", 27, chess.White, nil, bishopDestinations},
		{"own piece blocks", 27, chess.White, []placedPiece{{45, ownPawn}}, []chess.Square{36, 34, 41, 48, 20, 13, 6, 18, 9, 0}},
		{"capture ends ray", 27, chess.White, []placedPiece{{45, enemyPawn}}, []chess.Square{36, 45, 34, 41, 48, 20, 13, 6, 18, 9, 0}},
	})
}

func TestQueenMoves(t *testing.T) {
	ownPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemyPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	runMoveCases(t, chess.Queen, movegen.QueenMoves, []moveCase{
		{"center", 27, chess.White, nil, queenDestinations},
		{"own piece blocks", 27, chess.White, []placedPiece{{29, ownPawn}}, []chess.Square{28, 26, 25, 24, 35, 43, 51, 59, 19, 11, 3, 36, 45, 54, 63, 34, 41, 48, 20, 13, 6, 18, 9, 0}},
		{"capture ends ray", 27, chess.White, []placedPiece{{29, enemyPawn}}, []chess.Square{28, 29, 26, 25, 24, 35, 43, 51, 59, 19, 11, 3, 36, 45, 54, 63, 34, 41, 48, 20, 13, 6, 18, 9, 0}},
	})
}
