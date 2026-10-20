package render_test

import (
	"github.com/gabrzb/chess-go/internal/tui/render"
	"strings"
	"testing"
)

func TestMenuStylesOnlyItsArrow(t *testing.T) {
	options := []string{"Continuar", "Desistir", "Voltar ao menu"}
	for choice := range options {
		menu := render.Menu(options, choice, 40)
		if strings.Count(menu, render.Paint(render.MenuCursor, render.Accent)) != 1 {
			t.Fatal("menu must contain exactly one focused arrow")
		}
		if strings.Contains(menu, "48;") || strings.Contains(menu, "\x1b[7m") {
			t.Fatal("menu focus must not paint a row background")
		}
		for _, option := range options {
			if strings.Contains(menu, render.Paint(option, render.Accent)) {
				t.Fatal("menu label inherited focus color")
			}
		}
	}
}
