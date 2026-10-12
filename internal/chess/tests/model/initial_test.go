package model_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestNewInitialPositionState(t *testing.T) {
	position := chess.NewInitialPosition()
	rights := chess.WhiteKingSide | chess.WhiteQueenSide |
		chess.BlackKingSide | chess.BlackQueenSide

	if position.SideToMove != chess.White {
		t.Errorf("Expected White to move, got %v", position.SideToMove)
	}
	if position.CastlingRights != rights {
		t.Errorf("Expected all castling rights %v, got %v", rights, position.CastlingRights)
	}
	if position.EnPassant != chess.NoSquare {
		t.Errorf("Expected no en passant opportunity, got %v", position.EnPassant)
	}
	if position.HalfmoveClock != 0 {
		t.Errorf("Expected halfmove clock zero, got %d", position.HalfmoveClock)
	}
	if position.FullmoveNumber != 1 {
		t.Errorf("Expected fullmove number one, got %d", position.FullmoveNumber)
	}
}

func TestNewInitialPositionPieces(t *testing.T) {
	position := chess.NewInitialPosition()
	want := map[chess.Square]chess.Piece{
		chess.A1: {Type: chess.Rook, Color: chess.White},
		chess.B1: {Type: chess.Knight, Color: chess.White},
		chess.C1: {Type: chess.Bishop, Color: chess.White},
		chess.D1: {Type: chess.Queen, Color: chess.White},
		chess.E1: {Type: chess.King, Color: chess.White},
		chess.F1: {Type: chess.Bishop, Color: chess.White},
		chess.G1: {Type: chess.Knight, Color: chess.White},
		chess.H1: {Type: chess.Rook, Color: chess.White},
		chess.A2: {Type: chess.Pawn, Color: chess.White},
		chess.B2: {Type: chess.Pawn, Color: chess.White},
		chess.C2: {Type: chess.Pawn, Color: chess.White},
		chess.D2: {Type: chess.Pawn, Color: chess.White},
		chess.E2: {Type: chess.Pawn, Color: chess.White},
		chess.F2: {Type: chess.Pawn, Color: chess.White},
		chess.G2: {Type: chess.Pawn, Color: chess.White},
		chess.H2: {Type: chess.Pawn, Color: chess.White},
		chess.A7: {Type: chess.Pawn, Color: chess.Black},
		chess.B7: {Type: chess.Pawn, Color: chess.Black},
		chess.C7: {Type: chess.Pawn, Color: chess.Black},
		chess.D7: {Type: chess.Pawn, Color: chess.Black},
		chess.E7: {Type: chess.Pawn, Color: chess.Black},
		chess.F7: {Type: chess.Pawn, Color: chess.Black},
		chess.G7: {Type: chess.Pawn, Color: chess.Black},
		chess.H7: {Type: chess.Pawn, Color: chess.Black},
		chess.A8: {Type: chess.Rook, Color: chess.Black},
		chess.B8: {Type: chess.Knight, Color: chess.Black},
		chess.C8: {Type: chess.Bishop, Color: chess.Black},
		chess.D8: {Type: chess.Queen, Color: chess.Black},
		chess.E8: {Type: chess.King, Color: chess.Black},
		chess.F8: {Type: chess.Bishop, Color: chess.Black},
		chess.G8: {Type: chess.Knight, Color: chess.Black},
		chess.H8: {Type: chess.Rook, Color: chess.Black},
	}
	for square := chess.A1; square <= chess.H8; square++ {
		t.Run(square.String(), func(t *testing.T) {
			got, valid := position.Board.PieceAt(square)
			if !valid {
				t.Fatalf("Expected valid initial square %v", square)
			}
			if expected := want[square]; got != expected {
				t.Errorf("At %v: expected %v, got %v", square, expected, got)
			}
		})
	}
}
