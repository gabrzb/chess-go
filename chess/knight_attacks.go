package chess

func knightAttacksSquare(from, to Square) bool {
	if !from.Valid() || !to.Valid() {
		return false
	}

	colDiff := int(to)%8 - int(from)%8
	rowDiff := int(to)/8 - int(from)/8

	if colDiff < 0 {
		colDiff = -colDiff
	}

	if rowDiff < 0 {
		rowDiff = -rowDiff
	}

	return (colDiff == 1 && rowDiff == 2) || (colDiff == 2 && rowDiff == 1)
}
