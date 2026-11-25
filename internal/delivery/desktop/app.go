package desktop

import (
	"context"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/gabrzb/chess-go/internal/delivery/desktop/components"
	"github.com/gabrzb/chess-go/internal/delivery/desktop/pages"
)

type Options = components.Options
type App = pages.App

func New(ctx context.Context, options Options) (*App, error) {
	return pages.New(ctx, options)
}

// Run opens the window and blocks until the app quits.
func Run(options Options) error {
	options.Quit = mygo.App.Quit
	a, err := New(context.Background(), options)
	if err != nil {
		return err
	}
	mygo.App.WhenReady(func() {
		win := mygo.NewWindow(mygo.WindowOptions{
			Title:     "Chess Go",
			Width:     1100,
			Height:    780,
			MinWidth:  720,
			MinHeight: 560,
			Content:   ui.View(a.View),
		})
		a.SetUpdate(win.Update)
		a.StartAI()
	})
	return mygo.App.Run()
}
