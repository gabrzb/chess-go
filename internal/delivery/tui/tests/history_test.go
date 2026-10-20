package tui_test

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/delivery/tui"
	chess "github.com/gabrzb/chess-go/internal/domain"
)

func TestHistoryReviewAndResume(t *testing.T) {
	g := gameWith(t)
	m := tui.NewModel(g)
	views := []string{boardText(t, m)}
	for _, move := range []chess.Move{{From: chess.E2, To: chess.E4}, {From: chess.D7, To: chess.D5}, {From: chess.E4, To: chess.D5}} {
		if _, err := g.Play(move); err != nil {
			t.Fatal(err)
		}
		views = append(views, boardText(t, m))
	}
	before, length := g.Position(), g.HistoryLen()
	m = updateModel(t, m, tea.KeyPressMsg{Code: 'h'})
	for index := length - 1; index >= 0; index-- {
		if boardText(t, m) != views[index] || !strings.Contains(m.View().Content, fmt.Sprintf("position %d/%d", index+1, length)) {
			t.Fatalf("Wrong historical view at %d", index)
		}
		m = updateModel(t, m, tea.KeyPressMsg{Code: 'a'})
	}
	if boardText(t, m) != views[0] {
		t.Error("Review moved before initial position")
	}
	initial := m
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	assertSameView(t, m, initial)
	for range length + 1 {
		m = updateModel(t, m, tea.KeyPressMsg{Code: 'd'})
	}
	if boardText(t, m) != views[length-1] {
		t.Error("Review moved beyond latest position")
	}
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyLeft})
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEscape})
	assertUnchanged(t, g, before, length)
	if boardText(t, m) != views[length-1] {
		t.Error("Escape did not restore the live position")
	}
	m = updateModel(t, m, tea.KeyPressMsg{Code: 'h'})
	m = updateModel(t, m, tea.KeyPressMsg{Code: 'h'})
	var cmd tea.Cmd
	m, cmd = submit(t, m, chess.Move{From: chess.E7, To: chess.E5})
	if cmd != nil || g.HistoryLen() != length+1 {
		t.Error("Resuming play failed or truncated the history")
	}
}
