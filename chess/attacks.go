package chess

func pawnAttacksSquare(from, to Square, color Color) bool {
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

func pieceAttacksSquare(board Board, from, to Square) bool {
	piece, ok := board.PieceAt(from)
	if !ok {
		return false
	}

	switch piece.Type {
	case Pawn:
		return pawnAttacksSquare(from, to, piece.Color)
	case Knight:
		return knightAttacksSquare(from, to)
	case Bishop:
		return bishopAttacksSquare(board, from, to)
	case Rook:
		return rookAttacksSquare(board, from, to)
	case Queen:
		return queenAttacksSquare(board, from, to)
	case King:
		return kingAttacksSquare(from, to)
	default:
		return false
	}
}

func IsSquareAttacked(board Board, square Square, byColor Color) bool {
	if !square.Valid() || (byColor != White && byColor != Black) {
		return false
	}

	for from := Square(0); from < NoSquare; from++ {
		piece, _ := board.PieceAt(from)

		if piece.Type == NoPiece || piece.Color != byColor {
			continue
		}

		if pieceAttacksSquare(board, from, square) {
			return true
		}
	}

	return false
}
