package components

import (
	"errors"
	"time"

	"github.com/gabrzb/chess-go/internal/application"
)

const aiMinDelay = 500 * time.Millisecond

func (a *App) StartAI() {
	update := a.update
	if update == nil || a.Screen != GameScreen || a.paused {
		return
	}
	search, err := a.opponent.Start()
	if err != nil {
		a.message = err.Error()
		return
	}
	if search == nil {
		return
	}
	go func() {
		wait := time.After(aiMinDelay)
		result := search()
		<-wait
		update(func() { a.acceptAI(result) })
	}()
}

func (a *App) acceptAI(result application.SearchResult) {
	before, _ := a.game.Snapshot()
	accepted, err := a.opponent.Accept(result)
	if !accepted {
		return
	}
	switch {
	case errors.Is(err, application.ErrPositionChanged):
		a.message = "AI result discarded: the position changed."
	case err != nil:
		a.message = "AI failed: " + err.Error()
	default:
		move := result.Move
		a.sliding = &move
		a.record(before, move)
	}
}
