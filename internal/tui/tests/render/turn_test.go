package render_test

import (
	"github.com/gabrzb/chess-go/internal/tui/render"
	"reflect"
	"regexp"
	"strings"
	"testing"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestTurnBlockMatchesSideAndContext(t *testing.T) {
	for _, tc := range []struct {
		state                    render.TurnState
		title, label, background string
	}{
		{render.TurnState{Color: chess.White}, "VEZ DAS", "W", render.WhitePiece},
		{render.TurnState{Color: chess.Black}, "VEZ DAS", "B", render.BlackPiece},
		{render.TurnState{Color: chess.White, Reviewing: true}, "REVISÃO", "W", render.WhitePiece},
		{render.TurnState{Color: chess.Black, Finished: true}, "PARTIDA", "ENCERRADA", render.Background},
	} {
		block := render.Turn(tc.state, render.CompactBoard)
		text := ansi.Strip(block)
		if !strings.Contains(text, tc.title) || !regexp.MustCompile(`│\s+`+regexp.QuoteMeta(tc.label)+`\s+│`).MatchString(text) {
			t.Fatal("incorrect turn context", text)
		}
		colored := 0
		for _, line := range uv.NewStyledString(block).Lines(ansi.GraphemeWidth) {
			for _, cell := range line {
				if sameColor(cell.Style.Bg, tc.background) {
					colored++
				}
			}
		}
		if colored != 72 {
			t.Fatalf("turn block lost its rectangular fill: %d cells", colored)
		}
	}
}

func TestMatchFitsStandardWindows(t *testing.T) {
	for _, size := range []struct{ width, height int }{{80, 24}, {120, 40}, {48, 24}} {
		for _, color := range []chess.Color{chess.White, chess.Black} {
			panel := render.MatchPanel{Width: size.width, Height: size.height, Board: chess.NewInitialPosition().Board,
				BoardState: render.BoardState{Cursor: chess.A2}, Turn: render.TurnState{Color: color},
				Status: "Vez das brancas\nJogador 1", CompactStatus: "Vez das brancas", LastMove: "—",
				Mode: "LOCAL · JOGADOR VS JOGADOR", CompactMode: "LOCAL", Footer: "Setas/WASD mover · Enter selecionar\nEsc menu · F inverter · H histórico · ? ajuda"}
			content := render.Match(panel)
			if !render.Fits(content, size.width, size.height) {
				t.Fatalf("match overflows %dx%d", size.width, size.height)
			}
			if strings.Count(ansi.Strip(content), render.MenuCursor) != 1 {
				t.Fatal("cursor duplicated by turn indicator")
			}
		}
	}
}

func TestWindowDoesNotSpreadCursorBackground(t *testing.T) {
	panel := render.MatchPanel{Width: render.DefaultWidth, Height: render.DefaultHeight, Board: chess.NewInitialPosition().Board,
		BoardState: render.BoardState{Cursor: chess.A2}, Turn: render.TurnState{Color: chess.White},
		Status: "Vez das brancas", CompactStatus: "Vez das brancas", LastMove: "—"}
	before := uv.NewStyledString(render.Window(render.Match(panel), panel.Width, panel.Height).Content).Lines(ansi.GraphemeWidth)
	panel.BoardState.Cursor = chess.H8
	after := uv.NewStyledString(render.Window(render.Match(panel), panel.Width, panel.Height).Content).Lines(ansi.GraphemeWidth)
	for y, line := range before {
		for x, cell := range line {
			if !reflect.DeepEqual(cell.Style.Bg, after[y][x].Style.Bg) {
				t.Fatalf("focus leaked into window background at %d,%d", x, y)
			}
		}
	}
}
