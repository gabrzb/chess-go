package render

import (
	"strings"
	"testing"
)

func TestMenuStylesOnlyItsArrow(t *testing.T) {
	options := []string{"Continuar", "Desistir", "Voltar ao menu"}
	for choice := range options {
		menu := Menu(options, choice, 40)
		if strings.Count(menu, Paint(MenuCursor, Accent)) != 1 {
			t.Fatal("menu must contain exactly one focused arrow")
		}
		if strings.Contains(menu, "48;") || strings.Contains(menu, "\x1b[7m") {
			t.Fatal("menu focus must not paint a row background")
		}
		for _, option := range options {
			if strings.Contains(menu, Paint(option, Accent)) {
				t.Fatal("menu label inherited focus color")
			}
		}
	}
}
