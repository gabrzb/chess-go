package chess

func slidingAttacksSquare(
	board Board,
	from, to Square,
	dir direction,
) bool {
	if !from.Valid() || !to.Valid() || from == to {
		return false
	}

	nextCol := int(from)%8 + dir.colStep
	nextRow := int(from)/8 + dir.rowStep

	for nextCol >= 0 && nextCol < 8 && nextRow >= 0 && nextRow < 8 {
		square := Square(nextRow*8 + nextCol)

		if square == to {
			return true
		}

		occupant, _ := board.PieceAt(square)
		if occupant.Type != NoPiece {
			return false
		}

		nextCol += dir.colStep
		nextRow += dir.rowStep
	}

	return false
}
