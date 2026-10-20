package render

import (
	"fmt"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/gabrzb/chess-go/internal/chess"
)

type BoardState struct {
	Cursor, Selected chess.Square
	LegalMoves       []chess.Move
	Flipped          bool
	Scale            BoardScale
}

type BoardScale uint8

const (
	CompactBoard BoardScale = iota
	LargeBoard
	BoardGutter        = 2
	terminalCellAspect = 2
	cellContentWidth   = 2 // One marker and one piece glyph.
	compactCellHeight  = 1
	largeCellHeight    = 2
)

// Terminal columns are approximately half as wide as a text row is tall.
// Both scales preserve that ratio instead of stretching only one axis.
func (scale BoardScale) CellDimensions() (width, height int) {
	height = compactCellHeight

	if scale == LargeBoard {
		height = largeCellHeight
	}

	return height * terminalCellAspect, height
}

func Board(board chess.Board, state BoardState) string {
	var text strings.Builder
	cellWidth, cellHeight := state.Scale.CellDimensions()

	for rowIndex := 0; rowIndex < chess.BoardSize; rowIndex++ {
		row := chess.BoardSize - 1 - rowIndex

		if state.Flipped {
			row = rowIndex
		}

		cells := make([]string, 0, chess.BoardSize)

		for colIndex := 0; colIndex < chess.BoardSize; colIndex++ {
			col := colIndex

			if state.Flipped {
				col = chess.BoardSize - 1 - colIndex
			}

			square, _ := chess.A1.Offset(col, row)
			cells = append(cells, boardCell(board, square, state, cellWidth, cellHeight))
		}
		for lineIndex, line := range strings.Split(lipgloss.JoinHorizontal(lipgloss.Top, cells...), "\n") {

			if lineIndex == (cellHeight-1)/2 {
				fmt.Fprintf(&text, "%d ", row+1)
			} else {
				text.WriteString(strings.Repeat(" ", BoardGutter))
			}

			text.WriteString(line)
			text.WriteString("\n")
		}
	}

	text.WriteString(strings.Repeat(" ", BoardGutter))

	for index := range chess.BoardSize {
		file := index

		if state.Flipped {
			file = chess.BoardSize - 1 - index
		}

		padding := (cellWidth-cellContentWidth)/2 + 1
		text.WriteString(strings.Repeat(" ", padding))
		text.WriteRune(rune('a' + file))
		text.WriteString(strings.Repeat(" ", cellWidth-padding-1))
	}

	return text.String()
}

func boardCell(board chess.Board, square chess.Square, state BoardState, width, height int) string {
	piece, _ := board.PieceAt(square)
	background := LightSquare

	if (square.Column()+square.Row())%2 == 0 {
		background = DarkSquare
	}

	foreground := WhitePiece

	if piece.Color == chess.Black {
		foreground = BlackPiece
	}

	symbol := Piece(piece)

	if piece.Type == chess.NoPiece {
		symbol = " "
	}

	marker := " "
	switch {
	case square == state.Cursor:
		marker = MenuCursor
	case square == state.Selected:
		marker = "●"
	case slices.ContainsFunc(state.LegalMoves, func(move chess.Move) bool { return move.To == square }):
		marker = "·"
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(foreground)).Background(lipgloss.Color(background))
	left := (width - cellContentWidth) / 2
	right := width - cellContentWidth - left
	lines := make([]string, height)

	for row := range height {
		lines[row] = style.Render(strings.Repeat(" ", width))

		if row == (height-1)/2 {
			// Each span owns its complete style; nested ANSI resets would erase the tile.
			lines[row] = style.Render(strings.Repeat(" ", left)) + style.Foreground(lipgloss.Color(BoardMarker)).Render(marker) +
				style.Render(symbol+strings.Repeat(" ", right))
		}
	}

	return strings.Join(lines, "\n")
}
