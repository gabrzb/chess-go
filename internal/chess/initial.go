package chess

func NewInitialPosition() Position {
	position := Position{
		SideToMove:     White,
		CastlingRights: WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide,
		EnPassant:      NoSquare,
		FullmoveNumber: 1,
	}

	backRank := [BoardSize]PieceType{
		Rook, Knight, Bishop, Queen,
		King, Bishop, Knight, Rook,
	}

	sides := []struct {
		color                Color
		backStart, pawnStart Square
	}{
		{White, A1, A2},
		{Black, A8, A7},
	}

	for _, side := range sides {
		for col, pieceType := range backRank {
			backSquare, _ := side.backStart.Offset(col, 0)
			pawnSquare, _ := side.pawnStart.Offset(col, 0)

			position.Board.squares[backSquare] = Piece{
				Type:  pieceType,
				Color: side.color,
			}
			position.Board.squares[pawnSquare] = Piece{
				Type:  Pawn,
				Color: side.color,
			}
		}
	}

	return position
}
