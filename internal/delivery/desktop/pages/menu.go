package pages

import (
	"github.com/egoist/mygo/ui"

	"github.com/gabrzb/chess-go/internal/delivery/desktop/components"
)

func (a *App) menu(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Width(340).AlignItems(ui.Center).Gap(10).Children(func() {
		ui.Text(c, "♞").FontSize(96).TextColor(t.Accent)
		ui.Text(c, "CHESS GO").FontSize(44).Bold().LetterSpacing(6)
		ui.Text(c, "Two players on one board, or a game against the AI").TextColor(t.TextMuted).TextAlign(ui.Center)
		ui.Column(c).FillWidth().Gap(10).Margin(28, 0, 0, 0).FocusGroup(ui.Vertical).Children(func() {
			if ui.PrimaryButton(c, "Player vs Player").Height(46).FontSize(16).AutoFocus().Clicked() {
				a.StartGame(false)
			}
			if ui.Button(c, "Player vs AI").Height(46).FontSize(16).Clicked() {
				a.Screen = components.SetupScreen
			}
			if ui.Button(c, "Quit").Height(46).FontSize(16).Clicked() && a.Quit != nil {
				a.Quit()
			}
		})
	})
}
