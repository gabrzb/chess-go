package game

import (
	"errors"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/rules"
)

var ErrDrawNotClaimable = errors.New("no draw can be claimed in this position")

// ClaimableDraw reports a claim available to the side to move in the current
// position. A claim does not require the opponent's agreement.
func (g *Game) ClaimableDraw() EndReason {
	if g.status != chess.StatusOngoing || g.drawOffer != chess.NoColor {
		return NoReason
	}
	if g.repetitions[rules.RepetitionKey(g.position)] >= 3 {
		return Repetition
	}
	if g.position.HalfmoveClock >= 100 {
		return FiftyMoves
	}
	return NoReason
}

func (g *Game) ClaimDraw() error {
	if g.status != chess.StatusOngoing {
		return ErrGameOver
	}
	reason := g.ClaimableDraw()
	switch reason {
	case Repetition:
		g.status = chess.StatusRepetition
	case FiftyMoves:
		g.status = chess.StatusFiftyMoves
	default:
		return ErrDrawNotClaimable
	}
	g.outcome = Outcome{Reason: reason}
	return nil
}
