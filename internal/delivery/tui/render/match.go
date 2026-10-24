package render

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	chess "github.com/gabrzb/chess-go/internal/domain"
)

type MatchPanel struct {
	Width, Height            int
	Mode, CompactMode        string
	Status, CompactStatus    string
	LastMove, Footer, Notice string
	Board                    chess.Board
	BoardState               BoardState
	Turn                     TurnState
	Moves                    []string
}

func Match(panel MatchPanel) string {
	for _, scale := range []BoardScale{LargeBoard, CompactBoard} {
		content := matchContent(panel, scale)

		if Fits(content, panel.Width, panel.Height) {
			return content
		}
	}

	return matchContent(panel, CompactBoard)
}

func matchContent(panel MatchPanel, scale BoardScale) string {
	width := ContentWidth(panel.Width, MatchWidth)
	center := func(text string) string { return Center(ansi.Wrap(text, width, " "), width) }
	bar := Band(width, width, false)
	panel.BoardState.Scale = scale

	board := Board(panel.Board, panel.BoardState)
	board = Center("BOARD\n"+board, lipgloss.Width(board))
	body := Center(lipgloss.JoinHorizontal(lipgloss.Center, board, strings.Repeat(" ", ColumnGap), Turn(panel.Turn, scale)), width)
	footer := center(Paint(panel.Footer, Muted))

	compose := func(mode, info string) string {
		header := bar + "\n" + center("CHESS · "+mode) + "\n" + bar

		if panel.Notice != "" {
			info += "\n" + center(panel.Notice)
		}

		return strings.Join([]string{header, info, body, footer}, "\n\n")
	}

	if panel.Width >= WideMatchMinWidth {
		for rows := MaxHistoryRows; rows >= 0; rows-- {
			columns := []string{"GAME\n" + panel.Status, "LAST MOVE\n" + panel.LastMove}

			if rows > 0 {
				columns = append(columns, History(panel.Moves, rows))
			}

			columnWidth := (width - (len(columns)-1)*ColumnGap) / len(columns)

			for index, column := range columns {
				columns[index] = Center(ansi.Wrap(column, columnWidth, " "), columnWidth)
			}

			content := compose(panel.Mode, Center(lipgloss.JoinHorizontal(lipgloss.Top, interleave(columns, strings.Repeat(" ", ColumnGap))...), width))
			if Fits(content, panel.Width, panel.Height) {
				return content
			}
		}
	}

	return compose(panel.CompactMode, center(panel.CompactStatus+"\nLast move: "+panel.LastMove))
}

func interleave(items []string, separator string) []string {
	result := make([]string, 0, len(items)*2-1)

	for index, item := range items {
		if index > 0 {
			result = append(result, separator)
		}

		result = append(result, item)
	}

	return result
}

func History(labels []string, rows int) string {
	content := "HISTORY\nWhite   Black"

	if len(labels) == 0 {
		return content + "\n—         —"
	}

	last := (len(labels) - 1) / 2

	for full := max(0, last-rows+1); full <= last; full++ {
		white, black := labels[full*2], ""

		if full*2+1 < len(labels) {
			black = labels[full*2+1]
		}

		content += fmt.Sprintf("\n%d. %-8s %s", full+1, white, black)
	}

	return content
}
