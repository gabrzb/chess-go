package render

import (
	chess "github.com/gabrzb/chess-go/internal/domain"
)

func PieceName(kind chess.PieceType) string {
	switch kind {
	case chess.Pawn:
		return "Peão"
	case chess.Knight:
		return "Cavalo"
	case chess.Bishop:
		return "Bispo"
	case chess.Rook:
		return "Torre"
	case chess.Queen:
		return "Dama"
	case chess.King:
		return "Rei"
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
