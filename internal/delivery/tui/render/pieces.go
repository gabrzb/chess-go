package render

import (
	chess "github.com/gabrzb/chess-go/internal/domain"
)

func PieceName(kind chess.PieceType) string {
	switch kind {
	case chess.Pawn:
		return "Pawn"
	case chess.Knight:
		return "Knight"
	case chess.Bishop:
		return "Bishop"
	case chess.Rook:
		return "Rook"
	case chess.Queen:
		return "Queen"
	case chess.King:
		return "King"
	default:
		return ""
	}
}

func Piece(piece chess.Piece) string {
	symbols := map[chess.PieceType]struct {
		white string
		black string
	}{
		chess.Pawn:   {white: "♙", black: "♟"},
		chess.Knight: {white: "♘", black: "♞"},
		chess.Bishop: {white: "♗", black: "♝"},
		chess.Rook:   {white: "♖", black: "♜"},
		chess.Queen:  {white: "♕", black: "♛"},
		chess.King:   {white: "♔", black: "♚"},
	}

	symbol, ok := symbols[piece.Type]
	if !ok {
		return "."
	}

	if piece.Color == chess.Black {
		return symbol.black
	}

	return symbol.white
}
