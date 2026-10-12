package tui

import (
	"fmt"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestUpdateCursor(t *testing.T) {
	cases := []struct {
		name string
		keys []tea.KeyPressMsg
		want chess.Square
		edge chess.Square
	}{
		{"up", []tea.KeyPressMsg{{Code: tea.KeyUp}, {Code: 'w', Text: "w"}}, chess.D5, chess.A8},
		{"down", []tea.KeyPressMsg{{Code: tea.KeyDown}, {Code: 's', Text: "s"}}, chess.D3, chess.A1},
		{"left", []tea.KeyPressMsg{{Code: tea.KeyLeft}, {Code: 'a', Text: "a"}}, chess.C4, chess.A1},
		{"right", []tea.KeyPressMsg{{Code: tea.KeyRight}, {Code: 'd', Text: "d"}}, chess.E4, chess.H8},
	}
	for _, tc := range cases {
		for _, key := range tc.keys {
			for _, state := range []struct {
				name string
				from chess.Square
				want chess.Square
			}{
				{"inside board", chess.D4, tc.want},
				{"at edge", tc.edge, tc.edge},
			} {
				t.Run(fmt.Sprintf("%s/%s/%s", tc.name, key.String(), state.name), func(t *testing.T) {
					m := newModel()
					m.cursor = state.from
					before, err := m.game.Snapshot()
					if err != nil {
						t.Fatalf("Failed to get snapshot before key event: %v", err)
					}

					updated, cmd := m.Update(key)
					got, ok := updated.(model)
					if !ok {
						t.Fatalf("Expected TUI model, got %T", updated)
					}
					if cmd != nil {
						t.Error("Expected cursor movement without a command")
					}
					if got.cursor != state.want {
						t.Errorf("Expected cursor at %v, got %v", state.want, got.cursor)
					}
					after, err := got.game.Snapshot()
					if err != nil {
						t.Fatalf("Failed to get snapshot after key event: %v", err)
					}
					if after != before {
						t.Error("Expected cursor movement to preserve game state")
					}
				})
			}
		}
	}
}

func TestUpdateIgnoresNonMovementEvents(t *testing.T) {
	cases := []struct {
		name string
		msg  tea.Msg
	}{
		{"arrow release", tea.KeyReleaseMsg{Code: tea.KeyUp}},
		{"WASD release", tea.KeyReleaseMsg{Code: 'w', Text: "w"}},
		{"quit key release", tea.KeyReleaseMsg{Code: 'q', Text: "q"}},
		{"unrelated key", tea.KeyPressMsg{Code: 'x', Text: "x"}},
		{"window size event", tea.WindowSizeMsg{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := newModel()
			m.cursor = chess.D4
			updated, cmd := m.Update(tc.msg)
			got, ok := updated.(model)
			if !ok {
				t.Fatalf("Expected TUI model, got %T", updated)
			}
			if got.cursor != m.cursor {
				t.Errorf("Expected cursor to stay at %v, got %v", m.cursor, got.cursor)
			}
			if cmd != nil {
				t.Error("Expected ignored event without a command")
			}
		})
	}
}
