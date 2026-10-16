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
	column := -1
	for _, line := range strings.Split(ansi.Strip(m.View().Content), "\n") {
		for _, label := range []string{"Jogador vs Jogador", "Jogador vs IA", "Sair"} {
			if index := strings.Index(line, label); index >= 0 {
				got := ansi.StringWidth(line[:index])
				if column >= 0 && column != got {
					t.Fatal("Menu options are not aligned")
				}
				column = got
			}
		}
	}
	m = keys(t, m, tea.KeyEnter)
	contains(t, m, "Vez das brancas")
	m = keys(t, m, 'd', 'd', 'd', 'd', tea.KeyEnter, 'w', 'w', tea.KeyEnter)
	contains(t, m, "Vez das pretas")
	contains(t, m, "1. e4")
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
	m = keys(t, m, 'f', 'h', tea.KeyLeft)
	contains(t, m, "h g f e d c b a")
	contains(t, m, "REVISÃO · posição 1/2")
	m = keys(t, m, tea.KeyEscape)
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 48, Height: 24})
	contains(t, m, "PARTIDA ENCERRADA")
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 28, Height: 12})
	contains(t, m, "Terminal pequeno demais.")
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	contains(t, m, "Último lance: e2 → e4")
	m = keys(t, m, tea.KeyEscape, tea.KeyEnter)
	contains(t, m, "Vez das brancas")
}

func appMove(t *testing.T, m tea.Model, cursor *chess.Square, move chess.Move) tea.Model {
	t.Helper()
	for _, to := range []chess.Square{move.From, move.To} {
		for _, axis := range []struct {
			delta              int
			positive, negative rune
		}{
			{to.Column() - cursor.Column(), tea.KeyRight, tea.KeyLeft},
			{to.Row() - cursor.Row(), tea.KeyUp, tea.KeyDown},
		} {
			key, steps := axis.positive, axis.delta
			if steps < 0 {
				key, steps = axis.negative, -steps
			}
			for range steps {
				m = keys(t, m, key)
			}
		}
		*cursor = to
		m = keys(t, m, tea.KeyEnter)
	}
	return m
}

func TestApplicationPromotionAndNotation(t *testing.T) {
	m := app(t, tui.Options{StartGame: true})
	cursor := chess.A2
	for _, move := range whitePromotionSetup {
		m = appMove(t, m, &cursor, move)
	}
	m = appMove(t, m, &cursor, chess.Move{From: chess.B7, To: chess.A8})
	contains(t, m, "PROMOÇÃO DE PEÃO")
	m = keys(t, m, '?', tea.KeyEscape, tea.KeyEscape)
	contains(t, m, "Vez das brancas")
	m = appMove(t, m, &cursor, chess.Move{From: chess.B7, To: chess.A8})
	m = keys(t, m, tea.KeyDown, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "Vez das pretas")
	contains(t, m, "bxa8=B")
	contains(t, m, "Último lance: b7 → a8")
}

func TestApplicationSpecialNotation(t *testing.T) {
	for _, tc := range []struct {
		name, moves, want string
		ended             bool
	}{
		{"castling", "e2e4 e7e5 g1f3 b8c6 f1c4 g8f6 e1g1", "O-O", false},
		{"disambiguation", "d2d4 a7a6 g1f3 a6a5 f3d2", "Nfd2", false},
		{"mate", "f2f3 e7e5 g2g4 d8h4", "Qh4#", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := app(t, tui.Options{StartGame: true})
			cursor := chess.A2
			for _, text := range strings.Fields(tc.moves) {
				from, err := chess.ParseSquare(text[:2])
				if err != nil {
					t.Fatal(err)
				}
				to, err := chess.ParseSquare(text[2:])
				if err != nil {
					t.Fatal(err)
				}
				m = appMove(t, m, &cursor, chess.Move{From: from, To: to})
			}
			if tc.ended {
				contains(t, m, "PRETAS VENCERAM")
				m = keys(t, m, tea.KeyDown, tea.KeyEnter)
			}
			contains(t, m, tc.want)
		})
	}
}

func TestApplicationOpeningAndConfiguration(t *testing.T) {
	finished, err := tui.NewApplication(t.Context(), tui.Options{})
	if err != nil {
		t.Fatal(err)
	}
	command := finished.Init()
	for frames := 0; command != nil && frames < 30; frames++ {
		finished, command = finished.Update(command())
	}
	if command != nil {
		t.Fatal("Opening animation did not finish")
	}
	finished = keys(t, finished, tea.KeyEnter)
	contains(t, finished, "Vez das brancas")

	m, err := tui.NewApplication(t.Context(), tui.Options{})
	if err != nil {
		t.Fatal(err)
	}
	tick := m.Init()
	if tick == nil {
		t.Fatal("Missing opening animation")
	}
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("Skipping opening scheduled a command")
	}
	contains(t, m, "Jogador vs Jogador")
	m = updateModel(t, m, tick())
	m = keys(t, m, 'a', tea.KeyDown, tea.KeyEnter, tea.KeyRight, tea.KeyDown, tea.KeyRight)
	contains(t, m, "Pretas")
	contains(t, m, "Difícil")
	m = keys(t, m, tea.KeyDown)
	m, search := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if search == nil {
		t.Fatal("Configured black game did not start AI")
	}
	m = updateModel(t, m, search())
	contains(t, m, "Você joga de pretas.")
	m = keys(t, m, tea.KeyEscape, tea.KeyDown, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "ENCERRAR PARTIDA?")
	m = keys(t, m, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "Jogador vs Jogador")
	m = keys(t, m, tea.KeyEnter)
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
