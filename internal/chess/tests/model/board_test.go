package model_test

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestPieceAtEmptyBoard(t *testing.T) {
	var board chess.Board
	piece, ok := board.PieceAt(chess.Square(0))

	if piece != (chess.Piece{}) || !ok {
		t.Errorf("PieceAt(0) = %v, %v; expected Piece{} and true", piece, ok)
	}

	piece, ok = board.PieceAt(chess.NoSquare)

	if piece != (chess.Piece{}) || ok {
		t.Errorf("PieceAt(NoSquare) = %v, %v; expected Piece{} and false", piece, ok)
	}
}

func TestSetPiece(t *testing.T) {
	var board chess.Board
	piece := chess.Piece{Type: chess.Queen, Color: chess.White}
	if err := board.SetPiece(chess.Square(0), piece); err != nil {
		t.Fatalf("SetPiece(0, piece) = %v; expected nil", err)
	}

	got, ok := board.PieceAt(chess.Square(0))
	if got != piece || !ok {
		t.Errorf("PieceAt(0) = %v, %v; expected %v and true", got, ok, piece)
	}
}

func TestSetPieceInvalidSquare(t *testing.T) {
	var board chess.Board
	piece := chess.Piece{Type: chess.Queen, Color: chess.White}
	setPiece := board.SetPiece(chess.NoSquare, piece)

	if !errors.Is(setPiece, chess.ErrInvalidSquare) {
		t.Errorf("SetPiece(NoSquare, piece) = %v; expected ErrInvalidSquare", setPiece)
	}
}

func TestRemovePiece(t *testing.T) {
	var board chess.Board
	piece := chess.Piece{Type: chess.Queen, Color: chess.White}

	if err := board.SetPiece(chess.Square(0), piece); err != nil {
		t.Fatalf("SetPiece(0, piece) = %v; expected nil", err)
	}

	removed, err := board.RemovePiece(chess.Square(0))
	if removed != piece || err != nil {
		t.Errorf("RemovePiece(0) = %v, %v; expected %v and nil", removed, err, piece)
	}

	squarePiece, ok := board.PieceAt(chess.Square(0))

	if squarePiece != (chess.Piece{}) || !ok {
		t.Errorf("PieceAt(0) after removal = %v, %v; expected Piece{} and true", squarePiece, ok)
	}

	removedNoSquare, errNoSquare := board.RemovePiece(chess.NoSquare)
	if removedNoSquare != (chess.Piece{}) || !errors.Is(errNoSquare, chess.ErrInvalidSquare) {
		t.Errorf("RemovePiece(NoSquare) = %v, %v; expected Piece{} and ErrInvalidSquare", removedNoSquare, errNoSquare)
	}
}
