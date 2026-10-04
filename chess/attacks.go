package chess

func pawnAttackSquare(from, to Square, color Color) bool {
	if !from.Valid() || !to.Valid() {
		return false
	}

	colDiff := int(to)%8 - int(from)%8
	rowDiff := int(to)/8 - int(from)/8

	switch color {
	case White:
		return rowDiff == 1 && (colDiff == 1 || colDiff == -1)
	case Black:
		return rowDiff == -1 && (colDiff == 1 || colDiff == -1)
	default:
		return false
	}
}
