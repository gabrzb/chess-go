package game

import "github.com/gabrzb/chess-go/internal/chess"

type MoveResult struct {
	Move  chess.Move
	State Snapshot
}
