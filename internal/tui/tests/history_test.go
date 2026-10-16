package tui_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui"
)

func assertReviewPosition(t *testing.T, m tea.Model, expected tea.View, index, length int) {
	t.Helper()
	got := strings.Split(m.View().Content, "\n\n")
	want := strings.Split(expected.Content, "\n\n")
	if len(got) < 4 || len(want) < 3 {
		t.Fatalf("Expected board, status and controls, got %q", m.View().Content)
	}
	if got[1] != want[1] || strings.SplitN(got[2], "\n", 2)[0] != want[2] {
		t.Error("Expected reviewed board and status to match the historical position")
	}
	indicator := fmt.Sprintf("position %d/%d", index+1, length)
	if !strings.Contains(got[2], indicator) || !strings.Contains(got[3], "previous") || !strings.Contains(got[3], "next") {
		t.Errorf("Expected review indicator %q and navigation controls, got %q", indicator, m.View().Content)
	}
	if strings.Contains(got[3], "select/play") {
		t.Error("Expected review controls without move execution")
	}
}

func TestHistoryNavigation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		moves []chess.Move
	}{
		{name: "initial position only"},
		{"capture", []chess.Move{
			{From: chess.E2, To: chess.E4},
			{From: chess.D7, To: chess.D5},
			{From: chess.E4, To: chess.D5},
		}},
		{"check", []chess.Move{
			{From: chess.E2, To: chess.E4},
			{From: chess.D7, To: chess.D5},
			{From: chess.F1, To: chess.B5},
		}},
		{"checkmate", []chess.Move{
			{From: chess.F2, To: chess.F3},
			{From: chess.E7, To: chess.E5},
			{From: chess.G2, To: chess.G4},
			{From: chess.D8, To: chess.H4},
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			m := tui.NewModel(g)
			views := []tea.View{m.View()}
			snapshots := make([]game.Snapshot, 0, len(tc.moves)+1)
			initial, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			snapshots = append(snapshots, initial)
			for _, move := range tc.moves {
				result, err := g.Play(move)
				if err != nil {
					t.Fatal(err)
				}
				snapshots = append(snapshots, result.State)
				views = append(views, m.View())
			}
			length := len(views)
			m = updateModel(t, m, tea.KeyPressMsg{Code: 'h', Text: "h"})
			assertReviewPosition(t, m, views[length-1], length-1, length)
			for index := length - 2; index >= 0; index-- {
				key := tea.KeyPressMsg{Code: tea.KeyLeft}
				if index%2 == 0 {
					key = tea.KeyPressMsg{Code: 'a', Text: "a"}
				}
				m = updateModel(t, m, key)
				assertReviewPosition(t, m, views[index], index, length)
			}
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyLeft})
			assertReviewPosition(t, m, views[0], 0, length)
			for index := 1; index < length; index++ {
				key := tea.KeyPressMsg{Code: tea.KeyRight}
				if index%2 == 0 {
					key = tea.KeyPressMsg{Code: 'd', Text: "d"}
				}
				m = updateModel(t, m, key)
				assertReviewPosition(t, m, views[index], index, length)
			}
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyRight})
			assertReviewPosition(t, m, views[length-1], length-1, length)
			assertGameUnchanged(t, g, snapshots[length-1], length)
			for index, want := range snapshots {
				got, err := g.HistorySnapshot(index)
				if err != nil || !reflect.DeepEqual(got, want) {
					t.Errorf("Expected historical position %d to remain unchanged: %v", index, err)
				}
			}
		})
	}
}

