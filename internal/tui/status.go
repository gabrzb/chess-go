package tui

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func colorName(color chess.Color) string {
	switch color {
	case chess.White:
		return "White"
	case chess.Black:
		return "Black"
	default:
		return "-"
	}
}

func statusText(snapshot game.Snapshot) string {
	switch snapshot.Status {
	case chess.StatusCheckmate:
		winner := chess.White

		if snapshot.SideToMove == chess.White {
			winner = chess.Black
		}

		return "Checkmate! " + colorName(winner) + " wins."
	case chess.StatusStalemate:
		return "Stalemate! Draw."
	}

	text := "Turn: " + colorName(snapshot.SideToMove)
	if snapshot.Check {
		text += " (check)"
	}

	return text
}
