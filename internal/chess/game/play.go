package game

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

func (g *Game) Play(move chess.Move) (MoveResult, error) {
	if g.status != chess.StatusOngoing {
		return MoveResult{}, ErrGameOver
	}
	if g.drawOffer != chess.NoColor {
		return MoveResult{}, ErrDrawOffer
	}

	nextPosition, err := rules.PositionAfterMove(g.position, move)
	if err != nil {
		return MoveResult{}, err
	}

	nextStatus, err := rules.EvaluateStatus(nextPosition)
	if err != nil {
		return MoveResult{}, err
	}
	key := rules.RepetitionKey(nextPosition)
	if nextStatus == chess.StatusOngoing && g.repetitions[key]+1 >= 5 {
		nextStatus = chess.StatusFivefoldRepetition
	}

	next := Game{
		position: nextPosition,
		status:   nextStatus,
		lastMove: &move,
		outcome:  outcomeFor(nextStatus, nextPosition.SideToMove),
	}

	snapshot, err := next.Snapshot()
	if err != nil {
		return MoveResult{}, err
	}

	previous, err := g.Snapshot()
	if err != nil {
		return MoveResult{}, err
	}

	g.history = append(g.history, previous)
	if g.repetitions == nil {
		g.repetitions = make(map[chess.Position]int)
	}
	g.repetitions[key]++
	g.position = nextPosition
	g.status = nextStatus
	g.lastMove = &move
	g.outcome = next.outcome

	return MoveResult{
		Move:  move,
		State: snapshot,
	}, nil
}
