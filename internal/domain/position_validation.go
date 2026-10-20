package domain

import "fmt"

// Validate checks board and metadata consistency, not historical reachability.
// Positions remain editable values so tests and search can construct them freely.
func (p Position) Validate() error {
	if p.SideToMove != White && p.SideToMove != Black {
		return ErrInvalidColor
	}
	if p.FullmoveNumber == 0 || p.CastlingRights & ^(WhiteKingSide|WhiteQueenSide|BlackKingSide|BlackQueenSide) != 0 {
		return fmt.Errorf("%w: invalid clocks or castling flags", ErrInvalidPosition)
	}
	var kings [3]int
	for square := Square(0); square < NoSquare; square++ {
		piece, _ := p.Board.PieceAt(square)
		if piece == (Piece{}) {
			continue
		}
		if piece.Type < Pawn || piece.Type > King || (piece.Color != White && piece.Color != Black) {
			return fmt.Errorf("%w: invalid piece at %s", ErrInvalidPosition, square)
		}
		if piece.Type == Pawn && (square.Row() == 0 || square.Row() == BoardSize-1) {
			return fmt.Errorf("%w: unpromoted pawn at %s", ErrInvalidPosition, square)
		}
		if piece.Type == King {
			kings[piece.Color]++
		}
	}
	for _, color := range []Color{White, Black} {
		if kings[color] == 0 {
			return ErrKingNotFound
		}
		if kings[color] != 1 {
			return fmt.Errorf("%w: multiple kings of one color", ErrInvalidPosition)
		}
	}
	for _, side := range []struct {
		right      CastlingRights
		color      Color
		king, rook Square
	}{
		{WhiteKingSide, White, E1, H1}, {WhiteQueenSide, White, E1, A1},
		{BlackKingSide, Black, E8, H8}, {BlackQueenSide, Black, E8, A8},
	} {
		king, _ := p.Board.PieceAt(side.king)
		rook, _ := p.Board.PieceAt(side.rook)
		if p.CastlingRights.Has(side.right) && (king != (Piece{King, side.color}) || rook != (Piece{Rook, side.color})) {
			return fmt.Errorf("%w: castling right without king and rook", ErrInvalidPosition)
		}
	}
	if p.EnPassant != NoSquare {
		row, step, color := 5, -1, Black
		if p.SideToMove == Black {
			row, step, color = 2, 1, White
		}
		pawnSquare, _ := p.EnPassant.Offset(0, step)
		origin, _ := p.EnPassant.Offset(0, -step)
		target, _ := p.Board.PieceAt(p.EnPassant)
		pawn, _ := p.Board.PieceAt(pawnSquare)
		previous, _ := p.Board.PieceAt(origin)
		if !p.EnPassant.Valid() || p.EnPassant.Row() != row || target != (Piece{}) || previous != (Piece{}) || pawn != (Piece{Pawn, color}) || p.HalfmoveClock != 0 {
			return fmt.Errorf("%w: inconsistent en passant target", ErrInvalidPosition)
		}
	}
	return nil
}
