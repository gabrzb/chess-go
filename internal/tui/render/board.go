package render

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gabrzb/chess-go/internal/chess"
)

type BoardState struct {
	Cursor, Selected chess.Square
	LegalMoves       []chess.Move
	Flipped          bool
}

func Board(board chess.Board, state BoardState) string {
	var text strings.Builder

	for rowIndex := 0; rowIndex < chess.BoardSize; rowIndex++ {
		row := chess.BoardSize - 1 - rowIndex
		if state.Flipped {
			row = rowIndex
		}
		fmt.Fprintf(&text, "%d", row+1)

		for colIndex := 0; colIndex < chess.BoardSize; colIndex++ {
			col := colIndex
			if state.Flipped {
				col = chess.BoardSize - 1 - colIndex
			}
			square, _ := chess.A1.Offset(col, row)
			piece, _ := board.PieceAt(square)

			symbol := Piece(piece)
			if square == state.Cursor {
				symbol = cursorStyle.Render(symbol)
			} else if square == state.Selected {
				symbol = selectedStyle.Render(symbol)
			} else if slices.ContainsFunc(state.LegalMoves, func(move chess.Move) bool { return move.To == square }) {
				symbol = legalMoveStyle.Render(symbol)
			}

			fmt.Fprintf(&text, " %s", symbol)
		}

		text.WriteString("\n")
	}

	if state.Flipped {
		text.WriteString("  h g f e d c b a")
	} else {
		text.WriteString("  a b c d e f g h")
	}
	return text.String()
}
