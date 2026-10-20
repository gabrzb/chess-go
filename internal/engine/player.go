package engine

import (
	"context"
	"errors"
	"math/rand"
	"time"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/rules"
)

type Difficulty string

const (
	Easy   Difficulty = "easy"
	Medium Difficulty = "medium"
	Hard   Difficulty = "hard"
)

var (
	ErrInvalidDifficulty = errors.New("AI difficulty must be easy, medium or hard")
	ErrNoLegalMoves      = errors.New("AI has no legal moves")
)

var _ chess.MoveChooser = (*Player)(nil)

type Player struct {
	difficulty Difficulty
	seed       int64
	moveTime   time.Duration
}

func NewPlayer(difficulty Difficulty, seed int64) (*Player, error) {
	if difficulty.Depth() == 0 {
		return nil, ErrInvalidDifficulty
	}
	return &Player{difficulty: difficulty, seed: seed, moveTime: 2 * time.Second}, nil
}

// Depth counts individual turns (plies), including the move being chosen.
func (d Difficulty) Depth() int {
	switch d {
	case Easy:
		return 1
	case Medium:
		return 2
	case Hard:
		return 3
	default:
		return 0
	}
}

func (p *Player) ChooseMove(ctx context.Context, position chess.Position) (chess.Move, error) {
	if err := ctx.Err(); err != nil {
		return chess.Move{}, err
	}
	depth := p.difficulty.Depth()
	if depth == 0 {
		return chess.Move{}, ErrInvalidDifficulty
	}
	if err := position.Validate(); err != nil {
		return chess.Move{}, err
	}
	moves, err := legalMoves(ctx, position)
	if err != nil {
		return chess.Move{}, err
	}
	if len(moves) == 0 {
		return chess.Move{}, ErrNoLegalMoves
	}
	orderMoves(position, moves)
	searchCtx, cancel := context.WithTimeout(ctx, p.moveTime)
	defer cancel()
	best := moves[0]
	for currentDepth := 1; currentDepth <= depth; currentDepth++ {
		candidate, err := p.chooseAtDepth(searchCtx, position, moves, currentDepth)
		if parentErr := ctx.Err(); parentErr != nil {
			return chess.Move{}, parentErr
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return best, nil
		}
		if err != nil {
			return chess.Move{}, err
		}
		best = candidate
	}
	return best, nil
}

func (p *Player) chooseAtDepth(ctx context.Context, position chess.Position, moves []chess.Move, depth int) (chess.Move, error) {
	random := rand.New(rand.NewSource(p.seed))
	bestScore, ties := -searchInfinity, 0
	var best chess.Move
	for _, move := range moves {
		if err := ctx.Err(); err != nil {
			return chess.Move{}, err
		}
		next, err := rules.PositionAfterLegalMove(position, move)
		if err != nil {
			return chess.Move{}, err
		}
		// Each root candidate gets an exact score; a cutoff bound must not enter a random tie.
		score, err := negamax(ctx, next, depth-1, 1, -searchInfinity, searchInfinity)
		if err != nil {
			return chess.Move{}, err
		}
		score = -score
		if score > bestScore {
			best, bestScore, ties = move, score, 1
		} else if score == bestScore {
			ties++
			if random.Intn(ties) == 0 {
				best = move
			}
		}
	}
	if err := ctx.Err(); err != nil {
		return chess.Move{}, err
	}
	return best, nil
}
