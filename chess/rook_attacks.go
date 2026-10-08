package chess

func rookAttacksSquare(board Board, from, to Square) bool {
	rook, ok := board.PieceAt(from)

	if !ok || rook.Type != Rook {
		return false
	}

	directions := []direction{
		{colStep: 1, rowStep: 0},  // Right
		{colStep: -1, rowStep: 0}, // Left
		{colStep: 0, rowStep: 1},  // Up
		{colStep: 0, rowStep: -1}, // Down
	}

	for _, dir := range directions {
		if slidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}
