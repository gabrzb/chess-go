package application

import (
	"context"
	"errors"

	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

var (
	ErrInvalidAIConfig = errors.New("AI requires a context, game, move chooser and valid color")
	ErrPositionChanged = errors.New("AI result discarded: the position changed")
)

// Opponent coordinates searches on the caller's event loop. Only the returned
// search function runs in the background; it captures independent input values.
type Opponent struct {
	ctx              context.Context
	game             *game.Game
	chooser          chess.MoveChooser
	color            chess.Color
	session, request uint64
	cancel           context.CancelFunc
	thinking         bool
}

type SearchResult struct {
	Session  uint64
	Move     chess.Move
	Err      error
	request  uint64
	position chess.Position
}

func NewOpponent(ctx context.Context, g *game.Game, chooser chess.MoveChooser, color chess.Color, session uint64) (Opponent, error) {
	if ctx == nil || g == nil || chooser == nil || (color != chess.White && color != chess.Black) {
		return Opponent{}, ErrInvalidAIConfig
	}
	return Opponent{ctx: ctx, game: g, chooser: chooser, color: color, session: session}, nil
}

func (o Opponent) Enabled() bool      { return o.chooser != nil }
func (o Opponent) Color() chess.Color { return o.color }
func (o Opponent) Thinking() bool     { return o.thinking }
func (o Opponent) IsTurn() bool       { return o.Enabled() && o.game.Position().SideToMove == o.color }

func (o *Opponent) Start() (func() SearchResult, error) {
	if !o.IsTurn() || o.thinking {
		return nil, nil
	}
	snapshot, err := o.game.Snapshot()
	if err != nil {
		return nil, err
	}
	if snapshot.Status != chess.StatusOngoing {
		return nil, nil
	}
	ctx, cancel := context.WithCancel(o.ctx)
	o.cancel, o.thinking = cancel, true
	o.request++
	request, session, position, chooser := o.request, o.session, o.game.Position(), o.chooser
	return func() SearchResult {
		move, err := chooser.ChooseMove(ctx, position)
		return SearchResult{Session: session, Move: move, Err: err, request: request, position: position}
	}, nil
}

func (o *Opponent) Stop() {
	if o.cancel != nil {
		o.cancel()
		o.cancel = nil
	}
	o.thinking = false
}

// Accept ignores canceled, duplicate, and previous-session results. A current
// result is applied through Game.Play so even custom choosers must return legal moves.
func (o *Opponent) Accept(result SearchResult) (bool, error) {
	if !o.thinking || result.request != o.request || result.Session != o.session {
		return false, nil
	}
	o.Stop()
	if result.position != o.game.Position() {
		return true, ErrPositionChanged
	}
	if result.Err != nil {
		return true, result.Err
	}
	_, err := o.game.Play(result.Move)
	return true, err
}
