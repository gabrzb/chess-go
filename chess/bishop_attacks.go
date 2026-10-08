package chess

func bishopAttacksSquare(board Board, from, to Square) bool {
	bishop, ok := board.PieceAt(from)

	if !ok || bishop.Type != Bishop {
		return false
	}

	directions := []direction{
		{colStep: 1, rowStep: 1},   // Right
		{colStep: -1, rowStep: 1},  // Left
		{colStep: 1, rowStep: -1},  // Up
		{colStep: -1, rowStep: -1}, // Down
	}

	for _, dir := range directions {
		if slidingAttacksSquare(board, from, to, dir) {
			return true
		}
	}

	return false
}
