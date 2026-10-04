package chess

import "testing"

func checkRookMoves(t *testing.T, board Board, rook Piece) {
	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := rookMoves(board, Square(27))

	if len(moves) != 4 {
		t.Fatalf("Expected 4 moves, got %d", len(moves))
	}

	destinations := []Square{28, 29, 30, 31}

	for i, destination := range destinations {
		expectedMove := Move{From: Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestRookMovesToRight(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}

	checkRookMoves(t, board, rook)
}

func TestRookMovesToRightBlockedByOwnPiece(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}
	pawn := Piece{Type: Pawn, Color: White}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(29), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := rookMoves(board, Square(27))

	if len(moves) != 1 {
		t.Fatalf("Expected 1 move, got %d", len(moves))
	}

	expectedMove := Move{From: Square(27), To: Square(28)}
	if moves[0] != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, moves[0])
	}
}

func TestRookMovesToRightCapturesEnemy(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}
	pawn := Piece{Type: Pawn, Color: Black}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if err := board.SetPiece(Square(29), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := rookMoves(board, Square(27))
	expectedMoves := []Move{
		{From: Square(27), To: Square(28)},
		{From: Square(27), To: Square(29)},
	}

	if len(moves) != len(expectedMoves) {
		t.Fatalf("Expected %d moves, got %d", len(expectedMoves), len(moves))
	}

	for i, wantedMove := range expectedMoves {
		if moves[i] != wantedMove {
			t.Errorf("Expected move %v, got %v", wantedMove, moves[i])
		}
	}
}

func TestRookMovesFromEmptySquare(t *testing.T) {
	var board Board
	moves := rookMoves(board, Square(27))

	if len(moves) != 0 {
		t.Fatalf("Expected 0 moves, got %d", len(moves))
	}
}
