package tui_test

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui"
)

func boardText(t *testing.T, m tea.Model) string {
	t.Helper()
	sections := strings.SplitN(m.View().Content, "\n\n", 3)
	if len(sections) != 3 {
		t.Fatalf("Expected title, board and controls in view, got %q", m.View().Content)
	}
	return sections[1]
}

func visibleCursor(t *testing.T, m tea.Model) chess.Square {
	t.Helper()
	lines := strings.Split(boardText(t, m), "\n")
	if len(lines) != chess.BoardSize+1 {
		t.Fatalf("Expected %d board ranks and coordinate footer, got %q", chess.BoardSize, lines)
	}
	cursor := chess.NoSquare
	for _, line := range lines[:chess.BoardSize] {
		cells := strings.Fields(line)
		if len(cells) != chess.BoardSize+1 {
			t.Fatalf("Expected rank label and %d cells, got %q", chess.BoardSize, line)
		}
		for col, cell := range cells[1:] {
			if ansi.Strip(cell) == cell {
				continue
			}
			if cursor.Valid() {
				t.Fatal("Expected exactly one highlighted cursor cell")
			}
			square, err := chess.ParseSquare(fmt.Sprintf("%c%s", 'a'+rune(col), cells[0]))
			if err != nil {
				t.Fatalf("Invalid highlighted cell coordinate: %v", err)
			}
			cursor = square
		}
	}
	if !cursor.Valid() {
		t.Fatal("Expected a highlighted cursor cell")
	}
	return cursor
}

func updateModel(t *testing.T, m tea.Model, msg tea.Msg) tea.Model {
	t.Helper()
	updated, cmd := m.Update(msg)
	if cmd != nil {
		t.Fatalf("Expected event without a command, got %T", msg)
	}
	if updated == nil {
		t.Fatal("Expected updated TUI model")
	}
	return updated
}

func modelAt(t *testing.T, target chess.Square) (tea.Model, *game.Game) {
	t.Helper()
	g := game.NewGame()
	m := tui.NewModel(g)
	from := visibleCursor(t, m)
	for _, axis := range []struct {
		delta              int
		positive, negative rune
	}{
		{target.Column() - from.Column(), tea.KeyRight, tea.KeyLeft},
		{target.Row() - from.Row(), tea.KeyUp, tea.KeyDown},
	} {
		key, steps := axis.positive, axis.delta
		if steps < 0 {
			key, steps = axis.negative, -steps
		}
		for step := 0; step < steps; step++ {
			m = updateModel(t, m, tea.KeyPressMsg{Code: key})
		}
	}
	if got := visibleCursor(t, m); got != target {
		t.Fatalf("Failed to position cursor: expected %v, got %v", target, got)
	}
	return m, g
}
