package tui_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/delivery/tui"
	chess "github.com/gabrzb/chess-go/internal/domain"
)

func TestClaimDrawFromPauseMenu(t *testing.T) {
	m := app(t, tui.Options{StartGame: true})
	cursor := chess.A2
	for range 2 {
		for _, move := range []chess.Move{{From: chess.G1, To: chess.F3}, {From: chess.G8, To: chess.F6}, {From: chess.F3, To: chess.G1}, {From: chess.F6, To: chess.G8}} {
			m = appMove(t, m, &cursor, move)
		}
	}
	m = keys(t, m, tea.KeyEscape, tea.KeyUp)
	contains(t, m, "▸ Reivindicar empate")
	m = keys(t, m, tea.KeyEnter)
	contains(t, m, "EMPATE")
	contains(t, m, "tripla repetição")
}
