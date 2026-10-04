package chess

import "testing"

func TestKingMoves(t *testing.T) {
	var board Board
	king := Piece{Type: King, Color: White}

	if err := board.SetPiece(Square(27), king); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := kingMoves(board, Square(27))

	destinationSquares := []Square{
		28, 26, // e4, c4
		35, 19, // d5, d3
		36, 34, // e5, c5
		20, 18, // e3, c3
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := Move{From: Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKingMovesFromCorner(t *testing.T) {
	var board Board
	king := Piece{Type: King, Color: White}

	if err := board.SetPiece(Square(0), king); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := kingMoves(board, Square(0))

	destinationSquares := []Square{1, 8, 9} // b1, a2, b2

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := Move{From: Square(0), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKingMovesBlockedByOwnPiece(t *testing.T) {
	var board Board
	king := Piece{Type: King, Color: White}
	pawn := Piece{Type: Pawn, Color: White}

	if err := board.SetPiece(Square(27), king); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(28), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := kingMoves(board, Square(27))

	destinationSquares := []Square{
		26,     // c4
		35, 19, // d5, d3
		36, 34, // e5, c5
		20, 18, // e3, c3
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := Move{From: Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestKingMovesCapturesEnemy(t *testing.T) {
	var board Board
	king := Piece{Type: King, Color: White}
	pawn := Piece{Type: Pawn, Color: Black}

	if err := board.SetPiece(Square(27), king); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(28), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := kingMoves(board, Square(27))

	destinationSquares := []Square{
		28, 26, // e4, c4
		35, 19, // d5, d3
		36, 34, // e5, c5
		20, 18, // e3, c3
	}

	if len(moves) != len(destinationSquares) {
		t.Fatalf("Expected %d moves, got %d", len(destinationSquares), len(moves))
	}

	for i, destination := range destinationSquares {
		expectedMove := Move{From: Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}
