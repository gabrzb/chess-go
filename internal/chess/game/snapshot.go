package game

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

type Snapshot struct {
	Board      chess.Board
	SideToMove chess.Color
	Status     chess.Status
	Check      bool
	LastMove   *chess.Move
}

func (g *Game) Snapshot() (Snapshot, error) {
	inCheck, err := rules.IsInCheck(g.position.Board, g.position.SideToMove)
	if err != nil {
		return Snapshot{}, err
	}

	snapshot := Snapshot{
		Board:      g.position.Board,
		SideToMove: g.position.SideToMove,
		Status:     g.status,
		Check:      inCheck,
	}

	if g.lastMove != nil {
		move := *g.lastMove
		snapshot.LastMove = &move
	}

	return snapshot, nil
}
