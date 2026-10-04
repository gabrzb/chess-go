package chess

func rookMoves(board Board, from Square) []Move {
	rook, ok := board.PieceAt(from)

	if !ok || rook.Type != Rook {
		return nil
	}

	directions := []direction{
		{colStep: 1, rowStep: 0},  // Right
		{colStep: -1, rowStep: 0}, // Left
		{colStep: 0, rowStep: 1},  // Up
		{colStep: 0, rowStep: -1}, // Down
	}

	var moves []Move

	for _, dir := range directions {
		directionMoves := slidingMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}
