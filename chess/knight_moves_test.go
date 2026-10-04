package chess

import "testing"

func TestKnightMoves(t *testing.T) {
	var board Board
	knight := Piece{Type: Knight, Color: White}

	if err := board.SetPiece(Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := knightMoves(board, Square(27))

	destinationSquares := []Square{
		37, 21, // f5, f3
		33, 17, // b5, b3
		44, 42, // e6, c6
		12, 10, // e2, c2
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

func TestKnightMovesFromCorner(t *testing.T) {
	var board Board
	knight := Piece{Type: Knight, Color: White}

	if err := board.SetPiece(Square(0), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := knightMoves(board, Square(0))

	destinationSquares := []Square{10, 17} // c2, b3

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

func TestKnightMovesJumpOverPieces(t *testing.T) {
	var board Board
	knight := Piece{Type: Knight, Color: White}
	pawn := Piece{Type: Pawn, Color: White}

	if err := board.SetPiece(Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	for _, square := range []Square{18, 19, 20, 26, 28, 34, 35, 36} {
		if err := board.SetPiece(square, pawn); err != nil {
			t.Fatalf("Failed to set piece: %v", err)
		}
	}

	moves := knightMoves(board, Square(27))

	destinationSquares := []Square{
		37, 21, // f5, f3
		33, 17, // b5, b3
		44, 42, // e6, c6
		12, 10, // e2, c2
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

func TestKnightMovesBlockedByOwnPiece(t *testing.T) {
	var board Board
	knight := Piece{Type: Knight, Color: White}
	pawn := Piece{Type: Pawn, Color: White}

	if err := board.SetPiece(Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(37), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := knightMoves(board, Square(27))

	destinationSquares := []Square{
		21,     // f3
		33, 17, // b5, b3
		44, 42, // e6, c6
		12, 10, // e2, c2
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

func TestKnightMovesCapturesEnemy(t *testing.T) {
	var board Board
	knight := Piece{Type: Knight, Color: White}
	pawn := Piece{Type: Pawn, Color: Black}

	if err := board.SetPiece(Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(37), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := knightMoves(board, Square(27))

	destinationSquares := []Square{
		37, 21, // f5, f3
		33, 17, // b5, b3
		44, 42, // e6, c6
		12, 10, // e2, c2
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
