package pages

import (
	"github.com/egoist/mygo/ui"

	"github.com/gabrzb/chess-go/internal/delivery/desktop/components"
)

func (a *App) setup(c *ui.Context) {
	t := c.Theme()
	label := func(s string) { ui.Text(c, s).FontSize(12).Bold().TextColor(t.TextMuted).Margin(8, 0, 0, 0) }
	ui.Column(c).Width(420).Padding(28).Gap(8).Radius(14).Background(t.Surface).Border(1, t.Border).
		Shadow(0, 16, 40, 0, ui.RGBA(0, 0, 0, 0.18)).Children(func() {
		ui.Text(c, "New game").FontSize(26).Bold()
		ui.Text(c, "Player vs AI").TextColor(t.TextMuted)
		label("YOUR PIECES")
		ui.Segmented(c, &a.Side, components.SideLabels...).FillEqually().Label("Your pieces")
		label("DIFFICULTY")
		ui.Segmented(c, &a.Level, components.LevelLabels...).FillEqually().Label("Difficulty")
		ui.Text(c, components.DifficultyDescription(a.Level)).
			FontSize(12).TextColor(t.TextMuted)
		ui.Row(c).Gap(10).Margin(20, 0, 0, 0).Children(func() {
			if ui.Button(c, "Back").Height(40).Clicked() {
				a.Screen = components.MenuScreen
			}
			if ui.PrimaryButton(c, "Start game").Height(40).Grow(1).AutoFocus().Clicked() {
				a.StartGame(true)
			}
		})
	})
}
