package render

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	chess "github.com/gabrzb/chess-go/internal/domain"
)

func Card(title, body, footer, bar string, options []string, choice, width int) string {
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	content := bar + "\n\n" + center.Bold(true).Render(title) + "\n\n" + bar

	if body != "" {
		content += "\n\n" + center.Render(ansi.Wrap(body, width, " "))
	}

	if len(options) > 0 {
		content += "\n\n" + Menu(options, choice, width)
	}

	if footer != "" {
		content += "\n\n" + center.Render(Paint(ansi.Wrap(footer, width, " "), Muted))
	}

	return content
}

func Promotion(cursor chess.Square, color chess.Color, choice int, kinds []chess.PieceType, width int) string {
	bar := Band(width, width, false)
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	var items []string

	for _, kind := range kinds {
		option := Piece(chess.Piece{Type: kind, Color: color}) + "  " + PieceName(kind)
		items = append(items, option)
	}

	return bar + "\n\n" + center.Bold(true).Render("PROMOÇÃO DE PEÃO") + "\n" + center.Render("Escolha a nova peça para "+cursor.String()+".") + "\n\n" + bar + "\n\n" + Menu(items, choice, width) + "\n\n" + center.Render("↑↓ escolher · Enter confirmar\nQ/R/B/N atalhos · Esc cancelar")
}

const MenuCursor = "▸"

// Menu keeps the focus style on the arrow; labels and padding never inherit it.
func Menu(options []string, choice, width int) string {
	items := make([]string, len(options))

	for index, option := range options {
		marker := " "

		if index == choice {
			marker = Paint(MenuCursor, Accent)
		}

		items[index] = marker + " " + option
	}

	return Center(strings.Join(items, "\n"), width)
}
