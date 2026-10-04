package chess

func pawnMoves(board Board, from Square) []Move {
	pawn, ok := board.PieceAt(from)

	if !ok || pawn.Type != Pawn {
		return nil
	}

	var rowStep int

	switch pawn.Color {
	case White:
		rowStep = 1
	case Black:
		rowStep = -1
	default:
		return nil
	}

	col := int(from) % 8
	row := int(from) / 8
	nextRow := row + rowStep

	if nextRow < 0 || nextRow >= 8 {
		return nil
	}

	destination := Square(nextRow*8 + col)
	occupant, _ := board.PieceAt(destination)
	var moves []Move

	if occupant.Type == NoPiece {
		moves = append(moves, Move{From: from, To: destination})
		startRow := 1

		if pawn.Color == Black {
			startRow = 6
		}

		if row == startRow {
			doubleDestination := Square((row+2*rowStep)*8 + col)
			doubleOccupant, _ := board.PieceAt(doubleDestination)

			if doubleOccupant.Type == NoPiece {
				moves = append(moves, Move{From: from, To: doubleDestination})
			}
		}
	}

	for _, colStep := range []int{-1, 1} {
		captureCol := col + colStep

		if captureCol < 0 || captureCol >= 8 {
			continue
		}

		captureDestination := Square(nextRow*8 + captureCol)
		target, _ := board.PieceAt(captureDestination)

		if target.Type != NoPiece && target.Color != pawn.Color {
			moves = append(moves, Move{From: from, To: captureDestination})
		}
	}

	return moves
}
