package tui

import (
	"github.com/gabrzb/chess-go/internal/chess"
)

func pieceSymbol(piece chess.Piece) string {
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
