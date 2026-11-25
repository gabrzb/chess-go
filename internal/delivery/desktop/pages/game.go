package pages

import (
	"github.com/egoist/mygo/ui"

	"github.com/gabrzb/chess-go/internal/delivery/desktop/components"
	chess "github.com/gabrzb/chess-go/internal/domain"
)

const (
	gamePadding   = 24
	sidebarWidth  = 280
	playerGap     = 8
	minSquareSize = 36
)

func (a *App) game(c *ui.Context) {
	width, height := c.Size()
	boardWidth := width - 3*gamePadding - sidebarWidth
	boardHeight := height - 2*(gamePadding+components.PlayerBarHeight+playerGap)
	size := max(float32(int(min(boardWidth, boardHeight)/chess.BoardSize)), minSquareSize)
	top, bottom := chess.Black, chess.White
	if a.Flipped {
		top, bottom = bottom, top
	}
	ui.Row(c).Fill().Padding(gamePadding).Gap(gamePadding).Children(func() {
		ui.Column(c).Grow(1).Center().Gap(playerGap).Children(func() {
			a.Player(c, top, size*chess.BoardSize)
			a.Board(c, size)
			a.Player(c, bottom, size*chess.BoardSize)
		})
		ui.Column(c).Width(sidebarWidth).FillHeight().Gap(12).Children(func() { a.Sidebar(c) })
	})
}
