package tui_test

import (
	"reflect"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
)

func boardText(t *testing.T, m tea.Model) string {
	t.Helper()
	parts := strings.SplitN(m.View().Content, "\n\n", 3)
	if len(parts) != 3 {
		t.Fatalf("Missing board in %q", m.View().Content)
	}
	return parts[1]
}

func boardCells(t *testing.T, m tea.Model) map[chess.Square]string {
	t.Helper()
	lines := strings.Split(boardText(t, m), "\n")
	if len(lines) != chess.BoardSize+1 {
		t.Fatalf("Incorrect board dimensions: %q", lines)
	}
	cells := make(map[chess.Square]string)
	for row, line := range lines[:chess.BoardSize] {
		if ansi.StringWidth(line) != 2+chess.BoardSize*2 {
			t.Fatalf("Incorrect rank: %q", line)
		}
		for col := range chess.BoardSize {
			cell := strings.TrimSpace(ansi.Cut(line, 2+col*2, 2+(col+1)*2))
			square, _ := chess.A8.Offset(col, -row)
			cells[square] = cell
		}
	}
	return cells
}

func visibleCursor(t *testing.T, m tea.Model) chess.Square {
	t.Helper()
	cursor := chess.NoSquare
	for square, cell := range boardCells(t, m) {
		if strings.Contains(ansi.Strip(cell), "▸") {
			if cursor.Valid() {
				t.Fatal("Multiple cursors")
			}
			cursor = square
		}
	}
	if !cursor.Valid() {
		t.Fatal("Missing cursor")
	}
	return cursor
}

func updateModel(t *testing.T, m tea.Model, msg tea.Msg) tea.Model {
	t.Helper()
	next, cmd := m.Update(msg)
	if next == nil || cmd != nil {
		t.Fatalf("Expected state update without command for %T", msg)
	}
	return next
}

func moveCursorTo(t *testing.T, m tea.Model, to chess.Square) tea.Model {
	t.Helper()
	from := visibleCursor(t, m)
	for _, axis := range []struct {
		delta              int
		positive, negative rune
	}{
		{to.Column() - from.Column(), tea.KeyRight, tea.KeyLeft}, {to.Row() - from.Row(), tea.KeyUp, tea.KeyDown},
	} {
		code, steps := axis.positive, axis.delta
		if steps < 0 {
			code, steps = axis.negative, -steps
		}
		for range steps {
			m = updateModel(t, m, tea.KeyPressMsg{Code: code})
		}
	}
	if got := visibleCursor(t, m); got != to {
		t.Fatalf("Cursor: got %v, want %v", got, to)
	}
	return m
}

func submit(t *testing.T, m tea.Model, move chess.Move) (tea.Model, tea.Cmd) {
	t.Helper()
	m = moveCursorTo(t, m, move.From)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	m = moveCursorTo(t, m, move.To)
	return m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
}

func gameWith(t *testing.T, moves ...chess.Move) *game.Game {
	t.Helper()
	g := game.NewGame()
	for _, move := range moves {
		if _, err := g.Play(move); err != nil {
			t.Fatalf("Setup %v: %v", move, err)
		}
	}
	return g
}

func assertUnchanged(t *testing.T, g *game.Game, position chess.Position, length int) {
	t.Helper()
	if g.Position() != position || g.HistoryLen() != length {
		t.Fatal("Interface changed the live game or history")
	}
}

func assertHighlights(t *testing.T, m tea.Model, selected chess.Square, destinations ...chess.Square) {
	t.Helper()
	cursor := visibleCursor(t, m)
	for square, cell := range boardCells(t, m) {
		expected := square == cursor || square == selected
		for _, to := range destinations {
			expected = expected || square == to
		}
		highlighted := strings.ContainsAny(ansi.Strip(cell), "▸●·")
		if highlighted != expected {
			t.Errorf("Incorrect highlight at %v", square)
		}
	}
}

func assertSameView(t *testing.T, a, b tea.Model) {
	t.Helper()
	if !reflect.DeepEqual(a.View(), b.View()) {
		t.Error("Unexpected interface state")
	}
}
