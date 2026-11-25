package components

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

func (a *App) Sidebar(c *ui.Context) {
	t := c.Theme()
	live, _ := a.game.Snapshot()
	ui.Row(c).Gap(8).Children(func() {
		if ui.Button(c, "☰  Menu").Clicked() {
			a.openGameMenu()
		}
		ui.Spacer(c)
		if ui.Button(c, "⇅").Label("Flip board").Width(40).Clicked() {
			a.Flipped = !a.Flipped
		}
	})

	ui.Column(c).Padding(14).Gap(6).Radius(10).Background(t.Surface).Border(1, t.Border).Children(func() {
		ui.Text(c, a.status(live)).Bold()
		if a.review >= 0 {
			ui.Text(c, "Reviewing an earlier position — press → to return").FontSize(12).TextColor(t.Warning)
		}
		if a.message != "" {
			ui.Text(c, a.message).FontSize(12).TextColor(t.Danger)
		}
	})

	ui.Column(c).Grow(1).Padding(14).Gap(8).Radius(10).Background(t.Surface).Border(1, t.Border).Children(func() {
		ui.Row(c).Gap(4).Children(func() {
			ui.Text(c, "Moves").Bold().Grow(1)
			last := a.game.HistoryLen() - 1
			at := a.review
			if at < 0 {
				at = last
			}
			if ui.Button(c, "‹").Label("Previous move").Width(32).Disabled(at == 0).Clicked() {
				a.step(-1)
			}
			if ui.Button(c, "›").Label("Next move").Width(32).Disabled(at == last).Clicked() {
				a.step(1)
			}
		})
		ui.Scroll(c).Grow(1).Children(func() {
			if len(a.moveHistory) == 0 {
				ui.Text(c, "No moves yet").TextColor(t.TextMuted)
			}
			for i := 0; i < len(a.moveHistory); i += 2 {
				ui.Row(c.Key(i)).Gap(4).Children(func() {
					ui.Text(c, fmt.Sprintf("%d.", i/2+1)).Width(30).TextColor(t.TextMuted)
					for j := i; j < min(i+2, len(a.moveHistory)); j++ {
						a.moveChip(c, j)
					}
				})
			}
		})
	})

	if a.over() && a.resultDismissed {
		ui.Row(c).Gap(8).Children(func() {
			if ui.Button(c, "Main menu").Height(40).Clicked() {
				a.toMenu()
			}
			if ui.PrimaryButton(c, "Play again").Height(40).Grow(1).Clicked() {
				a.StartGame(a.vsAI)
			}
		})
	}
}

func (a *App) moveChip(c *ui.Context, index int) {
	t := c.Theme()
	current := a.review == index+1 || (a.review < 0 && index == len(a.moveHistory)-1)
	chip := ui.Box(c.Key(index)).Width(96).Padding(3, 8).Radius(5).Focusable().Cursor(ui.CursorPointer)
	if current {
		chip.Background(t.Accent)
	}
	if chip.Clicked() {
		a.setReview(index + 1)
	}
	chip.Children(func() {
		text := ui.Text(c, a.moveHistory[index])
		if current {
			text.TextColor(t.AccentText)
		}
	})
}