func TestHistoryReturnsToLiveGame(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{{Code: 'h', Text: "h"}, {Code: tea.KeyEscape}} {
		t.Run(key.String(), func(t *testing.T) {
			g := game.NewGame()
			m := tui.NewModel(g)
			for _, move := range []chess.Move{{From: chess.E2, To: chess.E4}, {From: chess.D7, To: chess.D5}} {
				if _, err := g.Play(move); err != nil {
					t.Fatal(err)
				}
			}
			m = moveCursorTo(t, m, chess.E4)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			m = moveCursorTo(t, m, chess.D4)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			if !strings.Contains(m.View().Content, chess.ErrIllegalMove.Error()) {
				t.Fatal("Expected selected piece and invalid move message before review")
			}
			before, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			length := g.HistoryLen()
			m = updateModel(t, m, tea.KeyPressMsg{Code: 'h', Text: "h"})
			assertSelectionHighlights(t, m, chess.NoSquare, nil, chess.D4)
			if strings.Contains(m.View().Content, chess.ErrIllegalMove.Error()) {
				t.Error("Expected entering review to clear the error")
			}
			for range length - 1 {
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyLeft})
			}
			m = updateModel(t, m, key)
			expected := moveCursorTo(t, tui.NewModel(g), chess.D4)
			if !reflect.DeepEqual(m.View(), expected.View()) {
				t.Error("Expected return to the live game with selection cleared")
			}
			assertGameUnchanged(t, g, before, length)
			move := chess.Move{From: chess.E4, To: chess.D5}
			m = moveCursorTo(t, m, move.From)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			m = moveCursorTo(t, m, move.To)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			after, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			if after.LastMove == nil || *after.LastMove != move || g.HistoryLen() != length+1 {
				t.Fatal("Expected resumed move to extend the live game exactly once")
			}
			previous, err := g.HistorySnapshot(length - 1)
			if err != nil || !reflect.DeepEqual(previous, before) {
				t.Error("Expected continued play to preserve the previous live position")
			}
			liveView := m.View()
			m = updateModel(t, m, tea.KeyPressMsg{Code: 'h', Text: "h"})
			assertReviewPosition(t, m, liveView, length, length+1)
		})
	}
}

func TestHistoryIgnoresOtherEvents(t *testing.T) {
	m, g := modelAt(t, chess.A2)
	m = updateModel(t, m, tea.KeyPressMsg{Code: 'h', Text: "h"})
	before, err := g.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	view := m.View()
	for _, msg := range []tea.Msg{
		tea.KeyPressMsg{Code: tea.KeyEnter},
		tea.KeyPressMsg{Code: tea.KeyUp},
		tea.KeyPressMsg{Code: tea.KeyDown},
		tea.KeyPressMsg{Code: 'w', Text: "w"},
		tea.KeyPressMsg{Code: 's', Text: "s"},
		tea.KeyPressMsg{Code: 'x', Text: "x"},
		tea.KeyReleaseMsg{Code: 'h', Text: "h"},
		tea.KeyReleaseMsg{Code: tea.KeyEscape},
		tea.KeyReleaseMsg{Code: 'q', Text: "q"},
		tea.WindowSizeMsg{},
	} {
		m = updateModel(t, m, msg)
		if !reflect.DeepEqual(m.View(), view) {
			t.Errorf("Expected ignored review event %v to preserve the view", msg)
		}
		assertGameUnchanged(t, g, before, 1)
	}
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEscape})
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	assertSelectionHighlights(t, m, chess.A2, []chess.Square{chess.A3, chess.A4}, chess.A2)
}

func TestQuitDuringReview(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{{Code: 'q', Text: "q"}, {Code: 'c', Mod: tea.ModCtrl}} {
		t.Run(key.String(), func(t *testing.T) {
			m, g := modelAt(t, chess.A2)
			m = updateModel(t, m, tea.KeyPressMsg{Code: 'h', Text: "h"})
			before, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			updated, cmd := m.Update(key)
			if cmd == nil {
				t.Fatal("Expected quit command during review")
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Fatal("Expected command to produce quit message")
			}
			if !reflect.DeepEqual(updated.View(), m.View()) {
				t.Error("Expected quit request to preserve the reviewed view")
			}
			assertGameUnchanged(t, g, before, 1)
		})
	}
}
