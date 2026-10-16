package render

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/gabrzb/chess-go/internal/chess"
	"strings"
)

func Card(title, body, footer, bar string, options []string, choice, width int) string {
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	content := bar + "\n\n" + center.Bold(true).Render(title) + "\n\n" + bar
	if body != "" {
		content += "\n\n" + center.Render(ansi.Wrap(body, width, " "))
	}
	if len(options) > 0 {
		items := make([]string, len(options))
		for index, option := range options {
			items[index] = "  " + option
			if index == choice {
				items[index] = Paint("▸ "+option, Accent)
			}
		}
		menu := strings.Join(items, "\n")
		content += "\n\n" + CenteredBlock(menu, width)
	}
	if footer != "" {
		content += "\n\n" + center.Render(Paint(ansi.Wrap(footer, width, " "), Muted))
	}
	return content
}

func Promotion(cursor chess.Square, color chess.Color, choice, width int) string {
	bar := Paint(strings.Repeat("=", width), Muted)
	center := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	names := []string{"Dama", "Torre", "Bispo", "Cavalo"}
	kinds := []chess.PieceType{chess.Queen, chess.Rook, chess.Bishop, chess.Knight}
	var items []string
	for index, name := range names {
		option := Piece(chess.Piece{Type: kinds[index], Color: color}) + "  " + name
		text := "  " + option
		if index == choice {
			text = Paint("▸ "+option, Accent)
		}
		items = append(items, text)
	}
	return bar + "\n\n" + center.Bold(true).Render("PROMOÇÃO DE PEÃO") + "\n" + center.Render("Escolha a nova peça para "+cursor.String()+".") + "\n\n" + bar + "\n\n" + CenteredBlock(strings.Join(items, "\n"), width) + "\n\n" + center.Render("↑↓ escolher · Enter confirmar\nQ/R/B/N atalhos · Esc cancelar")
}
