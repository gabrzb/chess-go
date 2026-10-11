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
		{"center", chess.D4, chess.White, nil, knightDestinations},
		{"corner", chess.A1, chess.White, nil, []chess.Square{chess.C2, chess.B3}},
		{"jump over pieces", chess.D4, chess.White, []placedPiece{{chess.C3, ownPawn}, {chess.D3, ownPawn}, {chess.E3, ownPawn}, {chess.C4, ownPawn}, {chess.E4, ownPawn}, {chess.C5, ownPawn}, {chess.D5, ownPawn}, {chess.E5, ownPawn}}, knightDestinations},
		{"own piece blocks landing", chess.D4, chess.White, []placedPiece{{chess.F5, ownPawn}}, []chess.Square{chess.F3, chess.B5, chess.B3, chess.E6, chess.C6, chess.E2, chess.C2}},
		{"capture enemy", chess.D4, chess.White, []placedPiece{{chess.F5, enemyPawn}}, knightDestinations},
	})
}

func TestKingMoves(t *testing.T) {
	ownPawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	enemyPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	runMoveCases(t, chess.King, movegen.KingMoves, []moveCase{
		{"center", chess.D4, chess.White, nil, kingDestinations},
		{"corner", chess.A1, chess.White, nil, []chess.Square{chess.B1, chess.A2, chess.B2}},
		{"own piece blocks landing", chess.D4, chess.White, []placedPiece{{chess.E4, ownPawn}}, []chess.Square{chess.C4, chess.D5, chess.D3, chess.E5, chess.C5, chess.E3, chess.C3}},
		{"capture enemy", chess.D4, chess.White, []placedPiece{{chess.E4, enemyPawn}}, kingDestinations},
	})
}
