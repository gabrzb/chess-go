package chess

func boardAfterMove(board Board, move Move) (Board, error) {
	if !move.From.Valid() || !move.To.Valid() {
		return Board{}, ErrInvalidSquare
	}

	piece, _ := board.PieceAt(move.From)
	if piece.Type == NoPiece {
		return Board{}, ErrEmptySource
	}

	board.squares[move.From] = Piece{}
	board.squares[move.To] = piece

	return board, nil
}
