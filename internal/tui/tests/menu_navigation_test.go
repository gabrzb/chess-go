package tui_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/tui"
)

func TestMenuActionsMatchDisplayedOptions(t *testing.T) {
	for _, difficulty := range []ai.Difficulty{"", ai.Easy} {
		m := app(t, tui.Options{StartGame: true, Difficulty: difficulty, HumanColor: chess.White})
		m = keys(t, m, tea.KeyEscape, tea.KeyDown)
		if difficulty == "" {
			m = keys(t, m, tea.KeyDown)
		}
		contains(t, m, "▸ Desistir")
		m = keys(t, m, tea.KeyEnter)
		contains(t, m, "CONFIRMAR DESISTÊNCIA?")
	}
}

func TestConfigurationCyclesInBothDirections(t *testing.T) {
	m := app(t, tui.Options{})
	m = keys(t, m, tea.KeyDown, tea.KeyEnter, tea.KeyLeft)
	contains(t, m, "Aleatório")
	m = keys(t, m, tea.KeyRight)
	contains(t, m, "Brancas")
	m = keys(t, m, tea.KeyDown, tea.KeyRight, tea.KeyRight, tea.KeyRight)
	contains(t, m, "Normal")
	m = keys(t, m, tea.KeyLeft)
	contains(t, m, "Fácil")
}
