package tui_test

import (
	"regexp"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/delivery/tui"
	chess "github.com/gabrzb/chess-go/internal/domain"
)

func TestTurnIndicatorFollowsDisplayedHistory(t *testing.T) {
	m := app(t, tui.Options{StartGame: true})
	cursor := chess.A2
	m = appMove(t, m, &cursor, chess.Move{From: chess.E2, To: chess.E4})
	if !regexp.MustCompile(`│\s*Player 2\(B\)\s*│`).MatchString(ansi.Strip(m.View().Content)) {
		t.Fatal("black turn indicator missing from quadrant")
	}
	m = keys(t, m, 'h', tea.KeyLeft)
	contains(t, m, "REVIEW")
	if !regexp.MustCompile(`│\s*Player 1\(W\)\s*│`).MatchString(ansi.Strip(m.View().Content)) {
		t.Fatal("white turn indicator missing from reviewed quadrant")
	}
	m = keys(t, m, tea.KeyEscape, tea.KeyEscape, tea.KeyDown, tea.KeyDown, tea.KeyEnter, tea.KeyDown, tea.KeyEnter)
	contains(t, m, "Player 1(W) WINS")
	m = keys(t, m, tea.KeyDown, tea.KeyEnter)
	view := ansi.Strip(m.View().Content)
	if strings.Contains(view, "TURN") || !strings.Contains(view, "OVER") {
		t.Fatal("finished game showed an active turn")
	}
}
