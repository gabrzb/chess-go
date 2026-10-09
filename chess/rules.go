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

func LegalMoves(board Board, from Square) ([]Move, error) {
	piece, ok := board.PieceAt(from)
	if !ok {
		return nil, ErrInvalidSquare
	}

	if piece.Type == NoPiece {
		return nil, ErrEmptySource
	}

	if _, err := IsInCheck(board, piece.Color); err != nil {
		return nil, err
	}

	var moves []Move

	for _, move := range pseudoLegalMoves(board, from) {
		target, _ := board.PieceAt(move.To)
		if target.Type == King {
			continue
		}

		nextMove, err := boardAfterMove(board, move)
		if err != nil {
			return nil, err
		}

		inCheck, err := IsInCheck(nextMove, piece.Color)
		if err != nil {
			return nil, err
		}

		if !inCheck {
			moves = append(moves, move)
		}
	}

	return moves, nil
}
