package tui

import (
	"fmt"
	"strings"

	"github.com/gabrzb/chess-go/internal/chess"
)

func renderBoard(board chess.Board) string {
	var text strings.Builder

	for row := chess.A8.Row(); row >= chess.A1.Row(); row-- {
		fmt.Fprintf(&text, "%d", row+1)

		for col := 0; col < chess.BoardSize; col++ {
			square, _ := chess.A1.Offset(col, row)
			piece, _ := board.PieceAt(square)
			fmt.Fprintf(&text, " %s", pieceSymbol(piece))
		}

		text.WriteString("\n")
	}

	text.WriteString("  a b c d e f g h")
	return text.String()
}
