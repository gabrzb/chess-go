package game

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func (g *Game) Play(move chess.Move) (MoveResult, error) {
	if g.status != chess.StatusOngoing {
		return MoveResult{}, ErrGameOver
	}

	nextPosition, err := rules.PositionAfterMove(g.position, move)
	if err != nil {
		return MoveResult{}, err
	}

	nextStatus, err := rules.EvaluateStatus(nextPosition)
	if err != nil {
		return MoveResult{}, err
	}

	next := Game{
		position: nextPosition,
		status:   nextStatus,
		lastMove: &move,
	}

	snapshot, err := next.Snapshot()
	if err != nil {
		return MoveResult{}, err
	}

	g.history = append(g.history, g.position)
	g.position = nextPosition
	g.status = nextStatus
	g.lastMove = &move

	return MoveResult{
		Move:  move,
		State: snapshot,
	}, nil
}
