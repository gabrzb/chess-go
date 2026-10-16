package render

import (
	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func ColorName(color chess.Color) string {
	switch color {
	case chess.White:
		return "White"
	case chess.Black:
		return "Black"
	default:
		return "-"
	}
}

func StatusText(snapshot game.Snapshot) string {
	switch snapshot.Status {
	case chess.StatusCheckmate:
		winner := chess.White

		if snapshot.SideToMove == chess.White {
			winner = chess.Black
		}

		return "Checkmate! " + ColorName(winner) + " wins."
	case chess.StatusStalemate:
		return "Stalemate! Draw."
	}

	text := "Turn: " + ColorName(snapshot.SideToMove)
	if snapshot.Check {
		text += " (check)"
	}

	return text
}
