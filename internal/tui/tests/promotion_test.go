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

type promotionCase struct {
	name  string
	setup []chess.Move
	move  chess.Move
}

var (
	whitePromotionSetup = []chess.Move{
		{From: chess.A2, To: chess.A4},
		{From: chess.B8, To: chess.C6},
		{From: chess.A4, To: chess.A5},
		{From: chess.H7, To: chess.H5},
		{From: chess.A5, To: chess.A6},
		{From: chess.H5, To: chess.H4},
		{From: chess.A6, To: chess.B7},
		{From: chess.H4, To: chess.H3},
	}
	blackPromotionSetup = append(slices.Clone(whitePromotionSetup),
		chess.Move{From: chess.B7, To: chess.B8, Promotion: chess.Knight},
		chess.Move{From: chess.H3, To: chess.G2},
		chess.Move{From: chess.G1, To: chess.F3},
	)
	promotionCases = []promotionCase{
		{"white advance", whitePromotionSetup, chess.Move{From: chess.B7, To: chess.B8}},
		{"white capture", whitePromotionSetup, chess.Move{From: chess.B7, To: chess.A8}},
		{"black advance", blackPromotionSetup, chess.Move{From: chess.G2, To: chess.G1}},
		{"black capture", blackPromotionSetup, chess.Move{From: chess.G2, To: chess.H1}},
	}
)

func pendingPromotion(t *testing.T, tc promotionCase) (tea.Model, *game.Game, game.Snapshot) {
	t.Helper()
	g := game.NewGame()
	for _, move := range tc.setup {
		if _, err := g.Play(move); err != nil {
			t.Fatalf("Failed to play setup move %v: %v", move, err)
		}
	}
	before, err := g.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	m := moveCursorTo(t, tui.NewModel(g), tc.move.From)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	m = moveCursorTo(t, m, tc.move.To)
	m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if !strings.Contains(m.View().Content, "Promotion:") {
		t.Fatalf("Expected promotion choices in view, got %q", m.View().Content)
	}
	assertGameUnchanged(t, g, before, len(tc.setup)+1)
	return m, g, before
}

func TestPromotionChoices(t *testing.T) {
	choices := []struct {
		key          string
		pieceType    chess.PieceType
		white, black string
	}{
		{"q", chess.Queen, "♕", "♛"},
		{"r", chess.Rook, "♖", "♜"},
		{"b", chess.Bishop, "♗", "♝"},
		{"n", chess.Knight, "♘", "♞"},
	}
	for _, tc := range promotionCases {
		for _, choice := range choices {
			t.Run(tc.name+"/"+choice.key, func(t *testing.T) {
				m, g, before := pendingPromotion(t, tc)
				m = updateModel(t, m, tea.KeyPressMsg{Code: rune(choice.key[0]), Text: choice.key})
				after, err := g.Snapshot()
				if err != nil {
					t.Fatal(err)
				}
				move := tc.move
				move.Promotion = choice.pieceType
				if after.LastMove == nil || *after.LastMove != move {
					t.Fatalf("Expected completed promotion %v, got %+v", move, after)
				}
				want := chess.Piece{Type: choice.pieceType, Color: before.SideToMove}
				if piece, _ := after.Board.PieceAt(move.To); piece != want {
					t.Errorf("Expected promoted piece %v at %v, got %v", want, move.To, piece)
				}
				if piece, _ := after.Board.PieceAt(move.From); piece.Type != chess.NoPiece {
					t.Errorf("Expected empty promotion source, got %v", piece)
				}
				if after.SideToMove == before.SideToMove || g.HistoryLen() != len(tc.setup)+2 {
					t.Error("Expected completed promotion to change turn and add exactly one history position")
				}
				symbol := choice.white
				if before.SideToMove == chess.Black {
					symbol = choice.black
				}
				if got := ansi.Strip(boardCells(t, m)[move.To]); got != symbol {
					t.Errorf("Expected promoted symbol %q, got %q", symbol, got)
				}
				expected := moveCursorTo(t, tui.NewModel(g), move.To)
				if !reflect.DeepEqual(m.View(), expected.View()) {
					t.Error("Expected completed promotion to clear choices, selection and messages")
				}
			})
		}
	}
}

