package tui_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/gabrzb/chess-go/internal/chess"
	"github.com/gabrzb/chess-go/internal/chess/game"
	"github.com/gabrzb/chess-go/internal/tui"
)

func TestViewInitialBoard(t *testing.T) {
	m := tui.NewModel(game.NewGame())
	want := `8 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜
7 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟
6 . . . . . . . .
5 . . . . . . . .
4 . . . . . . . .
3 . . . . . . . .
2 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙
1 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖
  a b c d e f g h`

	if got := ansi.Strip(boardText(t, m)); got != want {
		t.Errorf("Unexpected initial board rendering:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestViewBoardCursor(t *testing.T) {
	plain := ansi.Strip(boardText(t, tui.NewModel(game.NewGame())))
	for _, cursor := range []chess.Square{chess.A2, chess.D4, chess.H8} {
		t.Run(cursor.String(), func(t *testing.T) {
			m, _ := modelAt(t, cursor)
			got := boardText(t, m)
			if ansi.Strip(got) != plain {
				t.Fatal("Expected cursor styling to preserve board text and coordinates")
			}
			lines := strings.Split(got, "\n")
			for rowIndex, line := range lines[:chess.BoardSize] {
				cells := strings.Fields(line)
				if len(cells) != chess.BoardSize+1 {
					t.Fatalf("Expected rank label and %d cells, got %q", chess.BoardSize, line)
				}
				for col, cell := range cells[1:] {
					styled := ansi.Strip(cell) != cell
					wantStyled := rowIndex == chess.A8.Row()-cursor.Row() && col == cursor.Column()
					if styled != wantStyled {
						t.Errorf("At rank %s, column %d: expected styling %v, got %q", cells[0], col, wantStyled, cell)
					}
				}
			}
		})
	}
}
