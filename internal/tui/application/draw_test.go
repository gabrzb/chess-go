package application

import (
	"context"
	"strings"
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func TestClaimDrawFromPauseMenu(t *testing.T) {
	model, err := New(context.Background(), Options{StartGame: true, NoAnimations: true})
	if err != nil {
		t.Fatal(err)
	}
	m := model.(application)
	for range 2 {
		for _, move := range []chess.Move{{From: chess.G1, To: chess.F3}, {From: chess.G8, To: chess.F6}, {From: chess.F3, To: chess.G1}, {From: chess.F6, To: chess.G8}} {
			if _, err := m.match.Game().Play(move); err != nil {
				t.Fatal(err)
			}
		}
	}
	m.show(pauseScreen)
	options := m.menuOptions()
	m.choice = len(options) - 1
	if options[m.choice] != "Reivindicar empate" {
		t.Fatal("claim absent from menu")
	}
	m.activate()
	s, err := m.match.Game().Snapshot()
	if err != nil || m.screen != resultScreen || s.Outcome.Reason != game.Repetition {
		t.Fatal("claim did not end game", err)
	}
	if !strings.Contains(m.View().Content, "tripla repetição") {
		t.Fatal("draw reason absent from result")
	}
}
