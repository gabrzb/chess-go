package ai

import "github.com/gabrzb/chess-go/internal/chess"

const (
	pawnValue   = 100
	knightValue = 320
	bishopValue = 330
	rookValue   = 500
	queenValue  = 900

	pawnAdvanceBonus  = 8
	pawnFileBonus     = 2
	knightCenterBonus = 10
	bishopCenterBonus = 6
	queenCenterBonus  = 2
)

// Evaluate returns a score in centipawns: positive favors White, negative favors Black.
// Checkmate and stalemate belong to the search, not the static board evaluation.
func Evaluate(board chess.Board) int {
	score := 0
	for square := chess.A1; square <= chess.H8; square++ {
		piece, _ := board.PieceAt(square)
		if piece.Type == chess.NoPiece {
			continue
		}
		value := materialValue(piece.Type) + positionalBonus(piece, square)
		if piece.Color == chess.White {
			score += value
		} else if piece.Color == chess.Black {
			score -= value
		}
	}
	return score
}

func materialValue(piece chess.PieceType) int {
	switch piece {
	case chess.Pawn:
		return pawnValue
	case chess.Knight:
		return knightValue
	case chess.Bishop:
		return bishopValue
	case chess.Rook:
		return rookValue
	case chess.Queen:
		return queenValue
	default:
		return 0
	}
}

func positionalBonus(piece chess.Piece, square chess.Square) int {
	edge := chess.BoardSize - 1
	col, row := square.Column(), square.Row()
	fileCenter := min(col, edge-col)
	center := fileCenter + min(row, edge-row)
	switch piece.Type {
	case chess.Pawn:
		if piece.Color == chess.Black {
			row = edge - row
		}
		advance := max(0, row-chess.A2.Row())
		return advance*pawnAdvanceBonus + fileCenter*pawnFileBonus
	case chess.Knight:
		return center * knightCenterBonus
	case chess.Bishop:
		return center * bishopCenterBonus
	case chess.Queen:
		return center * queenCenterBonus
	default:
		return 0
	}
}
