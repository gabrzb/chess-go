package tui_test

import (
	"context"
	"errors"
	"io"
	"reflect"
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

func (f chooserFunc) ChooseMove(ctx context.Context, position chess.Position) (chess.Move, error) {
	return f(ctx, position)
}

func aiModel(t *testing.T, g *game.Game, chooser ai.MoveChooser, color chess.Color) tea.Model {
	t.Helper()
	m, err := tui.NewModelWithAI(context.Background(), g, chooser, color)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func initializeAI(t *testing.T, m tea.Model) (tea.Model, tea.Cmd) {
	t.Helper()
	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Expected AI initialization command")
	}
	return m.Update(cmd())
}

func TestAITurnsAndDuplicateResults(t *testing.T) {
	for _, color := range []chess.Color{chess.White, chess.Black} {
		name := map[chess.Color]string{chess.White: "White", chess.Black: "Black"}[color]
		t.Run(name, func(t *testing.T) {
			g := game.NewGame()
			move := chess.Move{From: chess.E2, To: chess.E4}
			if color == chess.Black {
				move = chess.Move{From: chess.D7, To: chess.D5}
			}
			calls := 0
			chooser := chooserFunc(func(ctx context.Context, position chess.Position) (chess.Move, error) {
				calls++
				if position != g.Position() || position.SideToMove != color {
					t.Error("Expected a copy of the current complete position for the AI color")
				}
				return move, ctx.Err()
			})
			m, cmd := initializeAI(t, aiModel(t, g, chooser, color))
			if color == chess.Black {
				if cmd != nil {
					t.Fatal("Expected White human to play before Black AI")
				}
				m = moveCursorTo(t, m, chess.E2)
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				m = moveCursorTo(t, m, chess.E4)
				m, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			}
			if cmd == nil || !strings.Contains(m.View().Content, "thinking...") || calls != 0 {
				t.Fatal("Expected a pending asynchronous AI command, without running search in Update")
			}
			before, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			length := g.HistoryLen()
			m = moveCursorTo(t, m, move.From)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			assertSelectionHighlights(t, m, chess.NoSquare, nil, move.From)
			result := cmd()
			assertGameUnchanged(t, g, before, length)
			m = updateModel(t, m, result)
			after, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			if calls != 1 || after.LastMove == nil || *after.LastMove != move || g.HistoryLen() != length+1 || after.SideToMove == color {
				t.Fatal("Expected one AI move, recorded exactly once, followed by the human turn")
			}
			assertSelectionHighlights(t, m, chess.NoSquare, nil, move.To)
			if strings.Contains(m.View().Content, "thinking...") {
				t.Error("Expected completed search to clear the thinking indicator")
			}
			view := m.View()
			m = updateModel(t, m, result)
			if !reflect.DeepEqual(m.View(), view) {
				t.Error("Expected a duplicate result to be ignored")
			}
			assertGameUnchanged(t, g, after, length+1)
		})
	}
}

func TestAIReviewCancellationAndStaleResults(t *testing.T) {
	for _, exit := range []tea.KeyPressMsg{{Code: 'h', Text: "h"}, {Code: tea.KeyEscape}} {
		t.Run(exit.String(), func(t *testing.T) {
			g := game.NewGame()
			var contexts []context.Context
			chooser := chooserFunc(func(ctx context.Context, _ chess.Position) (chess.Move, error) {
				contexts = append(contexts, ctx)
				return chess.Move{From: chess.E2, To: chess.E4}, nil
			})
			m, first := initializeAI(t, aiModel(t, g, chooser, chess.White))
			if first == nil {
				t.Fatal("Expected initial AI search")
			}
			m = updateModel(t, m, tea.KeyPressMsg{Code: 'h', Text: "h"})
			oldResult := first()
			if !errors.Is(contexts[0].Err(), context.Canceled) {
				t.Error("Expected entering review to cancel the search context")
			}
			view := m.View()
			m = updateModel(t, m, oldResult)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			if !reflect.DeepEqual(m.View(), view) || g.HistoryLen() != 1 {
				t.Fatal("Expected review to ignore late AI results and move attempts")
			}
			m, second := m.Update(exit)
			if second == nil {
				t.Fatal("Expected leaving review to restart the AI turn")
			}
			view = m.View()
			m = updateModel(t, m, oldResult)
			if !reflect.DeepEqual(m.View(), view) || g.HistoryLen() != 1 {
				t.Fatal("Expected old request ID to preserve the newer pending search")
			}
			m = updateModel(t, m, second())
			if g.HistoryLen() != 2 || !errors.Is(contexts[1].Err(), context.Canceled) {
				t.Error("Expected only the restarted search to play, and its context to be released")
			}
		})
	}
}

func TestAIFailureAndRetry(t *testing.T) {
	for _, illegal := range []bool{false, true} {
		name := "search error"
		if illegal {
			name = "illegal result"
		}
		t.Run(name, func(t *testing.T) {
			g := game.NewGame()
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
			m, cmd := initializeAI(t, aiModel(t, g, chooser, chess.White))
			m = updateModel(t, m, cmd())
			if !strings.Contains(m.View().Content, "AI failed:") || g.HistoryLen() != 1 {
				t.Fatal("Expected a visible error without changing the game")
			}
			m, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
			if cmd == nil {
				t.Fatal("Expected Enter to retry the AI turn")
			}
			m = updateModel(t, m, cmd())
			if g.HistoryLen() != 2 || calls != 2 || strings.Contains(m.View().Content, "AI failed:") {
				t.Error("Expected successful retry to play once and clear the error")
			}
		})
	}
}

func TestQuitCancelsAISearch(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{{Code: 'q', Text: "q"}, {Code: 'c', Mod: tea.ModCtrl}} {
		t.Run(key.String(), func(t *testing.T) {
			g := game.NewGame()
			chooser := chooserFunc(func(ctx context.Context, _ chess.Position) (chess.Move, error) {
				if !errors.Is(ctx.Err(), context.Canceled) {
					t.Error("Expected quitting to cancel the active search")
				}
				return chess.Move{}, ctx.Err()
			})
			m, search := initializeAI(t, aiModel(t, g, chooser, chess.White))
			updated, quit := m.Update(key)
			if quit == nil {
				t.Fatal("Expected quit command")
			}
			if _, ok := quit().(tea.QuitMsg); !ok {
				t.Fatal("Expected quit message")
			}
			updateModel(t, updated, search())
			if g.HistoryLen() != 1 {
				t.Error("Expected quitting to preserve the game")
			}
		})
	}
}

func TestRealAIPromotionAndFinishedGame(t *testing.T) {
	p, err := ai.NewPlayer(ai.Easy, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, finished := range []bool{false, true} {
		name := "automatic promotion"
		setup := whitePromotionSetup
		if finished {
			name = "finished game"
			setup = []chess.Move{
				{From: chess.F2, To: chess.F3}, {From: chess.E7, To: chess.E5},
				{From: chess.G2, To: chess.G4}, {From: chess.D8, To: chess.H4},
			}
		}
		t.Run(name, func(t *testing.T) {
			g := game.NewGame()
			for _, move := range setup {
				if _, err := g.Play(move); err != nil {
					t.Fatal(err)
				}
			}
			m, cmd := initializeAI(t, aiModel(t, g, p, chess.White))
			if finished {
				if cmd != nil || !strings.Contains(m.View().Content, "Checkmate!") || strings.Contains(m.View().Content, "retry AI") {
					t.Error("Expected no search or retry hint after checkmate")
				}
				return
			}
			if cmd == nil {
				t.Fatal("Expected promotion search")
			}
			m = updateModel(t, m, cmd())
			snapshot, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			if snapshot.LastMove == nil || snapshot.LastMove.Promotion == chess.NoPiece || g.HistoryLen() != len(setup)+2 {
				t.Fatal("Expected AI to include promotion choice in one atomic move")
			}
			if strings.Contains(m.View().Content, "Promotion:") {
				t.Error("Expected AI promotion without asking the human to choose")
			}
		})
	}
}

func TestAIConfigurationValidation(t *testing.T) {
	g := game.NewGame()
	chooser := chooserFunc(func(context.Context, chess.Position) (chess.Move, error) { return chess.Move{}, nil })
	for _, tc := range []struct {
		name    string
		ctx     context.Context
		game    *game.Game
		chooser ai.MoveChooser
		color   chess.Color
	}{
		{"missing context", nil, g, chooser, chess.Black},
		{"missing game", context.Background(), nil, chooser, chess.Black},
		{"missing chooser", context.Background(), g, nil, chess.Black},
		{"invalid color", context.Background(), g, chooser, chess.NoColor},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m, err := tui.NewModelWithAI(tc.ctx, tc.game, tc.chooser, tc.color)
			if !errors.Is(err, tui.ErrInvalidAIConfig) || m != nil {
				t.Errorf("Expected invalid configuration without a model, got %v", err)
			}
		})
	}
}

