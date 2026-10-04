package chess

func rookMoves(board Board, from Square) []Move {
	if !from.Valid() {
		return nil
	}

	col := int(from) % 8
	row := int(from) / 8
	var moves []Move
	rook, _ := board.PieceAt(from)

	if rook.Type != Rook {
		return nil
	}

	for nextCol := col + 1; nextCol < 8; nextCol++ {
		destination := Square(row*8 + nextCol)
		occupant, _ := board.PieceAt(destination)

		if occupant.Type != NoPiece && occupant.Color == rook.Color {
			break
		}

		moves = append(moves, Move{From: from, To: destination})

		if occupant.Type != NoPiece {
			break
		}
	}

	return moves
}
