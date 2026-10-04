package chess

func kingMoves(board Board, from Square) []Move {
	king, ok := board.PieceAt(from)

	if !ok || king.Type != King {
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
		directionMoves := jumpMoves(board, from, dir)
		moves = append(moves, directionMoves...)
	}

	return moves
}
