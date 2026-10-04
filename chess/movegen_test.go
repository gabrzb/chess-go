package chess

import "testing"

func checkSlidingMoves(t *testing.T, dir direction, destinations []Square) {
	t.Helper()

	var board Board
	rook := Piece{Type: Rook, Color: White}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := slidingMoves(
		board,
		Square(27),
		dir,
	)

	if len(moves) != len(destinations) {
		t.Fatalf("Expected %d moves, got %d", len(destinations), len(moves))
	}

	for i, destination := range destinations {
		expectedMove := Move{From: Square(27), To: destination}

		if moves[i] != expectedMove {
			t.Errorf("Expected move %v, got %v", expectedMove, moves[i])
		}
	}
}

func TestSlidingMovesToRight(t *testing.T) {
	checkSlidingMoves(t,
		direction{colStep: 1, rowStep: 0},
		[]Square{28, 29, 30, 31},
	)
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

	moves := slidingMoves(
		board,
		Square(27),
		direction{colStep: 1, rowStep: 0},
	)

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

	moves := slidingMoves(
		board,
		Square(27),
		direction{colStep: 1, rowStep: 0},
	)

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

func TestSlidingMovesToLeft(t *testing.T) {
	checkSlidingMoves(t,
		direction{colStep: -1, rowStep: 0},
		[]Square{26, 25, 24},
	)
}

func TestSlidingMovesUp(t *testing.T) {
	checkSlidingMoves(t,
		direction{colStep: 0, rowStep: 1},
		[]Square{35, 43, 51, 59},
	)
}

func TestSlidingMovesDown(t *testing.T) {
	checkSlidingMoves(t,
		direction{colStep: 0, rowStep: -1},
		[]Square{19, 11, 3},
	)
}

func TestRookMoves(t *testing.T) {
	var board Board
	rook := Piece{Type: Rook, Color: White}

	if err := board.SetPiece(Square(27), rook); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := rookMoves(board, Square(27))

	destinationSquares := []Square{
		28, 29, 30, 31, // Right
		26, 25, 24, // Left
		35, 43, 51, 59, // Up
		19, 11, 3, // Down
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
