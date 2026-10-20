package game

import chess "github.com/gabrzb/chess-go/internal/domain"

// Position returns an independent value, including castling and en passant state,
// so a search can simulate future turns without changing the live game.
func (g *Game) Position() chess.Position {
	return g.position
}
