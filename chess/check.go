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

func IsInCheck(board Board, color Color) (bool, error) {
	if color != White && color != Black {
		return false, ErrInvalidColor
	}

	kingSquare := findKing(board, color)
	if kingSquare == NoSquare {
		return false, ErrKingNotFound
	}

	enemyColor := Black
	if color == Black {
		enemyColor = White
	}

	return IsSquareAttacked(board, kingSquare, enemyColor), nil
}
