package application

import (
	"context"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestTurnIndicatorFollowsDisplayedHistory(t *testing.T) {
	model, err := New(context.Background(), Options{StartGame: true, NoAnimations: true})
	if err != nil {
		t.Fatal(err)
	}
	m := model.(application)
	if _, err := m.match.Game().Play(chess.Move{From: chess.E2, To: chess.E4}); err != nil {
		t.Fatal(err)
	}
	m.afterMatch()
	if !strings.Contains(ansi.Strip(m.View().Content), "PRETAS") {
		t.Fatal("indicator did not follow move")
	}
	for _, key := range []rune{'h', tea.KeyLeft} {
		next, _ := m.updateBoard(tea.KeyPressMsg{Code: key})
		m = next.(application)
	}
	view := ansi.Strip(m.View().Content)
	if !strings.Contains(view, "REVISÃO") || !strings.Contains(view, "BRANCAS") {
		t.Fatal("indicator followed live game instead of displayed position")
	}
	m.match.EndReview()
	if err := m.match.Game().Resign(chess.White); err != nil {
		t.Fatal(err)
	}
	m.show(finalBoardScreen)
	view = ansi.Strip(m.View().Content)
	if strings.Contains(view, "VEZ DAS") || !strings.Contains(view, "ENCERRADA") {
		t.Fatal("finished game showed an active turn")
	}
}
