package model_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestNewPosition(t *testing.T) {
	var board chess.Board

	if err := board.SetPiece(chess.Square(4), chess.Piece{Type: chess.King, Color: chess.White}); err != nil {
		t.Fatalf("Failed to set king: %v", err)
	}

	position, err := chess.NewPosition(board, chess.White)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if position.Board != board {
		t.Error("Expected the position to preserve the supplied board")
	}

	if position.SideToMove != chess.White {
		t.Error("Expected White to move")
	}

	if position.EnPassant != chess.NoSquare {
		t.Error("Expected no en passant target")
	}

	if position.CastlingRights != 0 {
		t.Error("Expected no castling rights")
	}

	if position.HalfmoveClock != 0 {
		t.Error("Expected halfmove clock to start at 0")
	}

	if position.FullmoveNumber != 1 {
		t.Error("Expected fullmove number to start at 1")
	}
}

func TestNewPositionBlackToMove(t *testing.T) {
	var board chess.Board

	position, err := chess.NewPosition(board, chess.Black)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if position.SideToMove != chess.Black {
		t.Error("Expected Black to move")
	}
}

func TestNewPositionInvalidColor(t *testing.T) {
	var board chess.Board

	_, err := chess.NewPosition(board, chess.NoColor)

	if !errors.Is(err, chess.ErrInvalidColor) {
		t.Errorf("Expected ErrInvalidColor, got %v", err)
	}
}
