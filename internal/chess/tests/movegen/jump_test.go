package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

func TestKnightMoves(t *testing.T) {
	ownPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemyPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	runMoveCases(t, chess.Knight, movegen.KnightMoves, []moveCase{
		{"center", 27, chess.White, nil, knightDestinations},
		{"corner", 0, chess.White, nil, []chess.Square{10, 17}},
		{"jump over pieces", 27, chess.White, []placedPiece{{18, ownPawn}, {19, ownPawn}, {20, ownPawn}, {26, ownPawn}, {28, ownPawn}, {34, ownPawn}, {35, ownPawn}, {36, ownPawn}}, knightDestinations},
		{"own piece blocks landing", 27, chess.White, []placedPiece{{37, ownPawn}}, []chess.Square{21, 33, 17, 44, 42, 12, 10}},
		{"capture enemy", 27, chess.White, []placedPiece{{37, enemyPawn}}, knightDestinations},
	})
}

func TestKingMoves(t *testing.T) {
	ownPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemyPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	runMoveCases(t, chess.King, movegen.KingMoves, []moveCase{
		{"center", 27, chess.White, nil, kingDestinations},
		{"corner", 0, chess.White, nil, []chess.Square{1, 8, 9}},
		{"own piece blocks landing", 27, chess.White, []placedPiece{{28, ownPawn}}, []chess.Square{26, 35, 19, 36, 34, 20, 18}},
		{"capture enemy", 27, chess.White, []placedPiece{{28, enemyPawn}}, kingDestinations},
	})
}
