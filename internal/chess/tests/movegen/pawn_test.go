package movegen_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/movegen"
)

func TestPawnMoves(t *testing.T) {
	whitePawn := chess.Piece{Type: chess.Pawn, Color: chess.White}
	blackPawn := chess.Piece{Type: chess.Pawn, Color: chess.Black}
	runMoveCases(t, chess.Pawn, movegen.PawnMoves, []moveCase{
		{"white forward", chess.D4, chess.White, nil, []chess.Square{chess.D5}},
		{"black forward", chess.D4, chess.Black, nil, []chess.Square{chess.D3}},
		{"blocked in front", chess.D4, chess.White, []placedPiece{{chess.D5, blackPawn}}, nil},
		{"white initial advance", chess.E2, chess.White, nil, []chess.Square{chess.E3, chess.E4}},
		{"white double advance blocked", chess.E2, chess.White, []placedPiece{{chess.E4, blackPawn}}, []chess.Square{chess.E3}},
		{"initial advance blocked in front", chess.E2, chess.White, []placedPiece{{chess.E3, blackPawn}}, nil},
		{"black initial advance", chess.E7, chess.Black, nil, []chess.Square{chess.E6, chess.E5}},
		{"white diagonal captures", chess.D4, chess.White, []placedPiece{{chess.D5, blackPawn}, {chess.C5, blackPawn}, {chess.E5, blackPawn}}, []chess.Square{chess.C5, chess.E5}},
		{"black diagonal captures", chess.D4, chess.Black, []placedPiece{{chess.D3, whitePawn}, {chess.C3, whitePawn}, {chess.E3, whitePawn}}, []chess.Square{chess.C3, chess.E3}},
		{"own pieces cannot be captured", chess.D4, chess.White, []placedPiece{{chess.D5, whitePawn}, {chess.C5, whitePawn}, {chess.E5, whitePawn}}, nil},
		{"left edge", chess.A4, chess.White, []placedPiece{{chess.H4, blackPawn}}, []chess.Square{chess.A5}},
		{"right edge", chess.H4, chess.White, []placedPiece{{chess.A6, blackPawn}}, []chess.Square{chess.H5}},
	})
}
