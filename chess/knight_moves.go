package chess

func knightMoves(board Board, from Square) []Move {
	knight, ok := board.PieceAt(from)

	if !ok || knight.Type != Knight {
		return nil
	}

	directions := []direction{
		{colStep: 2, rowStep: 1},
		{colStep: 2, rowStep: -1},
		{colStep: -2, rowStep: 1},
		{colStep: -2, rowStep: -1},
		{colStep: 1, rowStep: 2},
		{colStep: -1, rowStep: 2},
		{colStep: 1, rowStep: -2},
		{colStep: -1, rowStep: -2},
	}

	var moves []Move

	for _, dir := range directions {
		directionMoves := jumpMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}
