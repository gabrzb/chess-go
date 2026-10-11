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
		{"white forward", 27, chess.White, nil, []chess.Square{35}},
		{"black forward", 27, chess.Black, nil, []chess.Square{19}},
		{"blocked in front", 27, chess.White, []placedPiece{{35, blackPawn}}, nil},
		{"white initial advance", 12, chess.White, nil, []chess.Square{20, 28}},
		{"white double advance blocked", 12, chess.White, []placedPiece{{28, blackPawn}}, []chess.Square{20}},
		{"initial advance blocked in front", 12, chess.White, []placedPiece{{20, blackPawn}}, nil},
		{"black initial advance", 52, chess.Black, nil, []chess.Square{44, 36}},
		{"white diagonal captures", 27, chess.White, []placedPiece{{35, blackPawn}, {34, blackPawn}, {36, blackPawn}}, []chess.Square{34, 36}},
		{"black diagonal captures", 27, chess.Black, []placedPiece{{19, whitePawn}, {18, whitePawn}, {20, whitePawn}}, []chess.Square{18, 20}},
		{"own pieces cannot be captured", 27, chess.White, []placedPiece{{35, whitePawn}, {34, whitePawn}, {36, whitePawn}}, nil},
		{"left edge", 24, chess.White, []placedPiece{{31, blackPawn}}, []chess.Square{32}},
		{"right edge", 31, chess.White, []placedPiece{{40, blackPawn}}, []chess.Square{39}},
	})
}
