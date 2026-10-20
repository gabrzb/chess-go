package render

import (
	"reflect"
	"strings"
	"testing"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestTurnBlockMatchesSideAndContext(t *testing.T) {
	for _, tc := range []struct {
		state                    TurnState
		title, label, background string
	}{
		{TurnState{Color: chess.White}, "VEZ DAS", "BRANCAS", WhitePiece},
		{TurnState{Color: chess.Black}, "VEZ DAS", "PRETAS", BlackPiece},
		{TurnState{Color: chess.White, Reviewing: true}, "REVISÃO", "BRANCAS", WhitePiece},
		{TurnState{Color: chess.Black, Finished: true}, "PARTIDA", "ENCERRADA", Background},
	} {
		block := Turn(tc.state, CompactBoard)
		text := ansi.Strip(block)
		if !strings.Contains(text, tc.title) || !strings.Contains(text, tc.label) {
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
		if colored != compactTurnHeight*compactTurnHeight*terminalCellAspect {
			t.Fatalf("turn block lost its rectangular fill: %d cells", colored)
		}
	}
}

func TestMatchFitsStandardWindows(t *testing.T) {
	for _, size := range []struct{ width, height int }{{80, 24}, {120, 40}, {48, 24}} {
		for _, color := range []chess.Color{chess.White, chess.Black} {
			panel := MatchPanel{Width: size.width, Height: size.height, Board: chess.NewInitialPosition().Board,
				BoardState: BoardState{Cursor: chess.A2}, Turn: TurnState{Color: color},
				Status: "Vez das brancas\nJogador 1", CompactStatus: "Vez das brancas", LastMove: "—",
				Mode: "LOCAL · JOGADOR VS JOGADOR", CompactMode: "LOCAL", Footer: "Setas/WASD mover · Enter selecionar\nEsc menu · F inverter · H histórico · ? ajuda"}
			content := Match(panel)
			if !Fits(content, size.width, size.height) {
				t.Fatalf("match overflows %dx%d", size.width, size.height)
			}
			if strings.Count(ansi.Strip(content), MenuCursor) != 1 {
				t.Fatal("cursor duplicated by turn indicator")
			}
		}
	}
}

func TestWindowDoesNotSpreadCursorBackground(t *testing.T) {
	panel := MatchPanel{Width: DefaultWidth, Height: DefaultHeight, Board: chess.NewInitialPosition().Board,
		BoardState: BoardState{Cursor: chess.A2}, Turn: TurnState{Color: chess.White},
		Status: "Vez das brancas", CompactStatus: "Vez das brancas", LastMove: "—"}
	before := uv.NewStyledString(Window(Match(panel), panel.Width, panel.Height).Content).Lines(ansi.GraphemeWidth)
	panel.BoardState.Cursor = chess.H8
	after := uv.NewStyledString(Window(Match(panel), panel.Width, panel.Height).Content).Lines(ansi.GraphemeWidth)
	for y, line := range before {
		for x, cell := range line {
			if !reflect.DeepEqual(cell.Style.Bg, after[y][x].Style.Bg) {
				t.Fatalf("focus leaked into window background at %d,%d", x, y)
			}
		}
	}
}
