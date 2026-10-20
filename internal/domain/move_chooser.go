package domain

import "context"

// MoveChooser selects a move from an independent position value.
// Callers must validate the returned move before changing the live game.
type MoveChooser interface {
	ChooseMove(context.Context, Position) (Move, error)
}
