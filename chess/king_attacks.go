package chess

func kingAttacksSquare(from, to Square) bool {
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

	return colDiff <= 1 && rowDiff <= 1 && (colDiff != 0 || rowDiff != 0)
}
