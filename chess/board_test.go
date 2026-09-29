package chess

import "testing"

func TestPieceAtEmptyBoard(t *testing.T) {
	var board Board
	piece, ok := board.PieceAt(Square(0))

	if piece != (Piece{}) || !ok {
		t.Errorf("PieceAt(0) = %v, %v; expected Piece{} and true", piece, ok)
	}

	piece, ok = board.PieceAt(NoSquare)

	if piece != (Piece{}) || ok {
		t.Errorf("PieceAt(NoSquare) = %v, %v; expected Piece{} and false", piece, ok)
	}
}

func TestSetPiece(t *testing.T) {
	var board Board
	piece := Piece{Type: Queen, Color: White}
	setPiece := board.SetPiece(Square(0), piece)
	got, ok := board.PieceAt(Square(0))

	if setPiece != nil {
		t.Errorf("SetPiece(0, piece) = %v; expected nil", setPiece)
	}

	if got != piece || !ok {
		t.Errorf("PieceAt(0) = %v, %v; expected %v and true", got, ok, piece)
	}
}

func TestSetPieceInvalidSquare(t *testing.T) {
	var board Board
	piece := Piece{Type: Queen, Color: White}
	setPiece := board.SetPiece(NoSquare, piece)

	if setPiece == nil {
		t.Errorf("SetPiece(NoSquare, piece) = %v; expected error", setPiece)
	}
}

func TestRemovePiece(t *testing.T) {
	var board Board
	piece := Piece{Type: Queen, Color: White}

	if err := board.SetPiece(Square(0), piece); err != nil {
		t.Fatalf("SetPiece(0, piece) = %v; expected nil", err)
	}
	removed, err := board.RemovePiece(Square(0))

	if removed != piece || err != nil {
		t.Errorf("RemovePiece(0) = %v, %v; expected %v and nil", removed, err, piece)
	}

	squarePiece, ok := board.PieceAt(Square(0))

	if squarePiece != (Piece{}) || !ok {
		t.Errorf("PieceAt(0) after removal = %v, %v; expected Piece{} and true", squarePiece, ok)
	}

	removedNoSquare, errNoSquare := board.RemovePiece(NoSquare)

	if removedNoSquare != (Piece{}) || errNoSquare == nil {
		t.Errorf("RemovePiece(NoSquare) = %v, %v; expected Piece{} and error", removedNoSquare, errNoSquare)
	}
}