func TestHumanPromotionStartsAIAfterChoice(t *testing.T) {
	g := game.NewGame()
	for _, move := range whitePromotionSetup {
		if _, err := g.Play(move); err != nil {
			t.Fatal(err)
		}
	}
	chooser := chooserFunc(func(context.Context, chess.Position) (chess.Move, error) {
		return chess.Move{From: chess.D7, To: chess.D5}, nil
	})
	m, cmd := initializeAI(t, aiModel(t, g, chooser, chess.Black))
	if cmd != nil {
		t.Fatal("Expected the human promotion turn first")
	}
	m = moveCursorTo(t, m, chess.B7)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	m = moveCursorTo(t, m, chess.B8)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updateModel(t, m, tea.KeyPressMsg{Code: 'h', Text: "h"})
	if !strings.Contains(m.View().Content, "Promotion:") || g.HistoryLen() != len(whitePromotionSetup)+1 {
		t.Fatal("Expected pending promotion to block review and defer the AI turn")
	}
	m, cmd = m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	if cmd == nil || g.HistoryLen() != len(whitePromotionSetup)+2 {
		t.Fatal("Expected choosing the promotion piece to schedule the AI reply")
	}
	m = updateModel(t, m, cmd())
	if g.HistoryLen() != len(whitePromotionSetup)+3 || strings.Contains(m.View().Content, "AI failed:") {
		t.Error("Expected the AI reply after the completed human promotion")
	}
}

