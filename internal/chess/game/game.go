package game

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

type Game struct {
	position chess.Position
	history  []Snapshot
	status   chess.Status
	lastMove *chess.Move
}

func NewGame() *Game {
	return &Game{
		position: chess.NewInitialPosition(),
		status:   chess.StatusOngoing,
	}
}

func (g *Game) LegalMovesFrom(from chess.Square) ([]chess.Move, error) {
	return rules.LegalMovesFrom(g.position, from)
}
