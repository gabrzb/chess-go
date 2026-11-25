package components

import (
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

var drawReasons = map[game.EndReason]string{
	game.Stalemate: "Stalemate", game.Agreement: "Draw agreed",
	game.InsufficientMaterial: "Insufficient material", game.Repetition: "Threefold repetition",
	game.FiftyMoves: "Fifty-move rule", game.FivefoldRepetition: "Fivefold repetition",
	game.SeventyFiveMoves: "Seventy-five-move rule",
}

func (a *App) playerName(c chess.Color) string {
	switch {
	case !a.vsAI && c == chess.White:
		return "Player 1"
	case !a.vsAI:
		return "Player 2"
	case c == a.opponent.Color():
		return "AI · " + LevelLabels[a.Level]
	}
	return "You"
}

func (a *App) result(s game.Snapshot) (string, string) {
	switch s.Outcome.Reason {
	case game.Checkmate:
		return a.winnerTitle(s.Outcome.Winner), "Checkmate."
	case game.Resignation:
		return a.winnerTitle(s.Outcome.Winner), colorName(other(s.Outcome.Winner)) + " resigned."
	}
	return "Draw", drawReasons[s.Outcome.Reason] + "."
}

func (a *App) winnerTitle(winner chess.Color) string {
	switch {
	case a.vsAI && winner == a.opponent.Color():
		return "AI wins"
	case a.vsAI:
		return "You win!"
	}
	return colorName(winner) + " wins"
}

func (a *App) status(s game.Snapshot) string {
	switch {
	case s.Status != chess.StatusOngoing:
		title, why := a.result(s)
		return title + " — " + why
	case a.opponent.Thinking():
		return "AI is thinking…"
	case s.Check:
		return colorName(s.SideToMove) + " to move — check!"
	}
	return colorName(s.SideToMove) + " to move"
}
