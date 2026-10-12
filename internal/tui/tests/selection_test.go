package tui_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui"
)

func assertSelectionHighlights(t *testing.T, m tea.Model, origin chess.Square, destinations []chess.Square, cursor chess.Square) {
	t.Helper()
	if got := visibleCursor(t, m); got != cursor {
		t.Errorf("Expected cursor at %v, got %v", cursor, got)
	}
	cells := boardCells(t, m)
	for square, cell := range cells {
		plain := ansi.Strip(cell)
		wantStyled := square == cursor || square == origin || slices.Contains(destinations, square)
		if styled := cell != plain; styled != wantStyled {
			t.Errorf("At %v: expected highlighting %v, got %q", square, wantStyled, cell)
		}
		if square == cursor && cell != ansi.NewStyle(ansi.AttrReverse).Styled(plain) {
			t.Errorf("Expected cursor style to take priority at %v, got %q", square, cell)
		}
	}
	if origin != cursor {
		originStyle := strings.ReplaceAll(cells[origin], ansi.Strip(cells[origin]), "")
		for _, destination := range destinations {
			if destination == cursor {
				continue
			}
			destinationStyle := strings.ReplaceAll(cells[destination], ansi.Strip(cells[destination]), "")
			if destinationStyle == originStyle {
				t.Errorf("Expected distinct styles for selected origin and destination %v", destination)
			}
		}
	}
}

func TestSelectPieceHighlights(t *testing.T) {
	cases := []struct {
		name         string
		blackTurn    bool
		origin       chess.Square
		destinations []chess.Square
	}{
		{"white pawn", false, chess.A2, []chess.Square{chess.A3, chess.A4}},
		{"black pawn", true, chess.A7, []chess.Square{chess.A6, chess.A5}},
		{"knight", false, chess.B1, []chess.Square{chess.A3, chess.C3}},
		{"blocked rook", false, chess.A1, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := game.NewGame()
			if tc.blackTurn {
				if _, err := g.Play(chess.Move{From: chess.E2, To: chess.E4}); err != nil {
					t.Fatalf("Failed to set up Black turn: %v", err)
				}
			}
			m := moveCursorTo(t, tui.NewModel(g), tc.origin)
			before, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			length := g.HistoryLen()
			plain := ansi.Strip(boardText(t, m))
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			assertSelectionHighlights(t, m, tc.origin, tc.destinations, tc.origin)
			m = moveCursorTo(t, m, chess.D4)
			assertSelectionHighlights(t, m, tc.origin, tc.destinations, chess.D4)
			if len(tc.destinations) > 0 {
				m = moveCursorTo(t, m, tc.destinations[0])
				assertSelectionHighlights(t, m, tc.origin, tc.destinations, tc.destinations[0])
			}
			if got := ansi.Strip(boardText(t, m)); got != plain {
				t.Error("Expected selection and highlights to preserve board text")
			}
			assertGameUnchanged(t, g, before, length)
		})
	}
}

func TestInvalidSelectionPreservesState(t *testing.T) {
	for _, side := range []struct {
		name       string
		setup      []chess.Move
		own, enemy chess.Square
	}{
		{name: "white", own: chess.A2, enemy: chess.A7},
		{name: "black", setup: []chess.Move{{From: chess.E2, To: chess.E4}}, own: chess.A7, enemy: chess.A2},
	} {
		for _, invalid := range []struct {
			name   string
			square chess.Square
			err    error
		}{
			{"empty square", chess.D4, chess.ErrEmptySource},
			{"opponent piece", side.enemy, chess.ErrIllegalMove},
		} {
			t.Run(side.name+"/"+invalid.name, func(t *testing.T) {
				g := game.NewGame()
				for _, move := range side.setup {
					if _, err := g.Play(move); err != nil {
						t.Fatal(err)
					}
				}
				m := moveCursorTo(t, tui.NewModel(g), side.own)
				m = moveCursorTo(t, m, invalid.square)
				beforeBoard := boardText(t, m)
				before, err := g.Snapshot()
				if err != nil {
					t.Fatal(err)
				}
				length := g.HistoryLen()
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				if got := boardText(t, m); got != beforeBoard {
					t.Error("Expected failed selection to preserve existing highlights and cursor")
				}
				if !strings.Contains(m.View().Content, invalid.err.Error()) {
					t.Errorf("Expected error message %q in view, got %q", invalid.err, m.View().Content)
				}
				assertGameUnchanged(t, g, before, length)
			})
		}
	}
}

func TestSelectionClearsError(t *testing.T) {
	for _, cancel := range []bool{true, false} {
		name := "Escape cancels"
		if !cancel {
			name = "successful selection clears message"
		}
		t.Run(name, func(t *testing.T) {
			m, g := modelAt(t, chess.A2)
			wantErr := chess.ErrEmptySource
			if cancel {
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				wantErr = chess.ErrIllegalMove
			}
			m = moveCursorTo(t, m, chess.D4)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			if !strings.Contains(m.View().Content, wantErr.Error()) {
				t.Fatal("Expected setup error message")
			}
			before, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			length := g.HistoryLen()
			var expected tea.Model
			if cancel {
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEscape})
				expected, _ = modelAt(t, chess.D4)
			} else {
				m = moveCursorTo(t, m, chess.B2)
				m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
				expected, _ = modelAt(t, chess.B2)
				expected = updateModel(t, expected, tea.KeyPressMsg{Code: tea.KeyEnter})
			}
			if !reflect.DeepEqual(m.View(), expected.View()) {
				t.Error("Expected cleared error and matching selection state in view")
			}
			assertGameUnchanged(t, g, before, length)
		})
	}
}

func TestSelectionIgnoresKeyReleases(t *testing.T) {
	for _, key := range []tea.KeyReleaseMsg{{Code: tea.KeyEnter}, {Code: tea.KeyEscape}} {
		t.Run(key.String(), func(t *testing.T) {
			m, _ := modelAt(t, chess.A2)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			m = moveCursorTo(t, m, chess.D4)
			before := m.View()
			m = updateModel(t, m, key)
			if !reflect.DeepEqual(m.View(), before) {
				t.Error("Expected key release to preserve selection, cursor and message")
			}
		})
	}
}
