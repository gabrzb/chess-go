package chess

type Color uint8

const (
	NoColor Color = iota
	White
	Black
)

type PieceType uint8

const (
	NoPiece PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

type Piece struct {
	Type  PieceType
	Color Color
}

type Status uint8

const (
	StatusOngoing Status = iota
	StatusCheckmate
	StatusStalemate
)

// Direction describes a displacement in board columns and rows.
type Direction struct {
	ColStep int
	RowStep int
}
