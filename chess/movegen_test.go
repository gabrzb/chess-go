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

func TestCanOccupySquareEmpty(t *testing.T) {
	var board Board

	if !canOccupySquare(board, Square(27), White) {
		t.Errorf("Expected canOccupySquare to return true for empty square")
	}
}

func TestCanOccupySquareInvalid(t *testing.T) {
	var board Board

	if canOccupySquare(board, NoSquare, White) {
		t.Errorf("Expected canOccupySquare to return false for NoSquare")
	}
}

func TestCanOccupySquareOwnPiece(t *testing.T) {
	var board Board
	pawn := Piece{Type: Pawn, Color: White}

	if err := board.SetPiece(Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if canOccupySquare(board, Square(27), White) {
		t.Errorf("Expected canOccupySquare to return false for own piece")
	}
}

func TestCanOccupySquareEnemyPiece(t *testing.T) {
	var board Board
	pawn := Piece{Type: Pawn, Color: Black}

	if err := board.SetPiece(Square(27), pawn); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	if !canOccupySquare(board, Square(27), White) {
		t.Errorf("Expected canOccupySquare to return true for enemy piece")
	}
}

func TestJumpMoves(t *testing.T) {
	var board Board
	knight := Piece{Type: Knight, Color: White}

	if err := board.SetPiece(Square(27), knight); err != nil {
		t.Fatalf("Failed to set piece: %v", err)
	}

	moves := jumpMoves(
		board,
		Square(27),
		direction{colStep: 2, rowStep: 1},
	)

	if len(moves) != 1 {
		t.Fatalf("Expected 1 move, got %d", len(moves))
	}

	expectedMove := Move{From: Square(27), To: Square(37)}

	if moves[0] != expectedMove {
		t.Errorf("Expected move %v, got %v", expectedMove, moves[0])
	}
}
