package attacks_test

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/attacks"
)

func TestIsSquareAttackedEmptyBoard(t *testing.T) {
	var board chess.Board

	if attacks.IsSquareAttacked(board, chess.Square(27), chess.Black) {
		t.Error("Expected square 27 to not be attacked on an empty board")
	}
}

func TestIsSquareAttackedByColor(t *testing.T) {
	var board chess.Board
	rook := chess.Piece{Type: chess.Rook, Color: chess.Black}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece on board: %v", err)
	}

	if !attacks.IsSquareAttacked(board, chess.Square(31), chess.Black) {
		t.Error("Expected square 27 to be attacked by Black rook")
	}

	if attacks.IsSquareAttacked(board, chess.Square(31), chess.White) {
		t.Error("Expected square 27 to not be attacked by White rook")
	}
}

func TestIsSquareAttackedByPiece(t *testing.T) {
	cases := []struct {
		name      string
		pieceType chess.PieceType
		from      chess.Square
	}{
		{"pawn", chess.Pawn, chess.Square(34)},     // c5
		{"knight", chess.Knight, chess.Square(10)}, // c2
		{"bishop", chess.Bishop, chess.Square(0)},  // a1
		{"rook", chess.Rook, chess.Square(3)},      // d1
		{"queen", chess.Queen, chess.Square(63)},   // h8
		{"king", chess.King, chess.Square(26)},     // c4
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var board chess.Board
			piece := chess.Piece{Type: tc.pieceType, Color: chess.Black}

			if err := board.SetPiece(tc.from, piece); err != nil {
				t.Fatalf("Failed to set piece on board: %v", err)
			}

			if !attacks.IsSquareAttacked(board, chess.Square(27), chess.Black) {
				t.Errorf("Expected %s to attack d4", tc.name)
			}
		})
	}
}

func TestIsSquareAttackedContinuesSearching(t *testing.T) {
	var board chess.Board
	knight := chess.Piece{Type: chess.Knight, Color: chess.Black}
	rook := chess.Piece{Type: chess.Rook, Color: chess.Black}

	if err := board.SetPiece(chess.Square(0), knight); err != nil {
		t.Fatalf("Failed to set knight: %v", err)
	}

	if err := board.SetPiece(chess.Square(27), rook); err != nil {
		t.Fatalf("Failed to set rook: %v", err)
	}

	if !attacks.IsSquareAttacked(board, chess.Square(31), chess.Black) {
		t.Error("Expected the search to find the rook attacking h4")
	}
}
