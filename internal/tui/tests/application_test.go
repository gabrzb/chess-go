package tui_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/tui"
)

func app(t *testing.T, options tui.Options) tea.Model {
	t.Helper()
	options.NoAnimations = true
	m, err := tui.NewApplication(t.Context(), options)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func keys(t *testing.T, m tea.Model, codes ...rune) tea.Model {
	t.Helper()
	for _, code := range codes {
		m = updateModel(t, m, tea.KeyPressMsg{Code: code})
	}
	return m
}

func contains(t *testing.T, m tea.Model, text string) {
	t.Helper()
	if !strings.Contains(ansi.Strip(m.View().Content), text) {
		t.Fatalf("Expected %q in:\n%s", text, ansi.Strip(m.View().Content))
	}
}

func TestApplicationLocalGameAndResult(t *testing.T) {
	m := app(t, tui.Options{})
	contains(t, m, "Jogador vs Jogador")
	m = keys(t, m, tea.KeyEnter)
	contains(t, m, "Vez das brancas")
	m = keys(t, m, 'd', 'd', 'd', 'd', tea.KeyEnter, 'w', 'w', tea.KeyEnter)
	contains(t, m, "Vez das pretas")
	m = keys(t, m, '?', tea.KeyEscape, tea.KeyEscape, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "OFERTA DE EMPATE")
	m = keys(t, m, tea.KeyEnter)
	contains(t, m, "Vez das pretas")
	m = keys(t, m, tea.KeyEscape, tea.KeyDown, tea.KeyEnter, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "EMPATE")
	contains(t, m, "Empate por acordo.")
	m = keys(t, m, tea.KeyDown, tea.KeyEnter)
	before := m.View().Content
	m = keys(t, m, tea.KeyEnter)
	if m.View().Content != before {
		t.Fatal("Final board accepted a move")
	}
	m = keys(t, m, tea.KeyEscape, tea.KeyEnter)
	contains(t, m, "Vez das brancas")
}

func TestApplicationRejectsPreviousSessionAI(t *testing.T) {
	m := app(t, tui.Options{Difficulty: ai.Easy, HumanColor: chess.Black, StartGame: true})
	old := m.Init()
	if old == nil {
		t.Fatal("Missing initial AI search")
	}
	m = keys(t, m, tea.KeyEscape, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "As brancas serão as vencedoras.")
	m = keys(t, m, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "DERROTA")
	next, current := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if current == nil {
		t.Fatal("Rematch did not restart AI")
	}
	before := next.View().Content
	next = updateModel(t, next, old())
	if next.View().Content != before {
		t.Fatal("Old session changed the new session")
	}
	next = updateModel(t, next, current())
	contains(t, next, "Vez das pretas")
	contains(t, next, "Você joga de pretas.")
}
