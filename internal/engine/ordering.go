package engine

import (
	"cmp"
	"slices"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

const (
	capturePriority   = queenValue * 8
	promotionPriority = queenValue * 16
	victimWeight      = 16
)

func orderMoves(position chess.Position, moves []chess.Move) {
	slices.SortStableFunc(moves, func(a, b chess.Move) int {
		return cmp.Compare(movePriority(position, b), movePriority(position, a))
	})
}

func movePriority(position chess.Position, move chess.Move) int {
	piece, _ := position.Board.PieceAt(move.From)
	victim, _ := position.Board.PieceAt(move.To)
	if piece.Type == chess.Pawn && move.To == position.EnPassant && move.From.Column() != move.To.Column() {
		victim.Type = chess.Pawn
	}
	score := 0
	if victim.Type != chess.NoPiece {
		score += capturePriority + victimWeight*materialValue(victim.Type) - materialValue(piece.Type)
	}
	if move.Promotion != chess.NoPiece {
		score += promotionPriority + materialValue(move.Promotion)
	}
	return score
}
