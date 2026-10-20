package game

import chess "github.com/gabrzb/chess-go/internal/domain"

type MoveResult struct {
	Move  chess.Move
	State Snapshot
}
