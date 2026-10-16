package tui

import (
	"fmt"
	"strings"

	"github.com/gabrzb/chess-go/internal/chess"
)

func (m model) renderBoard(board chess.Board) string {
	var text strings.Builder

	for rowIndex := 0; rowIndex < chess.BoardSize; rowIndex++ {
		row := chess.BoardSize - 1 - rowIndex
		if m.flipped {
			row = rowIndex
		}
		fmt.Fprintf(&text, "%d", row+1)

		for colIndex := 0; colIndex < chess.BoardSize; colIndex++ {
			col := colIndex
			if m.flipped {
				col = chess.BoardSize - 1 - colIndex
			}
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

	if m.flipped {
		text.WriteString("  h g f e d c b a")
	} else {
		text.WriteString("  a b c d e f g h")
	}
	return text.String()
}
