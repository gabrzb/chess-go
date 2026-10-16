package render

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type MatchPanel struct {
	Width, Height                   int
	Mode, CompactMode               string
	Status, CompactStatus           string
	LastMove, Board, Footer, Notice string
	Moves                           []string
}

func Match(panel MatchPanel) string {
	width := min(76, max(18, panel.Width-4))
	center := func(text string) string { return Center(ansi.Wrap(text, width, " "), width) }
	bar := Paint(strings.Repeat("=", width), Muted)
	board := Center("TABULEIRO\n"+panel.Board, width)
	footer := center(Paint(panel.Footer, Muted))
	compose := func(mode, info string) string {
		header := bar + "\n" + center("CHESS · "+mode) + "\n" + bar
		if panel.Notice != "" {
			info += "\n" + center(panel.Notice)
		}
		return strings.Join([]string{header, info, board, footer}, "\n\n")
	}
	if panel.Width >= 72 {
		for rows := 4; rows >= 0; rows-- {
			columns := []string{"PARTIDA\n" + panel.Status, "ÚLTIMO LANCE\n" + panel.LastMove}
			if rows > 0 {
				columns = append(columns, History(panel.Moves, rows))
			}
			columnWidth := (width - (len(columns)-1)*2) / len(columns)
			for index, column := range columns {
				columns[index] = Center(ansi.Wrap(column, columnWidth, " "), columnWidth)
			}
			content := compose(panel.Mode, Center(lipgloss.JoinHorizontal(lipgloss.Top, interleave(columns, "  ")...), width))
			if lipgloss.Height(content) <= panel.Height && lipgloss.Width(content) <= panel.Width {
				return content
			}
		}
	}
	return compose(panel.CompactMode, center(panel.CompactStatus+"\nÚltimo lance: "+panel.LastMove))
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
	content := "HISTÓRICO\nBrancas   Pretas"
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
