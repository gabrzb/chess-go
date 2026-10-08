package chess

import "testing"

func TestIsSquareAttackedEmptyBoard(t *testing.T) {
	var board Board

	if IsSquareAttacked(board, Square(27), Black) {
		t.Error("Expected square 27 to not be attacked on an empty board")
	}
}

func TestIsSquareAttackedByColor(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: Black}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece on board: %v", err)
	}

	if !IsSquareAttacked(board, Square(31), Black) {
		t.Error("Expected square 27 to be attacked by Black rook")
	}

	if IsSquareAttacked(board, Square(31), White) {
		t.Error("Expected square 27 to not be attacked by White rook")
	}
}

func TestIsSquareAttackedByPiece(t *testing.T) {
	cases := []struct {
		name      string
		pieceType PieceType
		from      Square
	}{
		{"pawn", Pawn, Square(34)},     // c5
		{"knight", Knight, Square(10)}, // c2
		{"bishop", Bishop, Square(0)},  // a1
		{"rook", Rook, Square(3)},      // d1
		{"queen", Queen, Square(63)},   // h8
		{"king", King, Square(26)},     // c4
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board Board
			piece := Piece{Type: tc.pieceType, Color: Black}

			if err := board.SetPiece(tc.from, piece); err != nil {
				t.Fatalf("Failed to set piece on board: %v", err)
			}

			if !IsSquareAttacked(board, Square(27), Black) {
				t.Errorf("Expected %s to attack d4", tc.name)
			}
		})
	}
}

func TestIsSquareAttackedContinuesSearching(t *testing.T) {
	var board Board
	knight := Piece{Type: Knight, Color: Black}
	rook := Piece{Type: Rook, Color: Black}

	if err := board.SetPiece(Square(0), knight); err != nil {
		t.Fatalf("Failed to set knight: %v", err)
	}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set rook: %v", err)
	}

	if !IsSquareAttacked(board, Square(31), Black) {
		t.Error("Expected the search to find the rook attacking h4")
	}
}
