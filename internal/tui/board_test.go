package tui

import (
	"testing"

	"github.com/gabrzb/chess-go/internal/chess"
)

func TestRenderInitialBoard(t *testing.T) {
	position := chess.NewInitialPosition()
	want := `8 ♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜
7 ♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟
6 . . . . . . . .
5 . . . . . . . .
4 . . . . . . . .
3 . . . . . . . .
2 ♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙
1 ♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖
  a b c d e f g h`

	if got := renderBoard(position.Board); got != want {
		t.Errorf("Unexpected initial board rendering:\nwant:\n%s\ngot:\n%s", want, got)
	}
}
