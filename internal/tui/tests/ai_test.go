package tui_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/ai"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui"
)

type chooserFunc func(context.Context, chess.Position) (chess.Move, error)

func (f chooserFunc) ChooseMove(ctx context.Context, p chess.Position) (chess.Move, error) {
	return f(ctx, p)
}

func aiModel(t *testing.T, g *game.Game, chooser ai.MoveChooser, color chess.Color) (tea.Model, tea.Cmd) {
	t.Helper()
	m, err := tui.NewModelWithAI(context.Background(), g, chooser, color)
	if err != nil {
		t.Fatal(err)
	}
	init := m.Init()
	if init == nil {
		t.Fatal("Missing AI initialization")
	}
	return m.Update(init())
}

func TestAITurnsAndResults(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		g := gameWith(t)
		move := chess.Move{From: chess.E2, To: chess.E4}
		if color == chess.Black {
			move = chess.Move{From: chess.D7, To: chess.D5}
		}
		calls := 0
		chooser := chooserFunc(func(ctx context.Context, p chess.Position) (chess.Move, error) {
			calls++
			if p != g.Position() || p.SideToMove != color {
				t.Error("AI received an incorrect position")
			}
			return move, ctx.Err()
		})
		m, cmd := aiModel(t, g, chooser, color)
		if color == chess.Black {
			if cmd != nil {
				t.Fatal("AI started on the human turn")
			}
			m, cmd = submit(t, m, chess.Move{From: chess.E2, To: chess.E4})
		}
		if cmd == nil || calls != 0 || !strings.Contains(m.View().Content, "thinking...") {
			t.Fatal("Expected asynchronous search")
		}
		before, length := g.Position(), g.HistoryLen()
		m = moveCursorTo(t, m, move.From)
		m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
		assertHighlights(t, m, chess.NoSquare)
		result := cmd()
		assertUnchanged(t, g, before, length)
		m = updateModel(t, m, result)
		state, err := g.Snapshot()
		if err != nil || state.LastMove == nil || *state.LastMove != move || calls != 1 || g.HistoryLen() != length+1 || state.SideToMove == color {
			t.Fatal("AI move was not applied once")
		}
		completed := m
		m = updateModel(t, m, result)
		assertSameView(t, m, completed)
		if g.HistoryLen() != length+1 {
			t.Error("Duplicate result was applied")
		}
	}
	g := gameWith(t, chess.Move{From: chess.F2, To: chess.F3}, chess.Move{From: chess.E7, To: chess.E5}, chess.Move{From: chess.G2, To: chess.G4}, chess.Move{From: chess.D8, To: chess.H4})
	chooser := chooserFunc(func(context.Context, chess.Position) (chess.Move, error) {
		t.Error("Search started after checkmate")
		return chess.Move{}, nil
	})
	if _, cmd := aiModel(t, g, chooser, chess.White); cmd != nil {
		t.Error("Finished game scheduled a search")
	}
}

func TestAIReviewDiscardsCanceledResults(t *testing.T) {
	g := gameWith(t)
	var contexts []context.Context
	chooser := chooserFunc(func(ctx context.Context, _ chess.Position) (chess.Move, error) {
		contexts = append(contexts, ctx)
		return chess.Move{From: chess.E2, To: chess.E4}, nil
	})
	m, first := aiModel(t, g, chooser, chess.White)
	if first == nil {
		t.Fatal("Missing search")
	}
	m = updateModel(t, m, tea.KeyPressMsg{Code: 'h'})
	oldResult := first()
	if !errors.Is(contexts[0].Err(), context.Canceled) {
		t.Fatal("Review did not cancel the search")
	}
	reviewing := m
	m = updateModel(t, m, oldResult)
	assertSameView(t, m, reviewing)
	m, second := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if second == nil {
		t.Fatal("Leaving review did not restart search")
	}
	waiting := m
	m = updateModel(t, m, oldResult)
	assertSameView(t, m, waiting)
	if g.HistoryLen() != 1 {
		t.Fatal("Stale result changed the game")
	}
	m = updateModel(t, m, second())
	if g.HistoryLen() != 2 || strings.Contains(m.View().Content, "thinking...") {
		t.Error("Restarted search did not finish normally")
	}
}

func TestAIFailureAndRetry(t *testing.T) {
	for _, illegal := range []bool{false, true} {
		g := gameWith(t)
		calls := 0
		chooser := chooserFunc(func(context.Context, chess.Position) (chess.Move, error) {
			calls++
			if calls == 1 {
				if illegal {
					return chess.Move{From: chess.E2, To: chess.E5}, nil
				}
				return chess.Move{}, errors.New("search unavailable")
			}
			return chess.Move{From: chess.E2, To: chess.E4}, nil
		})
		m, cmd := aiModel(t, g, chooser, chess.White)
		m = updateModel(t, m, cmd())
		if !strings.Contains(m.View().Content, "AI failed:") || g.HistoryLen() != 1 {
			t.Fatal("Failure was hidden or changed the game")
		}
		m, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		if cmd == nil {
			t.Fatal("Missing retry command")
		}
		m = updateModel(t, m, cmd())
		if calls != 2 || g.HistoryLen() != 2 || strings.Contains(m.View().Content, "AI failed:") {
			t.Error("Retry failed to recover")
		}
	}
}

func TestHumanPromotionStartsAIOnlyAfterChoice(t *testing.T) {
	g := gameWith(t, whitePromotionSetup...)
	chooser := chooserFunc(func(context.Context, chess.Position) (chess.Move, error) {
		return chess.Move{From: chess.D7, To: chess.D5}, nil
	})
	m, cmd := aiModel(t, g, chooser, chess.Black)
	if cmd != nil {
		t.Fatal("AI started before the human played")
	}
	m, cmd = submit(t, m, chess.Move{From: chess.B7, To: chess.B8})
	if cmd != nil || g.HistoryLen() != len(whitePromotionSetup)+1 {
		t.Fatal("Pending promotion started an AI search")
	}
	m, cmd = m.Update(tea.KeyPressMsg{Code: 'n'})
	if cmd == nil {
		t.Fatal("Completed promotion did not schedule the AI reply")
	}
	m = updateModel(t, m, cmd())
	if g.HistoryLen() != len(whitePromotionSetup)+3 || strings.Contains(m.View().Content, "AI failed:") {
		t.Error("AI failed to reply after promotion")
	}
}

func TestEventLoopCancelsSearchOnQuit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	started, stopped := make(chan struct{}), make(chan struct{})
	chooser := chooserFunc(func(ctx context.Context, _ chess.Position) (chess.Move, error) {
		close(started)
		<-ctx.Done()
		close(stopped)
		return chess.Move{}, ctx.Err()
	})
	g := gameWith(t)
	m, err := tui.NewModelWithAI(ctx, g, chooser, chess.White)
	if err != nil {
		t.Fatal(err)
	}
	program := tea.NewProgram(m, tea.WithContext(ctx), tea.WithInput(nil), tea.WithOutput(io.Discard), tea.WithoutRenderer())
	done := make(chan error, 1)
	go func() { _, err := program.Run(); done <- err }()
	select {
	case <-started:
	case err := <-done:
		t.Fatalf("Event loop stopped before search: %v", err)
	case <-ctx.Done():
		t.Fatal("Search did not start")
	}
	program.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Quit failed during search: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("Search blocked the event loop")
	}
	select {
	case <-stopped:
	case <-ctx.Done():
		t.Fatal("Search was not canceled")
	}
	if g.HistoryLen() != 1 {
		t.Error("Quit altered the game")
	}
}
