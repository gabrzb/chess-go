package pages

import (
	"context"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/gabrzb/chess-go/internal/delivery/desktop/components"
)

const screenTransitionDuration = 220 * time.Millisecond

type App struct{ *components.App }

func New(ctx context.Context, options components.Options) (*App, error) {
	app, err := components.New(ctx, options)
	if err != nil {
		return nil, err
	}
	return &App{App: app}, nil
}

func (a *App) View(c *ui.Context) {
	t := c.Theme()
	root := ui.Box(c.Key(a.Screen)).Fill().Background(t.Background).Transition(ui.ElementTransition{Enter: &ui.Motion{Y: 16}, Duration: screenTransitionDuration})
	switch a.Screen {
	case components.MenuScreen:
		root.Center().Children(func() { a.menu(c) })
	case components.SetupScreen:
		if c.Shortcut(0, ui.KeyEscape) {
			a.Screen = components.MenuScreen
		}
		root.Center().Children(func() { a.setup(c) })
	case components.GameScreen:
		a.GameShortcuts(c)
		root.Children(func() { a.game(c) })
		a.Dialogs(c)
	}
}