func TestEventLoopRemainsResponsiveDuringSearch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	started := make(chan struct{})
	stopped := make(chan struct{})
	chooser := chooserFunc(func(searchContext context.Context, _ chess.Position) (chess.Move, error) {
		close(started)
		<-searchContext.Done()
		close(stopped)
		return chess.Move{}, searchContext.Err()
	})
	g := game.NewGame()
	m, err := tui.NewModelWithAI(ctx, g, chooser, chess.White)
	if err != nil {
		t.Fatal(err)
	}
	program := tea.NewProgram(m, tea.WithContext(ctx), tea.WithInput(nil), tea.WithOutput(io.Discard), tea.WithoutRenderer())
	done := make(chan error, 1)
	go func() {
		_, err := program.Run()
		done <- err
	}()
	select {
	case <-started:
	case err := <-done:
		t.Fatalf("Event loop stopped before AI search: %v", err)
	case <-ctx.Done():
		t.Fatal("AI search did not start")
	}
	program.Send(tea.KeyPressMsg{Code: 'q', Text: "q"})
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Expected normal quit while AI was computing: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("AI search blocked the event loop")
	}
	select {
	case <-stopped:
	case <-ctx.Done():
		t.Fatal("Quit did not release the AI search")
	}
	if g.HistoryLen() != 1 {
		t.Error("Expected quitting during search to preserve the live game")
	}
}
