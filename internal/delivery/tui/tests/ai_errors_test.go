package tui_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/delivery/tui/match"
	chess "github.com/gabrzb/chess-go/internal/domain"
	"github.com/gabrzb/chess-go/internal/domain/game"
)

func TestAIErrorPreservesCause(t *testing.T) {
	cause := errors.New("search failed")
	m, search := aiModel(t, game.NewGame(), chooserFunc(func(context.Context, chess.Position) (chess.Move, error) { return chess.Move{}, cause }), chess.White)
	if search == nil {
		t.Fatal("AI search did not start")
	}
	next, _ := m.Update(search())
	got := next.(match.Model)
	if !errors.Is(got.State().AIError, cause) {
		t.Fatal("AI failure lost its cause")
	}
	got.ClearMessage()
	if got.State().AIError != nil || got.State().Message != "" {
		t.Fatal("retry did not clear the failure")
	}
}
