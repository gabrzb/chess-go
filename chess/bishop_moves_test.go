package chess

import "testing"

func TestBishopMoves(t *testing.T) {
	var board Board
	bishop := Piece{Type: Bishop, Color: White}

	if err := board.SetPiece(Square(27), bishop); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := bishopMoves(board, Square(27))

	destinationSquares := []Square{
		36, 45, 54, 63, // diagonal up-right
		34, 41, 48, // diagonal up-left
		20, 13, 6, // diagonal down-right
		18, 9, 0, // diagonal down-left
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

func TestBishopMovesBlockedByOwnPiece(t *testing.T) {
	var board Board
	bishop := Piece{Type: Bishop, Color: White}
	pawn := Piece{Type: Pawn, Color: White}

	if err := board.SetPiece(Square(27), bishop); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(45), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := bishopMoves(board, Square(27))

	destinationSquares := []Square{
		36,         // diagonal up-right
		34, 41, 48, // diagonal up-left
		20, 13, 6, // diagonal down-right
		18, 9, 0, // diagonal down-left
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

func TestBishopCapturesEnemy(t *testing.T) {
	var board Board
	bishop := Piece{Type: Bishop, Color: White}
	pawn := Piece{Type: Pawn, Color: Black}

	if err := board.SetPiece(Square(27), bishop); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(45), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := bishopMoves(board, Square(27))

	destinationSquares := []Square{
		36, 45, // diagonal up-right
		34, 41, 48, // diagonal up-left
		20, 13, 6, // diagonal down-right
		18, 9, 0, // diagonal down-left
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