func TestPromotionIgnoresOtherEvents(t *testing.T) {
	cases := []struct {
		name string
		msg  tea.Msg
	}{
		{"arrow", tea.KeyPressMsg{Code: tea.KeyRight}},
		{"WASD", tea.KeyPressMsg{Code: 'd', Text: "d"}},
		{"Enter", tea.KeyPressMsg{Code: tea.KeyEnter}},
		{"unknown choice", tea.KeyPressMsg{Code: 'x', Text: "x"}},
		{"history shortcut", tea.KeyPressMsg{Code: 'h', Text: "h"}},
		{"queen key release", tea.KeyReleaseMsg{Code: 'q', Text: "q"}},
		{"Escape release", tea.KeyReleaseMsg{Code: tea.KeyEscape}},
		{"window event", tea.WindowSizeMsg{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, g, before := pendingPromotion(t, promotionCases[0])
			view := m.View()
			length := g.HistoryLen()
			m = updateModel(t, m, tc.msg)
			if !reflect.DeepEqual(m.View(), view) {
				t.Error("Expected ignored event to preserve pending promotion and cursor")
			}
			assertGameUnchanged(t, g, before, length)
			m = updateModel(t, m, tea.KeyPressMsg{Code: 'n', Text: "n"})
			after, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			move := promotionCases[0].move
			move.Promotion = chess.Knight
			if after.LastMove == nil || *after.LastMove != move || g.HistoryLen() != length+1 {
				t.Error("Expected original promotion to remain executable after ignored event")
			}
		})
	}
}

func TestPromotionCancellation(t *testing.T) {
	for _, tc := range []promotionCase{promotionCases[0], promotionCases[2]} {
		t.Run(tc.name, func(t *testing.T) {
			m, g, before := pendingPromotion(t, tc)
			length := g.HistoryLen()
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEscape})
			expected := moveCursorTo(t, tui.NewModel(g), tc.move.To)
			if !reflect.DeepEqual(m.View(), expected.View()) {
				t.Error("Expected Escape to clear choices and selection while keeping cursor")
			}
			assertGameUnchanged(t, g, before, length)
			m = moveCursorTo(t, m, tc.move.From)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			m = moveCursorTo(t, m, tc.move.To)
			m = updateModel(t, m, tea.KeyPressMsg{Code: tea.KeyEnter})
			if !strings.Contains(m.View().Content, "Promotion:") {
				t.Fatal("Expected promotion to be available again after cancellation")
			}
			m = updateModel(t, m, tea.KeyPressMsg{Code: 'r', Text: "r"})
			after, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			move := tc.move
			move.Promotion = chess.Rook
			if after.LastMove == nil || *after.LastMove != move || g.HistoryLen() != length+1 {
				t.Error("Expected retried promotion to execute exactly once")
			}
		})
	}
}

func TestQuitDuringPromotion(t *testing.T) {
	for _, pending := range []bool{false, true} {
		name := "q outside promotion"
		if pending {
			name = "Ctrl+C during promotion"
		}
		t.Run(name, func(t *testing.T) {
			m, g := modelAt(t, chess.A2)
			key := tea.KeyPressMsg{Code: 'q', Text: "q"}
			if pending {
				m, g, _ = pendingPromotion(t, promotionCases[0])
				key = tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
			}
			before, err := g.Snapshot()
			if err != nil {
				t.Fatal(err)
			}
			length := g.HistoryLen()
			updated, cmd := m.Update(key)
			if cmd == nil {
				t.Fatal("Expected quit command")
			}
			if _, ok := cmd().(tea.QuitMsg); !ok {
				t.Fatal("Expected command to produce quit message")
			}
			if !reflect.DeepEqual(updated.View(), m.View()) {
				t.Error("Expected quit request to preserve displayed state")
			}
			assertGameUnchanged(t, g, before, length)
		})
	}
}
