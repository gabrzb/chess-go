package chess

func findKing(board Board, color Color) Square {
	if color != White && color != Black {
		return NoSquare
	}

	for square := Square(0); square < NoSquare; square++ {
		piece, _ := board.PieceAt(square)

		if piece.Type == King && piece.Color == color {
			return square
		}
	}

	return NoSquare
}
