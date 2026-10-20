package render

import (
	"image/color"
	"reflect"
	"testing"

	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
)

func sameColor(got color.Color, hex string) bool {
	if got == nil {
		return false
	}
	r, g, b, a := got.RGBA()
	wr, wg, wb, wa := lipgloss.Color(hex).RGBA()
	return r == wr && g == wg && b == wb && a == wa
}

func TestBoardCellsKeepTheirColorsAndProportions(t *testing.T) {
	board := chess.NewInitialPosition().Board
	for _, scale := range []BoardScale{CompactBoard, LargeBoard} {
		for _, flipped := range []bool{false, true} {
			width, height := scale.CellDimensions()
			if width != 2*height {
				t.Fatal("terminal cells lost their square proportions")
			}
			state := BoardState{Cursor: chess.A2, Selected: chess.B2, LegalMoves: []chess.Move{{To: chess.C3}}, Scale: scale, Flipped: flipped}
			lines := uv.NewStyledString(Board(board, state)).Lines(ansi.GraphemeWidth)
			if len(lines) != chess.BoardSize*height+1 {
				t.Fatal("incorrect board height")
			}
			for y, line := range lines[:len(lines)-1] {
				if len(line) != BoardGutter+chess.BoardSize*width {
					t.Fatal("uneven board row width")
				}
				for x := BoardGutter; x < len(line); x++ {
					file, rank := (x-BoardGutter)/width, chess.BoardSize-1-y/height
					if flipped {
						file, rank = chess.BoardSize-1-file, chess.BoardSize-1-rank
					}
					background := LightSquare
					if (file+rank)%2 == 0 {
						background = DarkSquare
					}
					cell := line[x]
					if !sameColor(cell.Style.Bg, background) {
						t.Fatalf("cell (%d,%d) lost its tile background", x, y)
					}
					if cell.Content == " " || cell.Content == MenuCursor || cell.Content == "●" || cell.Content == "·" {
						continue
					}
					square, _ := chess.A1.Offset(file, rank)
					piece, _ := board.PieceAt(square)
					foreground := WhitePiece
					if piece.Color == chess.Black {
						foreground = BlackPiece
					}
					if !sameColor(cell.Style.Fg, foreground) {
						t.Fatalf("piece %v lost its color", square)
					}
				}
			}
		}
	}
}

func TestMovingCursorNeverChangesTileBackgrounds(t *testing.T) {
	board := chess.NewInitialPosition().Board
	before := uv.NewStyledString(Board(board, BoardState{Cursor: chess.A2})).Lines(ansi.GraphemeWidth)
	after := uv.NewStyledString(Board(board, BoardState{Cursor: chess.H8})).Lines(ansi.GraphemeWidth)
	for y, row := range before {
		for x, cell := range row {
			if !reflect.DeepEqual(cell.Style.Bg, after[y][x].Style.Bg) {
				t.Fatalf("focus changed background at %d,%d", x, y)
			}
		}
	}
}
