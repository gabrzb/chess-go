package ai

import (
	"context"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

const (
	mateScore      = 1_000_000
	searchInfinity = mateScore + 1
)

func legalMoves(ctx context.Context, position chess.Position) ([]chess.Move, error) {
	var moves []chess.Move
	for square := chess.A1; square <= chess.H8; square++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		piece, _ := position.Board.PieceAt(square)
		if piece.Type == chess.NoPiece || piece.Color != position.SideToMove {
			continue
		}
		fromMoves, err := rules.LegalMovesFrom(position, square)
		if err != nil {
			return nil, err
		}
		moves = append(moves, fromMoves...)
	}
	return moves, nil
}

// Negamax expresses minimax from the side-to-move perspective: the opponent's
// best score is negated. Alpha-beta bounds discard branches that cannot improve it.
func negamax(ctx context.Context, position chess.Position, depth, ply, alpha, beta int) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	moves, err := legalMoves(ctx, position)
	if err != nil {
		return 0, err
	}
	if len(moves) == 0 {
		inCheck, err := rules.IsInCheck(position.Board, position.SideToMove)
		if err != nil {
			return 0, err
		}
		if inCheck {
			return -mateScore + ply, nil
		}
		return 0, nil
	}
	if depth == 0 {
		score := Evaluate(position.Board)
		if position.SideToMove == chess.Black {
			score = -score
		}
		return score, nil
	}
	orderMoves(position, moves)
	best := -searchInfinity
	for _, move := range moves {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		next, err := rules.PositionAfterMove(position, move)
		if err != nil {
			return 0, err
		}
		score, err := negamax(ctx, next, depth-1, ply+1, -beta, -alpha)
		if err != nil {
			return 0, err
		}
		score = -score
		best = max(best, score)
		alpha = max(alpha, score)
		if alpha >= beta {
			break
		}
	}
	return best, nil
}
