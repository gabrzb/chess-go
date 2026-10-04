package chess

func queenMoves(board Board, from Square) []Move {
	queen, ok := board.PieceAt(from)

	if !ok || queen.Type != Queen {
		return nil
	}

	directions := []direction{
		{colStep: 1, rowStep: 0},
		{colStep: -1, rowStep: 0},
		{colStep: 0, rowStep: 1},
		{colStep: 0, rowStep: -1},
		{colStep: 1, rowStep: 1},
		{colStep: -1, rowStep: 1},
		{colStep: 1, rowStep: -1},
		{colStep: -1, rowStep: -1},
	}

	var moves []Move

	for _, dir := range directions {
		directionMoves := slidingMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}
