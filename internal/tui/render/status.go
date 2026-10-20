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
	case chess.StatusInsufficientMaterial:
		return "Draw: insufficient material."
	case chess.StatusRepetition, chess.StatusFivefoldRepetition:
		return "Draw: repetition."
	case chess.StatusFiftyMoves, chess.StatusSeventyFiveMoves:
		return "Draw: move-count rule."
	case chess.StatusDrawAgreed:
		return "Draw by agreement."
	case chess.StatusResigned:
		return "Resignation! " + ColorName(snapshot.Outcome.Winner) + " wins."
	}

	text := "Turn: " + ColorName(snapshot.SideToMove)
	if snapshot.Check {
		text += " (check)"
	}

	return text
}
