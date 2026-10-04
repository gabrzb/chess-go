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
