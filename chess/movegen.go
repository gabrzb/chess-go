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

func canOccupySquare(board Board, to Square, color Color) bool {
	occupant, ok := board.PieceAt(to)

	if !ok {
		return false
	}

	return occupant.Type == NoPiece || occupant.Color != color
}

func jumpMoves(board Board, from Square, dir direction) []Move {
	piece, ok := board.PieceAt(from)

	if !ok || piece.Type == NoPiece {
		return nil
	}

	col := int(from) % 8
	row := int(from) / 8

	nextCol := col + dir.colStep
	nextRow := row + dir.rowStep

	if nextCol < 0 || nextCol >= 8 || nextRow < 0 || nextRow >= 8 {
		return nil
	}

	destination := Square(nextRow*8 + nextCol)

	if !canOccupySquare(board, destination, piece.Color) {
		return nil
	}

	return []Move{{From: from, To: destination}}
}
