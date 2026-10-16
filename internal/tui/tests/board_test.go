package tui_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/tui"
)

func TestBoardRendering(t *testing.T) {
	m := tui.NewModel(gameWith(t))
	board := ansi.Strip(boardText(t, m))
	rawBoard := boardText(t, m)
	if strings.Count(rawBoard, "48;2;") != 1 || strings.Contains(rawBoard, ansi.NewStyle(ansi.AttrReverse).String()) {
		t.Fatal("Cursor background must affect only its cell, without reverse video")
	}
	if len(boardCells(t, m)) != chess.SquareCount || !strings.Contains(board, "  a b c d e f g h") || !strings.Contains(m.View().Content, "Turn: White") {
		t.Fatal("Missing initial board, coordinates or turn")
	}
	if got := ansi.Strip(boardCells(t, m)[chess.E1]); got != "\u2654" {
		t.Errorf("White king glyph: %q", got)
	}
	if got := ansi.Strip(boardCells(t, m)[chess.E8]); got != "\u265a" {
		t.Errorf("Black king glyph: %q", got)
	}
	m = moveCursorTo(t, m, chess.H8)
	if ansi.Strip(boardText(t, m)) != board {
		t.Error("Moving the cursor altered the rendered position")
	}
	assertHighlights(t, m, chess.NoSquare)
}

func TestCursorKeysAndBoundaries(t *testing.T) {
	m := tui.NewModel(gameWith(t))
	for _, tc := range []struct {
		code             rune
		from, want, edge chess.Square
	}{
		{tea.KeyUp, chess.D4, chess.D5, chess.A8}, {'w', chess.D4, chess.D5, chess.A8},
		{tea.KeyDown, chess.D4, chess.D3, chess.A1}, {'s', chess.D4, chess.D3, chess.A1},
		{tea.KeyLeft, chess.D4, chess.C4, chess.A1}, {'a', chess.D4, chess.C4, chess.A1},
		{tea.KeyRight, chess.D4, chess.E4, chess.H8}, {'d', chess.D4, chess.E4, chess.H8},
	} {
		m = moveCursorTo(t, m, tc.from)
		m = updateModel(t, m, tea.KeyPressMsg{Code: tc.code})
		if got := visibleCursor(t, m); got != tc.want {
			t.Errorf("Key %v: got %v, want %v", tc.code, got, tc.want)
		}
		m = moveCursorTo(t, m, tc.edge)
		m = updateModel(t, m, tea.KeyPressMsg{Code: tc.code})
		if got := visibleCursor(t, m); got != tc.edge {
			t.Errorf("Key %v moved beyond board", tc.code)
		}
	}
	before := m
	m = updateModel(t, m, tea.KeyReleaseMsg{Code: tea.KeyLeft})
	m = updateModel(t, m, tea.KeyReleaseMsg{Code: 'q'})
	m = updateModel(t, m, tea.WindowSizeMsg{})
	assertSameView(t, m, before)
}
