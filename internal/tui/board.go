package tui

import (
	"fmt"
	"strings"

	"github.com/gabrzb/chess-go/internal/chess"
)

func (m model) renderBoard(board chess.Board) string {
	var text strings.Builder

	for row := chess.A8.Row(); row >= chess.A1.Row(); row-- {
		fmt.Fprintf(&text, "%d", row+1)

		for col := 0; col < chess.BoardSize; col++ {
			square, _ := chess.A1.Offset(col, row)
			piece, _ := board.PieceAt(square)

			symbol := pieceSymbol(piece)
			if square == m.cursor {
				symbol = cursorStyle.Render(symbol)
			} else if square == m.selected {
				symbol = selectedStyle.Render(symbol)
			} else if m.isLegalDestination(square) {
				symbol = legalMoveStyle.Render(symbol)
			}

			fmt.Fprintf(&text, " %s", symbol)
		}

		text.WriteString("\n")
	}

	text.WriteString("  a b c d e f g h")
	return text.String()
}
