package chess

type direction struct {
	colStep int
	rowStep int
}

func slidingMoves(board Board, from Square, dir direction) []Move {
	if !from.Valid() {
		return nil
	}

	col := int(from) % 8
	row := int(from) / 8
	var moves []Move
	piece, _ := board.PieceAt(from)

	if piece.Type == NoPiece {
		return nil
	}

	nextCol := col + dir.colStep
	nextRow := row + dir.rowStep

	for nextCol >= 0 && nextCol < 8 && nextRow >= 0 && nextRow < 8 {
		destination := Square(nextRow*8 + nextCol)
		occupant, _ := board.PieceAt(destination)

		if occupant.Type != NoPiece && occupant.Color == piece.Color {
			break
		}

		moves = append(moves, Move{From: from, To: destination})

		if occupant.Type != NoPiece {
			break
		}

		nextCol += dir.colStep
		nextRow += dir.rowStep
	}

	return moves
}

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
