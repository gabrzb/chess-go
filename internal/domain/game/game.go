package game

import (
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/rules"
)

type Game struct {
	position    chess.Position
	history     []Snapshot
	status      chess.Status
	lastMove    *chess.Move
	outcome     Outcome
	drawOffer   chess.Color
	repetitions map[chess.Position]int
}

func NewGame() *Game {
	return &Game{
		position:    chess.NewInitialPosition(),
		status:      chess.StatusOngoing,
		repetitions: map[chess.Position]int{rules.RepetitionKey(chess.NewInitialPosition()): 1},
	}
}

func (g *Game) LegalMovesFrom(from chess.Square) ([]chess.Move, error) {
	return rules.LegalMovesFrom(g.position, from)
}
