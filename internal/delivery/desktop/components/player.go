package components

import (
	"github.com/egoist/mygo/ui"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

const PlayerBarHeight = 44

func (a *App) Player(c *ui.Context, color chess.Color, width float32) {
	t := c.Theme()
	live, _ := a.game.Snapshot()
	turn := live.SideToMove == color && live.Status == chess.StatusOngoing
	row := ui.Row(c.Key(color)).Width(width).Height(PlayerBarHeight).Gap(10).Padding(0, 12).Radius(8).Transition(ui.ElementTransition{})
	if turn {
		row.Background(t.Surface).Border(1, t.Border)
	}
	row.Children(func() {
		dot, ring := whiteInk, ui.RGBA(0, 0, 0, 0.25)
		if color == chess.Black {
			dot, ring = blackInk, ui.RGBA(255, 255, 255, 0.25)
		}
		ui.Box(c).Size(18, 18).Radius(9).Background(dot).Border(1, ring)
		ui.Text(c, a.playerName(color)).Bold()
		ui.Text(c, colorName(color)).FontSize(12).TextColor(t.TextMuted)
		ui.Spacer(c)
		if turn {
			text := "to move"
			if a.opponent.Thinking() && color == a.opponent.Color() {
				text = "thinking…"
			} else if live.Check {
				text = "in check"
			}
			ui.Text(c, text).FontSize(12).Bold().TextColor(t.Accent)
		}
	})
}
