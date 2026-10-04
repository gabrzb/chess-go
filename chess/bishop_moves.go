package chess

func bishopMoves(board Board, from Square) []Move {
	bishop, ok := board.PieceAt(from)

	if !ok || bishop.Type != Bishop {
		return nil
	}

	directions := []direction{
		{colStep: 1, rowStep: 1},   // diagonal up-right
		{colStep: -1, rowStep: 1},  // diagonal up-left
		{colStep: 1, rowStep: -1},  // diagonal down-right
		{colStep: -1, rowStep: -1}, // diagonal down-left
	}

	var moves []Move

	for _, dir := range directions {
		directionMoves := slidingMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}
