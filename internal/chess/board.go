package chess

type Board struct {
	squares [64]Piece
}

func (b Board) PieceAt(square Square) (Piece, bool) {
	if square.Valid() {
		return b.squares[square], true
	}

	return Piece{}, false
}

func (b *Board) SetPiece(square Square, piece Piece) error {
	if !square.Valid() {
		return ErrInvalidSquare
	}

	b.squares[square] = piece
	return nil
}

func (b *Board) RemovePiece(square Square) (Piece, error) {
	if !square.Valid() {
		return Piece{}, ErrInvalidSquare
	}

	removed := b.squares[square]
	b.squares[square] = Piece{}

	return removed, nil
}
