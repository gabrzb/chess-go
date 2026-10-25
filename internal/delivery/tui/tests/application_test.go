package tui_test

import (
	"regexp"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/delivery/tui"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/engine"
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
	contains(t, m, "Player vs Player")
	for _, line := range strings.Split(ansi.Strip(m.View().Content), "\n") {
		for _, label := range []string{"Player vs Player", "Player vs AI", "Quit"} {
			if index := strings.Index(line, label); index >= 0 {
				midpoint := ansi.StringWidth(line[:index]) + ansi.StringWidth(label)/2
				if midpoint < 39 || midpoint > 41 {
					t.Fatal("Menu option is not centered")
				}
			}
		}
	}
	m = keys(t, m, tea.KeyEnter)
	contains(t, m, "Turn: white")
	if strings.Count(ansi.Strip(m.View().Content), "▸") != 1 {
		t.Fatal("Expected exactly one board cursor")
	}
	view := ansi.Strip(m.View().Content)
	boardStart := strings.Index(view, "BOARD")
	if boardStart < 0 || strings.Index(view, "GAME") > boardStart || strings.Index(view, "HISTORY") > boardStart {
		t.Fatal("Match details must be above the board")
	}
	if ranks := regexp.MustCompile(`(?m)^\s*[1-8] `).FindAllString(view, -1); len(ranks) != chess.BoardSize {
		t.Fatal("Missing complete board")
	}
	m = keys(t, m, 'd', 'd', 'd', 'd', tea.KeyEnter, 'w', 'w', tea.KeyEnter)
	contains(t, m, "Turn: black")
	contains(t, m, "1. e4")
	m = keys(t, m, '?', tea.KeyEscape, tea.KeyEscape, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "DRAW OFFER")
	m = keys(t, m, tea.KeyEnter)
	contains(t, m, "Turn: black")
	m = keys(t, m, tea.KeyEscape, tea.KeyDown, tea.KeyEnter, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "DRAW")
	contains(t, m, "Draw by agreement.")
	m = keys(t, m, tea.KeyDown, tea.KeyEnter)
	before := m.View().Content
	m = keys(t, m, tea.KeyEnter)
	if m.View().Content != before {
		t.Fatal("Final board accepted a move")
	}
	m = keys(t, m, 'f', 'h', tea.KeyLeft)
	contains(t, m, "h g f e d c b a")
	contains(t, m, "REVIEW · position 1/2")
	m = keys(t, m, tea.KeyEscape)
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 48, Height: 24})
	contains(t, m, "GAME OVER")
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 28, Height: 12})
	contains(t, m, "Terminal too small.")
	m = updateModel(t, m, tea.WindowSizeMsg{Width: 80, Height: 24})
	contains(t, m, "e2 → e4")
	m = keys(t, m, tea.KeyEscape, tea.KeyEnter)
	contains(t, m, "Turn: white")
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
	contains(t, m, "PAWN PROMOTION")
	m = keys(t, m, '?', tea.KeyEscape, tea.KeyEscape)
	contains(t, m, "Turn: white")
	m = appMove(t, m, &cursor, chess.Move{From: chess.B7, To: chess.A8})
	m = keys(t, m, tea.KeyDown, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "Turn: black")
	contains(t, m, "bxa8=B")
	contains(t, m, "b7 → a8")
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
				contains(t, m, "Player 2(B) WINS")
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
	contains(t, finished, "Turn: white")

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
	contains(t, m, "Player vs Player")
	m = updateModel(t, m, tick())
	m = keys(t, m, 'a', tea.KeyDown, tea.KeyEnter, tea.KeyRight, tea.KeyDown, tea.KeyRight)
	contains(t, m, "Black")
	contains(t, m, "Hard")
	m = keys(t, m, tea.KeyDown)
	m, search := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if search == nil {
		t.Fatal("Configured black game did not start AI")
	}
	m = updateModel(t, m, search())
	contains(t, m, "You play as black.")
	m = keys(t, m, tea.KeyEscape, tea.KeyDown, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "END GAME?")
	m = keys(t, m, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "Player vs Player")
	m = keys(t, m, tea.KeyEnter)
	contains(t, m, "Turn: white")
}

func TestApplicationResignationNamesParticipants(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		m := app(t, tui.Options{Difficulty: engine.Easy, HumanColor: color, StartGame: true})
		m = keys(t, m, tea.KeyEscape, tea.KeyDown, tea.KeyEnter)
		contains(t, m, "AI will win.")
	}
	m := app(t, tui.Options{StartGame: true})
	m = keys(t, m, tea.KeyEscape, tea.KeyDown, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "Player 2(B) will win.")
}

func TestApplicationRejectsPreviousSessionAI(t *testing.T) {
	m := app(t, tui.Options{Difficulty: engine.Easy, HumanColor: chess.Black, StartGame: true})
	old := m.Init()
	if old == nil {
		t.Fatal("Missing initial AI search")
	}
	m = keys(t, m, tea.KeyEscape, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "AI will win.")
	m = keys(t, m, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "DEFEAT")
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
	contains(t, next, "Turn: black")
	contains(t, next, "You play as black.")
}
