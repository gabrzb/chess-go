package chess

import "testing"

func TestQueenMoves(t *testing.T) {
	var board Board
	queen := Piece{Type: Queen, Color: White}

	if err := board.SetPiece(Square(27), queen); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := queenMoves(board, Square(27))

	destinationSquares := []Square{
		28, 29, 30, 31, // Right
		26, 25, 24, // Left
		35, 43, 51, 59, // Up
		19, 11, 3, // Down
		36, 45, 54, 63, // Diagonal up-right
		34, 41, 48, // Diagonal up-left
		20, 13, 6, // Diagonal down-right
		18, 9, 0, // Diagonal down-left
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

func TestQueenMovesBlockedByOwnPiece(t *testing.T) {
	var board Board
	queen := Piece{Type: Queen, Color: White}
	pawn := Piece{Type: Pawn, Color: White}

	if err := board.SetPiece(Square(27), queen); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(29), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := queenMoves(board, Square(27))

	destinationSquares := []Square{
		28,         // Right
		26, 25, 24, // Left
		35, 43, 51, 59, // Up
		19, 11, 3, // Down
		36, 45, 54, 63, // Diagonal up-right
		34, 41, 48, // Diagonal up-left
		20, 13, 6, // Diagonal down-right
		18, 9, 0, // Diagonal down-left
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

func TestQueenMovesCapturesEnemy(t *testing.T) {
	var board Board
	queen := Piece{Type: Queen, Color: White}
	pawn := Piece{Type: Pawn, Color: Black}

	if err := board.SetPiece(Square(27), queen); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(29), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := queenMoves(board, Square(27))

	destinationSquares := []Square{
		28, 29, // Right
		26, 25, 24, // Left
		35, 43, 51, 59, // Up
		19, 11, 3, // Down
		36, 45, 54, 63, // Diagonal up-right
		34, 41, 48, // Diagonal up-left
		20, 13, 6, // Diagonal down-right
		18, 9, 0, // Diagonal down-left
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
