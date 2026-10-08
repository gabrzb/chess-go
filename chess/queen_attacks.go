package chess

func queenAttacksSquare(board Board, from, to Square) bool {
	queen, ok := board.PieceAt(from)

	if !ok || queen.Type != Queen {
		return false
	}

	directions := []direction{
		{colStep: 1, rowStep: 0},   // Right
		{colStep: -1, rowStep: 0},  // Left
		{colStep: 0, rowStep: 1},   // Up
		{colStep: 0, rowStep: -1},  // Down
		{colStep: 1, rowStep: 1},   // Up-Right
		{colStep: -1, rowStep: 1},  // Up-Left
		{colStep: 1, rowStep: -1},  // Down-Right
		{colStep: -1, rowStep: -1}, // Down-Left
	}

	for _, dir := range directions {
		if slidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}
