package match

import (
	"errors"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess/game"
)

func TestAIErrorPreservesCause(t *testing.T) {
	m := New(game.NewGame(), 1)
	m.thinking = true
	cause := errors.New("search failed")
	next, _ := m.receiveAIResult(aiResultMessage{sessionID: 1, position: m.game.Position(), err: cause})
	got := next.(Model)
	if !errors.Is(got.State().AIError, cause) {
		t.Fatal("AI failure lost its cause")
	}
	got.ClearMessage()
	if got.State().AIError != nil || got.State().Message != "" {
		t.Fatal("retry did not clear the failure")
	}
}
