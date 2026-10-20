package engine

import (
	"context"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/rules"
)

const (
	mateScore          = 1_000_000
	searchInfinity     = mateScore + 1
	maxQuiescencePlies = 8
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
	if depth == 0 {
		return quiescence(ctx, position, moves, ply, maxQuiescencePlies, alpha, beta)
	}
	if score, terminal, err := terminalScore(position, moves, ply); terminal || err != nil {
		return score, err
	}
	orderMoves(position, moves)
	best := -searchInfinity
	for _, move := range moves {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		next, err := rules.PositionAfterLegalMove(position, move)
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

func terminalScore(position chess.Position, moves []chess.Move, ply int) (int, bool, error) {
	if len(moves) == 0 {
		inCheck, err := rules.IsInCheck(position.Board, position.SideToMove)
		if err != nil {
			return 0, false, err
		}
		if inCheck {
			return -mateScore + ply, true, nil
		}
		return 0, true, nil
	}
	if rules.InsufficientMaterial(position.Board) || position.HalfmoveClock >= 150 {
		return 0, true, nil
	}
	return 0, false, nil
}

// Quiescence extends captures and promotions, or every evasion when in check.
// shortcut: cap extensions at eight plies, raise only with measured tactical need.
func quiescence(ctx context.Context, position chess.Position, moves []chess.Move, ply, remaining, alpha, beta int) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if score, terminal, err := terminalScore(position, moves, ply); terminal || err != nil {
		return score, err
	}
	static := Evaluate(position.Board)
	if position.SideToMove == chess.Black {
		static = -static
	}
	if remaining == 0 {
		return static, nil
	}
	inCheck, err := rules.IsInCheck(position.Board, position.SideToMove)
	if err != nil {
		return 0, err
	}
	best := -searchInfinity
	if !inCheck {
		best = static
		if best >= beta {
			return best, nil
		}
		alpha = max(alpha, best)
	}
	orderMoves(position, moves)
	for _, move := range moves {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		if !inCheck && !tacticalMove(position, move) {
			continue
		}
		next, err := rules.PositionAfterLegalMove(position, move)
		if err != nil {
			return 0, err
		}
		nextMoves, err := legalMoves(ctx, next)
		if err != nil {
			return 0, err
		}
		score, err := quiescence(ctx, next, nextMoves, ply+1, remaining-1, -beta, -alpha)
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

func tacticalMove(position chess.Position, move chess.Move) bool {
	victim, _ := position.Board.PieceAt(move.To)
	piece, _ := position.Board.PieceAt(move.From)
	return victim.Type != chess.NoPiece || move.Promotion != chess.NoPiece ||
		(piece.Type == chess.Pawn && move.From.Column() != move.To.Column())
}
